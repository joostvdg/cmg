package game

import (
	"errors"
	"fmt"
	"strings"

	"github.com/joostvdg/cmg/pkg/model"
	"github.com/prometheus/common/log"
)

func inflateGameFromCode(code string, gameLayout []int) (Board, error) {
	var boardMap [][]*model.Tile
	boardMap = make([][]*model.Tile, 0)

	// if we found a delimiter, and we happen to have exactly the number expected, we're probably good
	if strings.Contains(code, DefaultGameRulesNormal.Delimiter) && strings.Count(code, DefaultGameRulesNormal.Delimiter) == len(gameLayout) {
		code = strings.Replace(code, DefaultGameRulesNormal.Delimiter, "", len(gameLayout))
	}

	codeIndex := 0
	for column, row := range gameLayout {
		numberOfTiles := row
		tiles := make([]*model.Tile, numberOfTiles, numberOfTiles)
		for i := 0; i < numberOfTiles; i++ {

			landscapeCode := code[codeIndex : codeIndex+1]
			landscape := model.Landscapes[landscapeCode]
			if landscape.Name == "" {
				errorMessage := fmt.Sprintf("Inflation error: %v is not a valid code for a Landscape", landscapeCode)
				log.Warn(errorMessage)
				return Board{}, errors.New(errorMessage)
			}

			numberCode := code[codeIndex+1 : codeIndex+2]
			number := model.Numbers[numberCode]
			if number.Score == 0 {
				errorMessage := fmt.Sprintf("Inflation error: %v is not a valid code for a Number", numberCode)
				log.Warn(errorMessage)
				return Board{}, errors.New(errorMessage)
			}

			harborCode := code[codeIndex+2 : codeIndex+3]
			harbor := model.Harbors[harborCode]
			if harbor.Name == "" {
				errorMessage := fmt.Sprintf("Inflation error: %v is not a valid code for a Harbor", harborCode)
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
		}
		boardMap[column] = tiles
	}

	board := Board{
		Board: boardMap,
	}
	return board, nil
}
