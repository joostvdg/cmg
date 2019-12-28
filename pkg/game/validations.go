package game

import (
	"strconv"
	"strings"
	"time"

	"github.com/joostvdg/cmg/pkg/model"
	log "github.com/sirupsen/logrus"
)

// ValidateBoard function that validates certain attributes of a game Board
// the validate function should compare the current board against the rules for the request game
type ValidateBoard func(board *Board, gameRules GameRules) bool

var (
	Validations = []ValidateBoard{
		ValidateResourceScores,
		ValidateAdjacentTiles,
		ValidateTilesNumbers,
		ValidateResourceSpread,
	}
)

// ValidateResourceScores validates the scores of the resources
// we derived the scores from the probability scores of the Number of the tile they're associated with
// there is a maximum, and a minimum to validate, to make sure all resources fall within a certain distribution
func ValidateResourceScores(board *Board, rules GameRules) bool {
	start := time.Now()
	log.Debug(" > ValidateResourceScores start")
	isValid := true
	resourceCounts := make([]int, 6, 6)
	resourceScores := make([]int, 6, 6)

	resourceScores[0] = 0
	resourceCounts[0] = 0
	resourceScores[1] = 0
	resourceCounts[1] = 0
	resourceScores[2] = 0
	resourceCounts[2] = 0
	resourceScores[3] = 0
	resourceCounts[3] = 0
	resourceScores[4] = 0
	resourceCounts[4] = 0
	resourceScores[5] = 0
	resourceCounts[5] = 0

	for _, tile := range board.Tiles {
		codeInt, _ := strconv.Atoi(tile.Landscape.Code)
		codeInt-- // we don't use the All resource, which is 0
		resourceScores[codeInt] = resourceScores[codeInt] + tile.Number.Score
		resourceCounts[codeInt] = resourceCounts[codeInt] + 1
	}

	for resourceId, score := range resourceScores {
		skipDesertId, _ := strconv.Atoi(model.None.Code)
		skipDesertId-- // because we did the same in the resource scores/counts
		if resourceId == skipDesertId {
			// skip Desert tiles
			continue
		}

		log.WithFields(log.Fields{
			"score":          score,
			"resourceId":     resourceId,
			"resourceCounts": resourceCounts[resourceId],
		}).Debug("  - scoring for resource:")

		if resourceCounts[resourceId] == 0 {
			isValid = false
			continue
		}

		avgScore := score / resourceCounts[resourceId]
		if avgScore > rules.MaximumResourceScore || avgScore < rules.MinimumResourceScore {
			t := time.Now()
			elapsed := t.Sub(start)
			log.WithFields(log.Fields{
				"resourceId": resourceId,
				"avgScore":   avgScore,
				"Duration":   elapsed,
			}).Debug("  - Invalid scoring for resource:")
			isValid = false
		}
	}

	t := time.Now()
	elapsed := t.Sub(start)
	log.WithFields(log.Fields{
		"Duration": elapsed,
	}).Debug(" < ValidateResourceScores finish")
	return isValid
}

// ValidateAdjacentTiles validates the scores based on the tiles that lie next to each other
// we derive the scores from the probability score of the Number (based on a two dice roll)
// we do not want a too great a spot, we also do not want a too weak spot
// in addition, we also do not want too many spots (3 adjacent tiles) to be over a score of 300
// as this would signify a skewed distribution of resources and their scores
func ValidateAdjacentTiles(board *Board, rules GameRules) bool {
	start := time.Now()
	log.WithFields(log.Fields{
		"number of tile groups": len(board.GameType.AdjacentTileGroups),
	}).Debug(" > ValidateAdjacentTiles start")

	scoresOver300 := 0
	for _, tileGroup := range board.GameType.AdjacentTileGroups {
		valid, weightTotal := board.validateAdjectTileGroup(rules.MaximumScore, rules.MinimumScore, tileGroup[0], tileGroup[1], tileGroup[2])
		if weightTotal > 300 {
			scoresOver300++
		}

		log.WithFields(log.Fields{
			"weight": weightTotal,
			"G0":     tileGroup[0],
			"G1":     tileGroup[1],
			"G2":     tileGroup[2],
		}).Debug("  - Tile Group:")

		if !valid {
			t := time.Now()
			elapsed := t.Sub(start)
			log.WithFields(log.Fields{
				"Duration": elapsed,
			}).Debug(" < ValidateAdjacentTiles finish")
			return false
		}
	}

	t := time.Now()
	elapsed := t.Sub(start)
	log.WithFields(log.Fields{
		"Duration": elapsed,
	}).Debug(" < ValidateAdjacentTiles finish")
	if scoresOver300 > rules.MaxOver300 {
		return false
	}

	return true
}

