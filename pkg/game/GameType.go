package game

var NormalGame = createNormalGame()
var LargeGame = createLargeGame()

type GameType struct {
	Name          string
	TilesCount    int
	DesertCount   int
	ForestCount   int
	PastureCount  int
	FieldCount    int
	RiverCount    int
	MountainCount int
	HarborCount   int
}

func createLargeGame() GameType {
	return GameType{}
}

func createNormalGame() GameType {
	return GameType{
		Name:          "Normal",
		TilesCount:    19,
		DesertCount:   1,
		ForestCount:   4,
		PastureCount:  4,
		FieldCount:    4,
		RiverCount:    3,
		MountainCount: 3,
		HarborCount:   6,
	}
}
