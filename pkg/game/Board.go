package game

import (
	"fmt"
	"github.com/joostvdg/cmg/pkg/model"
	log "github.com/sirupsen/logrus"
	"strconv"
	"strings"
)

type Board struct {
	Tiles    []*model.Tile
	Board    map[string][]*model.Tile
	GameType GameType
	Harbors  map[string]*model.Harbor
}

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
	resourceScores[model.Desert] = 0
	resourceScores[model.Forest] = 0
	resourceScores[model.Pasture] = 0
	resourceScores[model.Field] = 0
	resourceScores[model.River] = 0
	resourceScores[model.Mountain] = 0
	resourceCounts[model.Desert] = 0
	resourceCounts[model.Forest] = 0
	resourceCounts[model.Pasture] = 0
	resourceCounts[model.Field] = 0
	resourceCounts[model.River] = 0
	resourceCounts[model.Mountain] = 0

	for _, tile := range b.Tiles {
		switch tile.Landscape {
		case model.Forest:
			resourceScores[model.Forest] = resourceScores[model.Forest] + tile.Number.Score
			resourceCounts[model.Forest] = resourceCounts[model.Forest] + 1
		case model.Pasture:
			resourceScores[model.Pasture] = resourceScores[model.Pasture] + tile.Number.Score
			resourceCounts[model.Pasture] = resourceCounts[model.Pasture] + 1
		case model.Field:
			resourceScores[model.Field] = resourceScores[model.Field] + tile.Number.Score
			resourceCounts[model.Field] = resourceCounts[model.Field] + 1
		case model.River:
			resourceScores[model.River] = resourceScores[model.River] + tile.Number.Score
			resourceCounts[model.River] = resourceCounts[model.River] + 1
		case model.Mountain:
			resourceScores[model.Mountain] = resourceScores[model.Mountain] + tile.Number.Score
			resourceCounts[model.Mountain] = resourceCounts[model.Mountain] + 1
		}
	}

	for resourceId, score := range resourceScores {
		if resourceId == 0 {
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
	if board[column][row].Resource == harborResource {
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
