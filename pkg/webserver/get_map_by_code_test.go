package webserver

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/joostvdg/cmg/pkg/game"
	"github.com/joostvdg/cmg/pkg/webserver/model"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

var (
	baseApiPath      = "/api"
	mapByCodeApiPath = "map/code"
)

func TestCodeIsUnrecognizable(t *testing.T) {
	unrecognizableCode := "abc"
	targetPath := fmt.Sprintf("%v/%v/%v", baseApiPath, mapByCodeApiPath, unrecognizableCode)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, targetPath, nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("code")
	c.SetParamValues(unrecognizableCode)
	if assert.NoError(t, GetMapByCode(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	}
	res := rec.Result()
	defer res.Body.Close()
	var gameMap model.Map

	expectedError := fmt.Sprintf("Could not inflate map base on game code %v, reason: Unrecognizable game code", unrecognizableCode)
	if assert.NoError(t, json.Unmarshal([]byte(rec.Body.String()), &gameMap)) {
		assert.Equal(t, expectedError, gameMap.Error)
	}
}

func TestCodeIsInvalid(t *testing.T) {
	invalidCode := fmt.Sprintf("%057d", 1)
	targetPath := fmt.Sprintf("%v/%v/%v", baseApiPath, mapByCodeApiPath, invalidCode)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, targetPath, nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("code")
	c.SetParamValues(invalidCode)
	if assert.NoError(t, GetMapByCode(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	}
	res := rec.Result()
	defer res.Body.Close()
	var gameMap model.Map

	expectedError := fmt.Sprintf("Could not inflate map base on game code %v, reason: Invalid code value", invalidCode)
	if assert.NoError(t, json.Unmarshal([]byte(rec.Body.String()), &gameMap)) {
		assert.Equal(t, gameMap.Error, expectedError)
	}
}

func TestCodeValidNormalGame(t *testing.T) {
	gameCode := "2j65f64a62e41b04c61h63i65d63f63g61h62d05g23b11e36z04c52i0"
	targetPath := fmt.Sprintf("%v/%v/%v", baseApiPath, mapByCodeApiPath, gameCode)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, targetPath, nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("code")
	c.SetParamValues(gameCode)
	if assert.NoError(t, GetMapByCode(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
	res := rec.Result()
	defer res.Body.Close()
	var gameMap model.Map
	expectedGameType := game.NormalGame.Name

	if assert.NoError(t, json.Unmarshal([]byte(rec.Body.String()), &gameMap)) {
		assert.Equal(t, gameMap.GameCode, gameCode)
		assert.Equal(t, gameMap.GameType, expectedGameType)
	}
}

func TestCodeValidNormalGameAlt(t *testing.T) {
	gameCode := "1g53d02b04c12f31e01d65i65h62g01a24e65b64h62c63i66z63f63j4"
	targetPath := fmt.Sprintf("%v/%v/%v", baseApiPath, mapByCodeApiPath, gameCode)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, targetPath, nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("code")
	c.SetParamValues(gameCode)
	if assert.NoError(t, GetMapByCode(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
	res := rec.Result()
	defer res.Body.Close()
	var gameMap model.Map
	expectedGameType := game.NormalGame.Name

	if assert.NoError(t, json.Unmarshal([]byte(rec.Body.String()), &gameMap)) {
		assert.Equal(t, gameMap.GameCode, gameCode)
		assert.Equal(t, gameMap.GameType, expectedGameType)
	}
}
