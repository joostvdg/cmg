package webserver

import (
	"github.com/joostvdg/cmg/pkg/mapgen"
	"github.com/labstack/echo/v4"

	"net/http"
)

// GetMap starts the Generation Cycle, which may or may not succeed with a valid map according to the supplied Game Rules
func GetMap(c echo.Context) error {

	requestInfo := GetRequestInfoFromRequest(c)
	rules := GetGameRulesFromRequest(c)

	gameMap, err := mapgen.ProcessMapGenerationRequest(rules, requestInfo)
	if err != nil {
		return AbortingMapGeneration(c, rules, requestInfo)
	}

	if requestInfo.JSONP {
		return c.JSONP(http.StatusOK, requestInfo.Callback, &gameMap)
	}
	return c.JSON(http.StatusOK, &gameMap)
}
