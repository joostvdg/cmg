package webserver

import (
	"fmt"
	"github.com/getsentry/sentry-go"
	sentryecho "github.com/getsentry/sentry-go/echo"
	"github.com/joostvdg/cmg/pkg/game"
	"github.com/joostvdg/cmg/pkg/mapgen"
	"github.com/joostvdg/cmg/pkg/webserver/model"
	"github.com/labstack/echo/v4"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"time"
)

func extractIntParamOrDefault(context echo.Context, paramName string, defaultValue int) int {
	paramValue :=  context.QueryParam(paramName)
	if len(paramValue) <= 0 {
		return defaultValue
	}
	intValue, err := strconv.Atoi(paramValue)
	if err != nil {
		return defaultValue
	}
	return intValue
}

func GetMap(c echo.Context) error {
	callback := c.QueryParam("callback")
	jsonp := c.QueryParam("jsonp")
	uuid, _ := uuid.NewUUID()
	start := time.Now()

	gameTypeValue := 0
	gameTypeParam := c.QueryParam("type")
	if gameTypeParam == "large" {
		gameTypeValue = 1
	}

	min := extractIntParamOrDefault(c, "min", 165)
	max := extractIntParamOrDefault(c, "max", 361)
	max300 := extractIntParamOrDefault(c, "max300", 14)
	maxr := extractIntParamOrDefault(c, "maxr", 130)
	minr := extractIntParamOrDefault(c, "minr", 30)

	rules := game.GameRules{
		GameType:             gameTypeValue,
		MinimumScore:         min,
		MaximumScore:         max,
		MaxOver300:           max300,
		MaximumResourceScore: maxr,
		MinimumResourceScore: minr,
	}

	gameType := game.NormalGame
	if rules.GameType == 0 {
		gameType = game.NormalGame
	} else if rules.GameType == 1 {
		gameType = game.LargeGame
	}
	verbose := false
	numberOfLoops := 1
	totalGenerations := 0

	log.WithFields(log.Fields{
		"GameRules": rules,
		"UUID": uuid,
		"RequestURI": c.Request().RequestURI,
		"HOST": c.Request().Host,
		"RemoteAddr": c.Request().RemoteAddr,
	}).Info("Attempt to generate a fair map:")

	board := mapgen.MapGenerationAttempt(gameType, verbose)
	for i := 0; i < numberOfLoops; i++ {
		for !board.IsValid(rules, gameType, verbose) {
			totalGenerations++
			if totalGenerations > 1501 {
				return abortingMapGeneration(c, rules, gameType.Name, totalGenerations, jsonp, callback)
			}
			board = mapgen.MapGenerationAttempt(gameType, verbose)
		}
	}
	var content = model.Map{
		GameType: gameType.Name,
		Board:    board.Board,
	}

	t := time.Now()
	elapsed := t.Sub(start)
	log.WithFields(log.Fields{
		"UUID": uuid,
		"Duration": elapsed,
	}).Info("Created a new map")

	if jsonp == "true" {
		return c.JSONP(http.StatusOK, callback, &content)
	}
	return c.JSON(http.StatusOK, &content)
}

func abortingMapGeneration(ctx echo.Context, rules game.GameRules, gameType string, totalGenerations int, jsonp string, callback string) error {
	if hub := sentryecho.GetHubFromContext(ctx); hub != nil {
		hub.WithScope(func(scope *sentry.Scope) {
			scope.SetExtra("GameRules", rules)
			scope.SetExtra("RequestURI", ctx.Request().RequestURI )
			hub.CaptureMessage(fmt.Sprintf("Could not generate map of type %v, tried %v times", gameType, totalGenerations))
			hub.Flush(time.Second * 5)
		})
	}

	log.Warnf("Can not generate a map even after %v tries, perhaps try less strict requirements?", totalGenerations)
	var content = model.Map{
		GameType: gameType,
		Board:    nil,
	}
	if jsonp == "true" {
		return ctx.JSONP(http.StatusLoopDetected, callback, &content)
	}
	return ctx.JSON(http.StatusLoopDetected, &content)
}


