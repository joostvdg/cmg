package webserver

import (
	"github.com/joostvdg/cmg/pkg/game"
	"github.com/joostvdg/cmg/pkg/mapgen"
	"github.com/joostvdg/cmg/pkg/webserver/model"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"net/http"
)

func GetMap(c echo.Context) error {
	callback := c.QueryParam("callback")
	rules := game.GameRules{
		GameType:             0,
		MinimumScore:         165,
		MaximumScore:         361,
		MaxOver300:           14,
		MaximumResourceScore: 130,
		MinimumResourceScore: 30,
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
		GameType: "Normal",
		Board:    board.Board,
	}

	return c.JSONP(http.StatusOK, callback, &content)
}
