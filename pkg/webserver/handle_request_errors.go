package webserver

import (
	"fmt"
	"net/http"
	"time"

	"github.com/getsentry/sentry-go"
	sentryecho "github.com/getsentry/sentry-go/echo"
	"github.com/joostvdg/cmg/pkg/game"
	"github.com/joostvdg/cmg/pkg/webserver/model"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

// AbortingMapGeneration handles aborting the attempt to generate a map, handle the error and send a response to the affected client
func AbortingMapGeneration(ctx echo.Context, rules game.GameRules, requestInfo model.RequestInfo) error {
	if hub := sentryecho.GetHubFromContext(ctx); hub != nil {
		hub.WithScope(func(scope *sentry.Scope) {
			scope.SetExtra("GameRules", rules)
			scope.SetExtra("RequestURI", ctx.Request().RequestURI)
			hub.CaptureMessage(fmt.Sprintf("Could not generate map of type %v, tried %v times", rules.GameTypeString, rules.Generations))
			hub.Flush(time.Second * 5)
		})
	}

	message := fmt.Sprintf("Can not generate a map even after %v tries, perhaps try less strict requirements?", rules.Generations)
	log.Warn(message)
	var content = model.Map{
		GameType: rules.GameTypeString,
		Board:    nil,
		Error:    message,
	}
	if requestInfo.JSONP {
		return ctx.JSONP(http.StatusLoopDetected, requestInfo.Callback, &content)
	}
	return ctx.JSON(http.StatusLoopDetected, &content)
}
