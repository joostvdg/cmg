package game

import (
	"fmt"
	"github.com/joostvdg/cmg/pkg/model"
	log "github.com/sirupsen/logrus"
	"sort"
	"strconv"
	"strings"
)

// Board the Catan game Board, contains the Tiles and how they are distributed on the Board
type Board struct {
	Tiles    []*model.Tile
	Board    map[string][]*model.Tile
	GameType GameType
	Harbors  map[string]*model.Harbor
	GameCode string
}

// IsValid wrapper function for encapsulating all the validations for the map
func (b *Board) IsValid(rules GameRules, game GameType, verbose bool) bool {
	return len(b.Tiles) == game.TilesCount &&
		b.validateAdjacentTiles(rules, verbose) &&
		b.validateResourceScores(rules, verbose) &&
		b.validateHarbors()
}

func (b *Board) validateResourceScores(rules GameRules, verbose bool) bool {
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

	for _, tile := range b.Tiles {
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
			log.WithFields(log.Fields{
				"resourceId": resourceId,
				"avgScore":   avgScore,
			}).Debug("Invalid scoring for resource:")
			isValid = false
		}
	}

	return isValid
}

func (b *Board) validateAdjacentTiles(rules GameRules, verbose bool) bool {

	scoresOver300 := 0
	for _, tileGroup := range b.GameType.AdjacentTileGroups {
		valid, weightTotal := b.validateAdjectTileGroup(rules.MaximumScore, rules.MinimumScore, tileGroup[0], tileGroup[1], tileGroup[2])
		if weightTotal > 300 {
			scoresOver300++
		}
		if verbose {
			tileGroupSet := fmt.Sprintf("[%s, %s, %s]", tileGroup[0], tileGroup[1], tileGroup[2])
			log.WithFields(log.Fields{
				"tileGroup": tileGroupSet,
				"weight":    weightTotal,
			}).Debug("Validating Tile Group:")
		}
		if !valid {
			return false
		}
	}

	if scoresOver300 > rules.MaxOver300 {
		return false
	}

	return true

}

func (b *Board) validateAdjectTileGroup(max int, min int, tileCodeA string, tileCodeB string, tileCodeC string) (bool, int) {
	weightTileA, _ := strconv.Atoi(b.element(tileCodeA))
	weightTileB, _ := strconv.Atoi(b.element(tileCodeB))
	weightTileC, _ := strconv.Atoi(b.element(tileCodeC))
	weightTotal := weightTileA + weightTileB + weightTileC
	if weightTotal > max || weightTotal < min {
		log.WithFields(log.Fields{
			"Score":       weightTotal,
			"Max allowed": max,
			"Min allowed": min,
		}).Debug("Invalid tile group")
		return false, weightTotal
	}
	return true, weightTotal
}

// validateHarbors validates whether or not a harbor is linked to a resource tile with the same resource as the harbor
func (b *Board) validateHarbors() bool {
	for k, v := range b.Harbors {
		harborResource := v.Resource
		tileCodeA := k
		tileCodeB := ""
		if strings.Contains(tileCodeA, ",") {
			tileCodes := strings.Split(tileCodeA, ",")
			tileCodeA = tileCodes[0]
			tileCodeB = tileCodes[1]
		}
		if sameResource(tileCodeA, harborResource, b.Board) || sameResource(tileCodeB, harborResource, b.Board) {
			return false
		}
	}
	return true
}

func sameResource(tileCode string, harborResource model.Resource, board map[string][]*model.Tile) bool {
	if tileCode == "" {
		return false
	}
	runeCode := []rune(tileCode)
	column := string(runeCode[0:1])
	row, _ := strconv.Atoi(string(runeCode[1:2]))
	if board[column][row].Landscape.Resource == harborResource {
		return true
	}

	return false
}

func (b *Board) PrintToConsole() {
	b.GameType.ToConsole(b)
}

func (board *Board) element(code string) string {
	runeCode := []rune(code)
	row, _ := strconv.Atoi(string(runeCode[0:1]))
	column := string(runeCode[1:2])
	elementType := string(runeCode[2:3])
	switch elementType {
	case "l":
		return fmt.Sprintf("%v", board.Board[column][row].Landscape)
	case "n":
		number := board.Board[column][row].Number.Number
		padding := ""
		if number < 10 {
			padding = "."
		}
		element := fmt.Sprintf("%s%v", padding, number)
		return element
	case "w": // weight of the number
		return fmt.Sprintf("%v", board.Board[column][row].Number.Score)
	default: // todo: panic?
		return ""
	}

}

func (board *Board) GetGameCode(delimiter bool) string {
	if board.GameCode == "" {
		code := ""

		rows := make([]string, 0)
		for row := range board.Board {
			rows = append(rows, row)
		}
		sort.Strings(rows)
		for _, rowKey := range rows {
			for _, tile := range board.Board[rowKey] {
				code += fmt.Sprintf("%v", tile.Landscape.Code)
				code += tile.Number.Code
				code += fmt.Sprintf("%v", tile.Harbor.Code)
			}
			if delimiter {
				code += "_"
			}
		}
		board.GameCode = code
	}
	return board.GameCode
}
