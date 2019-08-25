package webserver

import (
	"net/http"
	"time"

	"github.com/joostvdg/cmg/pkg/mapgen"
	"github.com/joostvdg/cmg/pkg/webserver/model"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

// GetMapCode starts the Generation Cycle, which may or may not succeed with a valid map according to the supplied Game Rules
// It is different from GetMap because it only returns the game code
func GetMapCode(c echo.Context) error {
	requestInfo := GetRequestInfoFromRequest(c)
	rules := GetGameRulesFromRequest(c)

	wholeMap, err := mapgen.ProcessMapGenerationRequest(rules, requestInfo)
	if err != nil {
		return AbortingMapGeneration(c, rules, requestInfo)
	}

	gameCode := model.GameCode{GameCode: wholeMap.GameCode}

	if requestInfo.JSONP {
		return c.JSONP(http.StatusOK, requestInfo.Callback, &gameCode)
	}
	return c.JSON(http.StatusOK, &gameCode)
}

// GetMapViaCodeGeneration Alternative approach to generating the map
// instead of generating the map via the structs, we generate a game code
// and then inflate the game code before validating the board
// WARNING: currently only supports default normal map
func GetMapViaCodeGeneration(c echo.Context) error {
	start := time.Now()
	log.Info(" > Generate Game by Game Code start")
	requestInfo := GetRequestInfoFromRequest(c)
	rules := GetGameRulesFromRequest(c)

	board := mapgen.GenerateBoardByGameCode(rules)
	var content = model.Map{
		GameType: board.GameType.Name,
		Board:    board.Board,
		GameCode: board.GameCode,
	}

	t := time.Now()
	elapsed := t.Sub(start)
	log.WithFields(log.Fields{
		"Duration": elapsed,
	}).Info(" < Generate Game by Game Code finish")

	if requestInfo.JSONP {
		return c.JSONP(http.StatusOK, requestInfo.Callback, &content)
	}
	return c.JSON(http.StatusOK, &content)
}
