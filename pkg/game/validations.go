package game

import (
	"strconv"
	"strings"
	`time`

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
		avgScore := score / resourceCounts[resourceId]
		if avgScore > rules.MaximumResourceScore || avgScore < rules.MinimumResourceScore {
			t := time.Now()
			elapsed := t.Sub(start)
			log.WithFields(log.Fields{
				"resourceId": resourceId,
				"avgScore":   avgScore,
				"Duration":  elapsed,
			}).Debug("  - Invalid scoring for resource:")
			isValid = false
		}
	}

	t := time.Now()
	elapsed := t.Sub(start)
	log.WithFields(log.Fields{
		"Duration":  elapsed,
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
				"Duration":  elapsed,
			}).Info(" < ValidateAdjacentTiles finish")
			return false
		}
	}

	t := time.Now()
	elapsed := t.Sub(start)
	log.WithFields(log.Fields{
		"Duration":  elapsed,
	}).Info(" < ValidateAdjacentTiles finish")
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
