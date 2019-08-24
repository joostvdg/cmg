package game

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/joostvdg/cmg/pkg/model"
	log "github.com/sirupsen/logrus"
)

// Board the Catan game Board, contains the Tiles and how they are distributed on the Board
type Board struct {
	Tiles     []*model.Tile
	Board     [][]*model.Tile
	GameType  *GameType
	Harbors   []*model.Harbor
	GameCode  string
	WaitGroup sync.WaitGroup
}

// IsValid wrapper function for encapsulating all the validations for the map
func (b *Board) IsValid(rules *GameRules, game *GameType) bool {

	// TODO parallelize via go routines
	start := time.Now()

	isValid := true
	log.Debug("Validating map")
	for _, validationFunc := range Validations {
		b.WaitGroup.Add(1)
		go func(validation ValidateBoard) {
			defer b.WaitGroup.Done()
			valid := validation(b, rules)
			if !valid {
				isValid = false
			}
		}(validationFunc)
	}
	log.Debug("Wait for validations to finish")
	b.WaitGroup.Wait()

	t := time.Now()
	elapsed := t.Sub(start)
	log.WithFields(log.Fields{
		"Valid":       isValid,
		"Validations": len(Validations),
		"Duration":    elapsed,
		"Rules":       rules,
	}).Debug("Validated map")
	return isValid
}

func (b *Board) validateAdjacentTiles(max int, min int, tileCodeA *model.TileCode, tileCodeB *model.TileCode, tileCodeC *model.TileCode) (bool, int) {
	weightTileA := b.tileResourceProbabilityScore(*tileCodeA)
	weightTileB := b.tileResourceProbabilityScore(*tileCodeB)
	weightTileC := b.tileResourceProbabilityScore(*tileCodeC)
	weightTotal := weightTileA + weightTileB + weightTileC
	if weightTotal > max || weightTotal < min {
		log.WithFields(log.Fields{
			"Score":       weightTotal,
			"Max allowed": max,
			"Min allowed": min,
		}).Debug("  - Invalid tile group")
		return false, weightTotal
	}
	return true, weightTotal
}

//func (b *Board) validateAdjacentTiles2(max int, min int, tileCodeA string, tileCodeB string, tileCodeC string) (bool, int) {
//	weightTileA, _ := strconv.Atoi(b.element(tileCodeA))
//	weightTileB, _ := strconv.Atoi(b.element(tileCodeB))
//	weightTileC, _ := strconv.Atoi(b.element(tileCodeC))
//	weightTotal := weightTileA + weightTileB + weightTileC
//	if weightTotal > max || weightTotal < min {
//		log.WithFields(log.Fields{
//			"Score":       weightTotal,
//			"Max allowed": max,
//			"Min allowed": min,
//		}).Debug("  - Invalid tile group")
//		return false, weightTotal
//	}
//	return true, weightTotal
//}

func sameResource(tileCode *model.TileCode, resource model.Resource, board [][]*model.Tile) bool {

	if board[tileCode.Column][tileCode.Row].Landscape.Resource == resource {
		return true
	}

	return false
}

//func sameResource(tileCode string, resource model.Resource, board map[string][]*model.Tile) bool {
//	if tileCode == "" {
//		return false
//	}
//	runeCode := []rune(tileCode)
//	column := string(runeCode[0:1])
//	row, _ := strconv.Atoi(string(runeCode[1:2]))
//	if board[column][row].Landscape.Resource == resource {
//		return true
//	}
//
//	return false
//}

//func (b *Board) PrintToConsole() {
//	b.GameType.ToConsole(b)
//}

func (board *Board) tileResourceProbabilityScore(tileCode model.TileCode) int {
	return board.Board[tileCode.Column][tileCode.Row].Number.Score
}

func (board *Board) element(code string) string {
	runeCode := []rune(code)
	row, _ := strconv.Atoi(string(runeCode[0:1]))
	column, _ := strconv.Atoi(string(runeCode[1:2]))
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

		for _, row := range board.Board {
			for _, tile := range row {
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
