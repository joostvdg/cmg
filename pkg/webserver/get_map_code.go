package webserver

import (
	"net/http"

	"github.com/joostvdg/cmg/pkg/mapgen"
	"github.com/joostvdg/cmg/pkg/webserver/model"
	"github.com/labstack/echo/v4"
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
