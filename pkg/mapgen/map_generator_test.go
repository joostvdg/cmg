package mapgen

import (
	"github.com/joostvdg/cmg/pkg/game"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGameCode(t *testing.T) {
	board := MapGenerationAttempt(game.NormalGame, false)
	gameCode := board.GetGameCode(false)
	expectedLengthOfGameCode := game.NormalGame.TilesCount * 3 // landscape, number, harbor
	assert.Equal(t, expectedLengthOfGameCode, len(gameCode))
	inflatedBoard, err := game.InflateNormalGameFromCode(gameCode)
	assert.NotEmpty(t, inflatedBoard)
	assert.Empty(t, err)
}

func TestGameCodeWithDelimiter(t *testing.T) {
	board := MapGenerationAttempt(game.NormalGame, false)
	gameCode := board.GetGameCode(true)
	expectedLengthOfGameCode := game.NormalGame.TilesCount*3 + len(board.Board) // landscape, number, harbor + rows
	assert.Equal(t, expectedLengthOfGameCode, len(gameCode))
	inflatedBoard, err := game.InflateNormalGameFromCode(gameCode)
	assert.NotEmpty(t, inflatedBoard)
	assert.Empty(t, err)
}

func TestGameCodeLage(t *testing.T) {
	board := MapGenerationAttempt(game.NormalGame, false)
	gameCode := board.GetGameCode(false)
	expectedLengthOfGameCode := game.NormalGame.TilesCount * 3 // landscape, number, harbor
	assert.Equal(t, expectedLengthOfGameCode, len(gameCode))
	inflatedBoard, err := game.InflateNormalGameFromCode(gameCode)
	assert.NotEmpty(t, inflatedBoard)
	assert.Empty(t, err)
}

func TestGameCodeWithDelimiterLage(t *testing.T) {
	board := MapGenerationAttempt(game.LargeGame, false)
	gameCode := board.GetGameCode(true)
	expectedLengthOfGameCode := game.LargeGame.TilesCount*3 + len(board.Board) // landscape, number, harbor + rows
	assert.Equal(t, expectedLengthOfGameCode, len(gameCode))
	inflatedBoard, err := game.InflateLargeGameFromCode(gameCode)
	assert.NotEmpty(t, inflatedBoard)
	assert.Empty(t, err)
}
