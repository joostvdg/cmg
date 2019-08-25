package mapgen

import (
	"fmt"
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

	gameTypeTime := time.Now()
	gameTypeElapsed := gameTypeTime.Sub(start)
	log.WithFields(log.Fields{
		"Duration": gameTypeElapsed,
	}).Debug("Setup Game Type ")

	verbose := false
	totalGenerations := 0

	board := MapGenerationAttempt(gameType, verbose)
	elapsedGen, _ := time.ParseDuration("0ns")
	for !board.IsValid(rules, gameType) {
		totalGenerations++
		if totalGenerations > rules.Generations {
			return model.Map{}, errors.New("Stuck in generation loop")
		}
		startGen := time.Now()
		board = MapGenerationAttempt(gameType, verbose)
		finishGen := time.Now()
		elapsedGen += finishGen.Sub(startGen)
	}

	var content = model.Map{
		GameType: gameType.Name,
		Board:    board.Board,
		GameCode: board.GetGameCode(requestInfo.Delimiter),
	}

	t := time.Now()
	elapsed := t.Sub(start)
	avgDurationNanaseconds := int(elapsedGen.Nanoseconds()) / totalGenerations
	avgDuration, _ := time.ParseDuration(fmt.Sprintf("%dns", avgDurationNanaseconds))

	log.WithFields(log.Fields{
		"RequestId":              requestInfo.RequestId,
		"Total Generations":      totalGenerations,
		"Avg. Creation Duration": avgDuration,
		"Total Duration":         elapsed,
	}).Info("Created a new map")

	return content, nil
}
