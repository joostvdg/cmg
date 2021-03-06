package game

import (
	"github.com/joostvdg/cmg/pkg/model"
)

var NormalGame = CreateNormalGame()
var LargeGame = CreateLargeGame()

type PrintBoardToConsole func(b *Board)

// GameType the information for the type of game
// Should be exhaustive and will be expanded for supporting alternative game types such as Seafarers
type GameType struct {
	Name               string
	TilesCount         int
	DesertCount        int
	ForestCount        int
	PastureCount       int
	FieldCount         int
	RiverCount         int
	MountainCount      int
	HarborCount        int
	AdjacentTileGroups [][]string
	NumberSet          []*model.Number
	HarborSet          []*model.Harbor
	HarborLayout       []string
	BoardLayout        map[string]int
	ToConsole          PrintBoardToConsole
}
