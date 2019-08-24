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
	AdjacentTileGroups [][]*model.TileCode
	//AdjacentTileGroups [][]string
	NumberSet []*model.Number
	HarborSet []*model.Harbor
	//HarborLayout       []string
	HarborLayout []*model.TileCode
	BoardLayout  []int
	//ToConsole          PrintBoardToConsole
}

// TC creates a TileCode via Column (named) and Row (indexed)
func TC(c int, r int) *model.TileCode {
	return &model.TileCode{Column: c, Row: r}
}
