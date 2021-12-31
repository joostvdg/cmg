package game

import (
	"errors"
	"fmt"
	"github.com/kennygrant/sanitize"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/joostvdg/cmg/pkg/model"
	log "github.com/sirupsen/logrus"
)

func inflateGameFromCode(code string, gameLayout map[string]int, gameType *GameType) (Board, error) {
	start := time.Now()
	log.Debug(" > Inflate Game from Game Code start")
	boardMap := make(map[string][]*model.Tile)

	// if we found a delimiter, and we happen to have exactly the number expected, we're probably good
	if strings.Contains(code, DefaultGameRulesNormal.Delimiter) && strings.Count(code, DefaultGameRulesNormal.Delimiter) == len(gameLayout) {
		code = strings.Replace(code, DefaultGameRulesNormal.Delimiter, "", len(gameLayout))
	}

	codeIndex := 0

	columns := make([]string, 0)
	for column := range gameLayout {
		columns = append(columns, column)
	}
	sort.Strings(columns)

	// TODO do we need these
	//allHarbors := make([]*model.Harbor, 0)
	allTiles := make([]*model.Tile, 0)
	for _, column := range columns {
		numberOfTiles := gameLayout[column]
		tiles := make([]*model.Tile, numberOfTiles)
		for i := 0; i < numberOfTiles; i++ {

			landscapeCode := code[codeIndex : codeIndex+1]
			landscape := model.Landscapes[landscapeCode]
			if landscape.Name == "" {
				errorMessage := fmt.Sprintf("Inflation error: %s is not a valid code for a Landscape", sanitize.Name(landscapeCode))
				log.Warn(errorMessage)
				return Board{}, errors.New(errorMessage)
			}

			numberCode := code[codeIndex+1 : codeIndex+2]
			number := model.Numbers[numberCode]
			if number.Score == 0 {
				errorMessage := fmt.Sprintf("Inflation error: %s is not a valid code for a Number", sanitize.Name(numberCode))
				log.Warn(errorMessage)
				return Board{}, errors.New(errorMessage)
			}

			harborCode := code[codeIndex+2 : codeIndex+3]
			harbor := model.Harbors[harborCode]
			if harbor.Name == "" {
				errorMessage := fmt.Sprintf("Inflation error: %s is not a valid code for a Harbor", sanitize.Name(harborCode))
				log.Warn(errorMessage)
				return Board{}, errors.New(errorMessage)
			}
			codeIndex += 3

			tile := model.Tile{
				Landscape: landscape,
				Harbor:    harbor,
				Number:    number,
			}
			tiles[i] = &tile

			allTiles = append(allTiles, &tile)
			// TODO do we need these?
			//allHarbors = append(allHarbors, &harbor)
		}
		boardMap[column] = tiles
	}

	var group sync.WaitGroup
	board := Board{
		Board:     boardMap,
		GameType:  *gameType,
		WaitGroup: group,
		Tiles:     allTiles,
	}
	t := time.Now()
	elapsed := t.Sub(start)
	log.WithFields(log.Fields{
		"Duration": elapsed,
	}).Debug(" < Inflate Game from Game Code finish")
	return board, nil
}
