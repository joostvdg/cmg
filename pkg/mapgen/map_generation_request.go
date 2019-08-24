package mapgen

import (
	"time"

	"github.com/go-errors/errors"
	"github.com/joostvdg/cmg/pkg/game"
	"github.com/joostvdg/cmg/pkg/webserver/model"
	log "github.com/sirupsen/logrus"
)

func ProcessMapGenerationRequest(rules game.GameRules, requestInfo model.RequestInfo) (model.Map, error) {
	start := time.Now()

	log.WithFields(log.Fields{
		"GameRules":  rules,
		"RequestId":  requestInfo.RequestId,
		"RequestURI": requestInfo.RequestURI,
		"HOST":       requestInfo.Host,
		"RemoteAddr": requestInfo.RemoteAddr,
	}).Info("Attempt to generate a fair map:")

	gameType := game.NormalGame
	if rules.GameType == 0 {
		gameType = game.NormalGame
	} else if rules.GameType == 1 {
		gameType = game.LargeGame
	}
	verbose := false
	numberOfLoops := 1
	totalGenerations := 0

	board := MapGenerationAttempt(gameType, verbose)
	for i := 0; i < numberOfLoops; i++ {
		for !board.IsValid(rules, gameType) {
			totalGenerations++
			if totalGenerations > rules.Generations {
				return model.Map{}, errors.New("Stuck in generation loop")
			}
			board = MapGenerationAttempt(gameType, verbose)
		}
	}

	log.WithFields(log.Fields{
		"NumberOfTiles": len(board.Tiles),
		"Tiles": board.Tiles,
	}).Debug("Final Board: ")

	var content = model.Map{
		GameType: gameType.Name,
		Board:    board.Board,
		GameCode: board.GetGameCode(requestInfo.Delimiter),
	}

	t := time.Now()
	elapsed := t.Sub(start)
	log.WithFields(log.Fields{
		"RequestId": requestInfo.RequestId,
		"Duration":  elapsed,
	}).Info("Created a new map")

	return content, nil
}
