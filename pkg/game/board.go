package game

import (
	"fmt"
	"github.com/joostvdg/cmg/pkg/model"
	log "github.com/sirupsen/logrus"
	"sort"
	"strconv"
	"sync"
	"time"
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
func (b *Board) IsValid(rules GameRules, game GameType) bool {

	// TODO parallelize via go routines
	start := time.Now()

	isValid := true
	var wg sync.WaitGroup
	wg.Add(len(Validations))
	for _, validationFunc := range Validations {
		go func(board *Board, gameRules GameRules) {
			valid := validationFunc(board, gameRules)
			if !valid {
				isValid = false
			}
			wg.Done()
		}(b, rules)
	}
	wg.Wait()

	t := time.Now()
	elapsed := t.Sub(start)
	log.WithFields(log.Fields{
		"Valid":       isValid,
		"Validations": len(Validations),
		"Duration":    elapsed,
	}).Info("Validated map")
	return isValid
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
