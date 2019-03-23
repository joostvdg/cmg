package webserver

import (
	"github.com/joostvdg/cmg/pkg/game"
	"github.com/joostvdg/cmg/pkg/mapgen"
	"github.com/joostvdg/cmg/pkg/webserver/model"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"net/http"
	"strconv"
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

	board := mapgen.MapGenerationAttempt(gameType, verbose)
	for i := 0; i < numberOfLoops; i++ {
		for !board.IsValid(rules, gameType, verbose) {
			totalGenerations++
			if totalGenerations > 1501 {
				log.Fatal("Can not generate a map... (1000+ runs)")
			}
			board = mapgen.MapGenerationAttempt(gameType, verbose)
		}
	}
	var content = model.Map{
		GameType: gameType.Name,
		Board:    board.Board,
	}

	if jsonp == "true" {
		return c.JSONP(http.StatusOK, callback, &content)
	}
	return c.JSON(http.StatusOK, &content)
}