// ValidateTilesNumbers validates if the number of tiles in the board matches the expected tiles for a given game type
func ValidateTilesNumbers(board *Board, rules GameRules) bool {
	log.Debug(" > ValidateTilesNumbers start")
	return len(board.Tiles) == board.GameType.TilesCount
}

// ValidateHarbors validates whether or not a harbor is linked to a resource tile with the same resource as the harbor
func ValidateHarbors(board *Board, rules GameRules) bool {
	log.Debug(" > ValidateHarbors start")
	for k, v := range board.Harbors {
		harborResource := v.Resource
		tileCodeA := k
		tileCodeB := ""
		if strings.Contains(tileCodeA, ",") {
			tileCodes := strings.Split(tileCodeA, ",")
			tileCodeA = tileCodes[0]
			tileCodeB = tileCodes[1]
		}
		if sameResource(tileCodeA, harborResource, board.Board) || sameResource(tileCodeB, harborResource, board.Board) {
			log.Debug(" < ValidateHarbors finish")
			return false
		}
	}

	log.Debug(" < ValidateHarbors finish")
	return true
}

// ValidateResourceSpread validates whether resources are spread on the board.
// There shouldn't be too many of the same resource next to each other.
func ValidateResourceSpread(board *Board, rules GameRules) bool {

	return ValidateResourcesPerColumn(board, rules) && ValidateResourcesPerRow(board, rules)
}

// ValidateResourcesPerColumn validates if there's not too many of the same landscape type
// in a Column.
//    a   b   c   d   f
// 0          x
// 1      x       x
// 2  x       x       x
// 3      x       x
// 4  x       x       x
// 5      x       x
// 6  x       x       x
// 7      x       x
// 8          x
// c0, b0, a0
// c1, b1, a1
// c2, b2, a2
// c0, d0, f0
// c1, d1, f1
// c2, d2, f2
func ValidateResourcesPerColumn(board *Board, rules GameRules) bool {

	initialRuneId := int('a')
	numberOfRows := len(board.Board)
	halfNumberOfRows := numberOfRows / 2
	// first, do left arc
	leftIsValid := validateResourcePerColumnArc(board, rules, initialRuneId)
	// then do right arc
	rightIsValid := validateResourcePerColumnArc(board, rules, initialRuneId+halfNumberOfRows)
	return leftIsValid && rightIsValid
}

func validateResourcePerColumnArc(board *Board, rules GameRules, initialRune int) bool {
	numberOfRows := len(board.Board)
	// we do c0, b0, a0 -> 5 / 2 = 2 + 1 -> 3
	numberOfRowsToCheck := (numberOfRows / 2) + 1
	numberOfColumns := len(board.Board["a"])
	for i := 0; i < numberOfColumns; i++ {
		countField := 0
		countForest := 0
		countPasture := 0
		countMountain := 0
		countHill := 0
		rowRune := initialRune

		for j := 0; j < numberOfRowsToCheck; j++ {
			rowIndex := string(rowRune)
			rowRune++
			tile := board.Board[rowIndex][i]
			switch tile.Landscape.Code {
			case model.Brick.Code:
				countHill++
			case model.Field.Code:
				countField++
			case model.Pasture.Code:
				countPasture++
			case model.Mountain.Code:
				countMountain++
			case model.Forest.Code:
				countForest++
			}
		}

		if countMountain > rules.MaxSameLandscapePerColumn ||
			countForest > rules.MaxSameLandscapePerColumn ||
			countPasture > rules.MaxSameLandscapePerColumn ||
			countField > rules.MaxSameLandscapePerColumn ||
			countHill > rules.MaxSameLandscapePerColumn {
			log.Warnf("Too many tiles of the same landscape type in a column arc: %v\n", rules.MaxSameLandscapePerColumn)
			return false
		}
	}

	return true
}

// ValidateResourcesPerRow validates if there's not too many of the same landscape type
// per row. Rows being a, b, c and so on.
func ValidateResourcesPerRow(board *Board, rules GameRules) bool {

	for _, column := range board.Board {
		if len(column) > 0 {
			countField := 0
			countForest := 0
			countPasture := 0
			countMountain := 0
			countHill := 0
			for _, tile := range column {
				switch tile.Landscape.Code {
				case model.Brick.Code:
					countHill++
				case model.Field.Code:
					countField++
				case model.Pasture.Code:
					countPasture++
				case model.Mountain.Code:
					countMountain++
				case model.Forest.Code:
					countForest++
				}
			}
			if countMountain > rules.MaxSameLandscapePerRow ||
				countForest > rules.MaxSameLandscapePerRow ||
				countPasture > rules.MaxSameLandscapePerRow ||
				countField > rules.MaxSameLandscapePerRow ||
				countHill > rules.MaxSameLandscapePerRow {
				log.Warnf("Too many tiles of the same landscape type in a Row: %v\n", rules.MaxSameLandscapePerRow)
				return false
			}
		}
	}
	return true
}
