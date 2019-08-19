package webserver

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/joostvdg/cmg/pkg/game"
	"github.com/joostvdg/cmg/pkg/webserver/model"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestGetNormalMap(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	targetPath := "/api/map"

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, targetPath, nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, GetMap(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
	res := rec.Result()
	defer res.Body.Close()
	var gameMap model.Map
	expectedGameType := game.NormalGame.Name

	if assert.NoError(t, json.Unmarshal([]byte(rec.Body.String()), &gameMap)) {
		assert.Equal(t, gameMap.GameType, expectedGameType)
		assert.NotEmpty(t, gameMap.GameCode)
		assert.Equal(t, 57, len(gameMap.GameCode))
		assert.Empty(t, gameMap.Error)
		assert.Equal(t, 5, len(gameMap.Board))
		assert.Equal(t, 3, len(gameMap.Board["a"]))
		assert.Equal(t, 4, len(gameMap.Board["b"]))
		assert.Equal(t, 5, len(gameMap.Board["c"]))
		assert.Equal(t, 4, len(gameMap.Board["d"]))
		assert.Equal(t, 3, len(gameMap.Board["e"]))
	}
}

func TestGetLargeMap(t *testing.T) {
	targetPath := "/api/map?type=large"

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, targetPath, nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, GetMap(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
	res := rec.Result()
	defer res.Body.Close()
	var gameMap model.Map
	expectedGameType := game.LargeGame.Name

	if assert.NoError(t, json.Unmarshal([]byte(rec.Body.String()), &gameMap)) {
		assert.Equal(t, gameMap.GameType, expectedGameType)
		assert.NotEmpty(t, gameMap.GameCode)
		assert.Equal(t, 90, len(gameMap.GameCode))
		assert.Empty(t, gameMap.Error)
		assert.Equal(t, 7, len(gameMap.Board))
		assert.Equal(t, 3, len(gameMap.Board["a"]))
		assert.Equal(t, 4, len(gameMap.Board["b"]))
		assert.Equal(t, 5, len(gameMap.Board["c"]))
		assert.Equal(t, 6, len(gameMap.Board["d"]))
		assert.Equal(t, 5, len(gameMap.Board["e"]))
		assert.Equal(t, 4, len(gameMap.Board["f"]))
		assert.Equal(t, 3, len(gameMap.Board["g"]))

	}
}

func TestGetNormalMapImpossibleFails(t *testing.T) {
	targetPath := "/api/map"

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, targetPath, nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, GetMap(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
	res := rec.Result()
	defer res.Body.Close()
	var gameMap model.Map
	expectedGameType := game.NormalGame.Name

	if assert.NoError(t, json.Unmarshal([]byte(rec.Body.String()), &gameMap)) {
		assert.Equal(t, gameMap.GameType, expectedGameType)
		assert.NotEmpty(t, gameMap.GameCode)
		assert.Equal(t, 57, len(gameMap.GameCode))
		assert.Empty(t, gameMap.Error)
		assert.Equal(t, 5, len(gameMap.Board))
		assert.Equal(t, 3, len(gameMap.Board["a"]))
		assert.Equal(t, 4, len(gameMap.Board["b"]))
		assert.Equal(t, 5, len(gameMap.Board["c"]))
		assert.Equal(t, 4, len(gameMap.Board["d"]))
		assert.Equal(t, 3, len(gameMap.Board["e"]))
	}
}

//if gameTypeParam == "large" {
//gameTypeValue = 1
//}
//
//min := extractIntParamOrDefault(c, "min", 165)
//max := extractIntParamOrDefault(c, "max", 361)
//max300 := extractIntParamOrDefault(c, "max300", 14)
//maxr := extractIntParamOrDefault(c, "maxr", 130)
//minr := extractIntParamOrDefault(c, "minr", 30)

func TestGetNormalMapJsonPAndCallback(t *testing.T) {
	callback := "MyFunction"
	targetPath := fmt.Sprintf("/api/map?jsonp=true&callback=%v", callback)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, targetPath, nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, GetMap(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
	res := rec.Result()
	defer res.Body.Close()
	var gameMap model.Map
	expectedGameType := game.NormalGame.Name

	prefix := fmt.Sprintf("%v(", callback)
	assert.True(t, strings.HasPrefix(rec.Body.String(), prefix))
	assert.True(t, strings.HasSuffix(rec.Body.String(), ");"))
	strippedBody := strings.TrimPrefix(rec.Body.String(), prefix)
	strippedBody = strings.TrimSuffix(strippedBody, ");")

	if assert.NoError(t, json.Unmarshal([]byte(strippedBody), &gameMap)) {
		assert.Equal(t, gameMap.GameType, expectedGameType)
	}
}
