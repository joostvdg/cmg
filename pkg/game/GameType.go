package game

import "github.com/joostvdg/cmg/pkg/model"

var NormalGame = createNormalGame()
var LargeGame = createLargeGame()

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
	HarborLayout       [][]string
	BoardLayout        map[string]int
}

func createLargeGame() GameType {
	return GameType{}
}

// createNormalGame creates a Normal game for up to four players.
// Will create a board layout as shown below.
// Harbors: [c0], [a0, a1], [a2], [b3, c4], [d3, c4], [e2], [e0, e1]
// 			c, a0, a1, a2, b3, d3, e2, e1,  e0
//............H...........
//........../.3\..........a- b- c0 d- e-
//....H./11\\.2//.8\.H....a- b0 c0 d0 e-
//../.6\\.3//.3\\.2//.0\..a0 b0 c1 d0 e0
//.H\.1//.6\\.4//.9\\.0/H.a0 b1 c1 d1 e0
//../.4\\.1//.9\\.5//.8\..a1 b1 c2 d1 e1
//..\.2//.4\\.3//10/\.5/..a1 b2 c2 d2 e1
//../.5\\.1//.5\\.3//.2\..a2 b2 c3 d2 e2
//.H\.3//12\\.1//11\\.2/H.a2 b3 c3 d3 e2
//......\.4//10\\.4/......a- b3 c4 d3 e-
//........H.\.5/.H........a- b- c4 d- e-
//........................
// a -> an+1, bn, bn+1
// b -> an, an-1, bn+1, cn, cn+1
// c -> cn+1, dn, dn-1
// d -> dn+1, en, en-1
// where n => 0
// where an < 3
// where bn < 4
// where cn < 5
// where dn < 4
// where en < 3
func createNormalGame() GameType {
	game := GameType{
		Name:          "Normal",
		TilesCount:    19,
		DesertCount:   1,
		ForestCount:   4,
		PastureCount:  4,
		FieldCount:    4,
		RiverCount:    3,
		MountainCount: 3,
		HarborCount:   9,
		NumberSet:     generateNumberSetNormal(19),
		HarborSet:     generateHarborSetNormal(9),
		BoardLayout:   generateNormalGameLayout(),
		HarborLayout:  generateHarborLayoutNormal(),
	}
	game.AdjacentTileGroups = [][]string{
		{"0aw", "1aw", "1bw"},
		{"0aw", "0bw", "1bw"},
		{"1aw", "0bw", "1bw"},
		{"1aw", "2aw", "1bw"},
		{"1aw", "1bw", "2bw"},
		{"2aw", "2bw", "3bw"},
		{"0bw", "0cw", "1cw"},
		{"0bw", "1bw", "0cw"},
		{"1bw", "1cw", "2cw"},
		{"1bw", "2bw", "1cw"},
		{"2bw", "2cw", "3cw"},
		{"2bw", "3bw", "2cw"},
		{"3bw", "3cw", "4cw"},
		{"0cw", "1cw", "0dw"},
		{"1cw", "2cw", "1dw"},
		{"1cw", "0dw", "1dw"},
		{"2cw", "3cw", "2dw"},
		{"2cw", "2dw", "3dw"},
		{"3cw", "4cw", "3dw"},
		{"3cw", "2dw", "3dw"},
		{"0dw", "0ew", "1ew"},
		{"0dw", "1dw", "0ew"},
		{"1dw", "2dw", "1ew"},
		{"2dw", "1ew", "2ew"},
		{"2dw", "3dw", "2ew"},
		{"3dw", "2dw", "2ew"},
		{"0ew", "1ew", "1dw"},
		{"0ew", "0dw", "1dw"},
		{"1ew", "2ew", "1dw"},
		{"2ew", "3dw", "2dw"},
	}
	return game
}
func generateNormalGameLayout() map[string]int {
	var boardLayout map[string]int
	boardLayout = make(map[string]int)
	boardLayout["a"] = 3
	boardLayout["b"] = 4
	boardLayout["c"] = 5
	boardLayout["d"] = 4
	boardLayout["e"] = 3

	return boardLayout
}

func generateHarborSetNormal(numberOfHarbors int) []*model.Harbor {
	harbors := make([]*model.Harbor, 0, numberOfHarbors)
	harbors = append(harbors, &model.Harbor{Name: "2:1 Grain", Resource: model.Grain})
	harbors = append(harbors, &model.Harbor{Name: "2:1 Brick", Resource: model.Brick})
	harbors = append(harbors, &model.Harbor{Name: "2:1 Ore", Resource: model.Ore})
	harbors = append(harbors, &model.Harbor{Name: "2:1 Wool", Resource: model.Wool})
	harbors = append(harbors, &model.Harbor{Name: "2:1 Lumber", Resource: model.Lumber})
	harbors = append(harbors, &model.Harbor{Name: "3:1", Resource: model.All})
	harbors = append(harbors, &model.Harbor{Name: "3:1", Resource: model.All})
	harbors = append(harbors, &model.Harbor{Name: "3:1", Resource: model.All})
	harbors = append(harbors, &model.Harbor{Name: "3:1", Resource: model.All})
	return harbors
}

// generateHarborPositionsNormal creates the matrix of the harbors positions
func generateHarborLayoutNormal() [][]string {
	return [][]string{
		{"c0"},
		{"a0", "b0"},
		{"a1", "a0"},
		{"a2"},
		{"b3", "c4"},
		{"d3", "c4"},
		{"e2"},
		{"e1", "e0"},
		{"e0", "d0"},
	}
}

func generateNumberSetNormal(numberOfTiles int) []*model.Number {
	numbers := make([]*model.Number, 0, numberOfTiles-1)

	numbers = append(numbers, &model.Number{Number: 2, Score: 27})
	numbers = append(numbers, &model.Number{Number: 3, Score: 55})
	numbers = append(numbers, &model.Number{Number: 3, Score: 55})
	numbers = append(numbers, &model.Number{Number: 4, Score: 83})
	numbers = append(numbers, &model.Number{Number: 4, Score: 83})
	numbers = append(numbers, &model.Number{Number: 5, Score: 111})
	numbers = append(numbers, &model.Number{Number: 5, Score: 111})
	numbers = append(numbers, &model.Number{Number: 6, Score: 139})
	numbers = append(numbers, &model.Number{Number: 6, Score: 139})
	numbers = append(numbers, &model.Number{Number: 8, Score: 139})
	numbers = append(numbers, &model.Number{Number: 8, Score: 139})
	numbers = append(numbers, &model.Number{Number: 9, Score: 111})
	numbers = append(numbers, &model.Number{Number: 9, Score: 111})
	numbers = append(numbers, &model.Number{Number: 10, Score: 83})
	numbers = append(numbers, &model.Number{Number: 10, Score: 83})
	numbers = append(numbers, &model.Number{Number: 11, Score: 55})
	numbers = append(numbers, &model.Number{Number: 11, Score: 55})
	numbers = append(numbers, &model.Number{Number: 12, Score: 27})

	return numbers
}
