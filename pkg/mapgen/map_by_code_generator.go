package mapgen

import (
	"github.com/joostvdg/cmg/pkg/game"
	log "github.com/sirupsen/logrus"
)

func GenerateBoardByGameCode(rules game.GameRules) game.Board {
	log.Debug(" > GenerateBoardByGameCode start")
	totalGenerations := 0
	code := game.GenerateGameCodeNormalGame()
	log.Debugf("GameCode: %v", code)

	var board game.Board
	gameType := game.NormalGame
	var err error
	for i := 0; i < rules.Generations; i++ {
		board, err = game.InflateNormalGameFromCode(code, &gameType)
		if err != nil {
			log.Debugf("Board not inflated correctly %v", err)
			continue
		}

		if board.IsValid(rules, gameType) {
			log.Info("Required iterations: ", totalGenerations)
			board.GameCode = code
			return board
		}
		totalGenerations++
		code = game.GenerateGameCodeNormalGame()
	}
	log.Debug(" > GenerateBoardByGameCode finish")
	return board
}
