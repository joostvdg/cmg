package webserver

import (
	"fmt"
	"github.com/getsentry/sentry-go"
	sentryecho "github.com/getsentry/sentry-go/echo"
	"github.com/google/uuid"
	"github.com/joostvdg/cmg/pkg/game"
	"github.com/joostvdg/cmg/pkg/webserver/model"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

// GetMap starts the Generation Cycle, which may or may not succeed with a valid map according to the supplied Game Rules
func GetMapByCode(ctx echo.Context) error {
	code := ctx.Param("code")
	callback := ctx.QueryParam("callback")
	jsonp := ctx.QueryParam("jsonp")
	requestUuid, _ := uuid.NewUUID()
	start := time.Now()

	log.WithFields(log.Fields{
		"UUID":       requestUuid,
		"Code":       code,
		"RequestURI": ctx.Request().RequestURI,
		"HOST":       ctx.Request().Host,
		"RemoteAddr": ctx.Request().RemoteAddr,
	}).Info("Attempt to inflate map from a code:")

	gameType := game.NormalGame

	// TODO: have a map with game types
	var board = game.Board{}
	switch len(code) {
	case 57:
		// normal game
		gameType = game.NormalGame
		inflatedBoard, err := game.InflateNormalGameFromCode(code)
		if err != nil {
			return invalidGameCode(ctx, "Invalid code value", code, jsonp, callback)
		}
		board = inflatedBoard
	case 90:
		// large game
		gameType = game.LargeGame
	default:
		return invalidGameCode(ctx, "Unrecognizable game code", code, jsonp, callback)
	}

	var content = model.Map{
		GameType: gameType.Name,
		Board:    board.Board,
		GameCode: code,
	}

	t := time.Now()
	elapsed := t.Sub(start)
	log.WithFields(log.Fields{
		"UUID":     requestUuid,
		"Duration": elapsed,
	}).Info("Created a new map")

	if jsonp == "true" {
		return ctx.JSONP(http.StatusOK, callback, &content)
	}
	return ctx.JSON(http.StatusOK, &content)
}

func invalidGameCode(ctx echo.Context, reason string, code string, jsonp string, callback string) error {
	message := fmt.Sprintf("Could not inflate map base on game code %v, reason: %v", code, reason)
	if hub := sentryecho.GetHubFromContext(ctx); hub != nil {
		hub.WithScope(func(scope *sentry.Scope) {
			scope.SetExtra("Code", code)
			scope.SetExtra("RequestURI", ctx.Request().RequestURI)
			hub.CaptureMessage(message)
			hub.Flush(time.Second * 5)
		})
	}

	log.Warn(message)
	var content = model.Map{
		GameCode: code,
		Error:    message,
	}
	if jsonp == "true" {
		return ctx.JSONP(http.StatusBadRequest, callback, &content)
	}
	return ctx.JSON(http.StatusBadRequest, &content)
}
