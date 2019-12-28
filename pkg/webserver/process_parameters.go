package webserver

import (
	"strconv"

	"github.com/google/uuid"
	"github.com/joostvdg/cmg/pkg/game"
	"github.com/joostvdg/cmg/pkg/webserver/model"
	"github.com/labstack/echo/v4"
)

func extractIntParamOrDefault(context echo.Context, paramName string, defaultValue int) int {
	paramValue := context.QueryParam(paramName)
	if len(paramValue) <= 0 {
		return defaultValue
	}
	intValue, err := strconv.Atoi(paramValue)
	if err != nil {
		return defaultValue
	}
	return intValue
}

func GetGameRulesFromRequest(c echo.Context) game.GameRules {
	gameTypeValue := 0
	gameTypeParam := c.QueryParam("type")
	if gameTypeParam == "large" {
		gameTypeValue = 1
	}

	min := extractIntParamOrDefault(c, "min", game.DefaultGameRulesNormal.MinimumScore)
	max := extractIntParamOrDefault(c, "max", game.DefaultGameRulesNormal.MaximumScore)
	max300 := extractIntParamOrDefault(c, "max300", game.DefaultGameRulesNormal.MaxOver300)
	maxr := extractIntParamOrDefault(c, "maxr", game.DefaultGameRulesNormal.MaximumResourceScore)
	minr := extractIntParamOrDefault(c, "minr", game.DefaultGameRulesNormal.MinimumResourceScore)
	maxRow := extractIntParamOrDefault(c, "maxRow", game.DefaultGameRulesNormal.MaxSameLandscapePerRow)
	maxColumn := extractIntParamOrDefault(c, "maxColumn", game.DefaultGameRulesNormal.MaxSameLandscapePerColumn)
	adjacentSame := extractIntParamOrDefault(c, "adjacentSame", game.DefaultGameRulesNormal.AdjacentSame)

	if gameTypeParam == "large" {
		min = extractIntParamOrDefault(c, "min", game.DefaultGameRulesLarge.MinimumScore)
		max = extractIntParamOrDefault(c, "max", game.DefaultGameRulesLarge.MaximumScore)
		max300 = extractIntParamOrDefault(c, "max300", game.DefaultGameRulesLarge.MaxOver300)
		maxr = extractIntParamOrDefault(c, "maxr", game.DefaultGameRulesLarge.MaximumResourceScore)
		minr = extractIntParamOrDefault(c, "minr", game.DefaultGameRulesLarge.MinimumResourceScore)
		maxRow = extractIntParamOrDefault(c, "maxRow", game.DefaultGameRulesLarge.MaxSameLandscapePerRow)
		maxColumn = extractIntParamOrDefault(c, "maxColumn", game.DefaultGameRulesLarge.MaxSameLandscapePerColumn)
		adjacentSame = extractIntParamOrDefault(c, "adjacentSame", game.DefaultGameRulesLarge.AdjacentSame)
	}

	rules := game.GameRules{
		GameType:                  gameTypeValue,
		MinimumScore:              min,
		MaximumScore:              max,
		MaxOver300:                max300,
		MaximumResourceScore:      maxr,
		MinimumResourceScore:      minr,
		MaxSameLandscapePerRow:    maxRow,
		MaxSameLandscapePerColumn: maxColumn,
		AdjacentSame:              adjacentSame,
		Generations:               game.DefaultGameRulesNormal.Generations,
		GameTypeString:            gameTypeParam,
	}

	return rules
}

func GetRequestInfoFromRequest(c echo.Context) model.RequestInfo {
	callback := c.QueryParam("callback")
	jsonpInput := c.QueryParam("jsonp")
	jsonp := false
	delimiterInput := c.QueryParam("delimiter")
	delimiter := false
	requestId, _ := uuid.NewUUID()

	if jsonpInput == "true" {
		jsonp = true
	}

	if delimiterInput == "true" {
		delimiter = true
	}

	return model.RequestInfo{
		RequestId:  requestId,
		Callback:   callback,
		JSONP:      jsonp,
		Delimiter:  delimiter,
		RequestURI: c.Request().RequestURI,
		Host:       c.Request().Host,
		RemoteAddr: c.Request().RemoteAddr,
	}
}
