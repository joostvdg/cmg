package game

import (
	"fmt"
	"github.com/joostvdg/cmg/pkg/model"
)

var NormalGame = createNormalGame()
var LargeGame = createLargeGame()

type PrintBoardToConsole func(b *Board)

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
	ToConsole          PrintBoardToConsole
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
		ToConsole:     PrintNormalGameToConsole,
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

func generateLargeGameLayout() map[string]int {
	var boardLayout map[string]int
	boardLayout = make(map[string]int)
	boardLayout["a"] = 3
	boardLayout["b"] = 4
	boardLayout["c"] = 5
	boardLayout["d"] = 6
	boardLayout["e"] = 5
	boardLayout["f"] = 4
	boardLayout["g"] = 3

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

func generateNumberSetLarge(numberOfTiles int) []*model.Number {
	numbers := make([]*model.Number, 0, numberOfTiles-2) // two desert tiles

	numbers = append(numbers, &model.Number{Number: 2, Score: 27})
	numbers = append(numbers, &model.Number{Number: 2, Score: 27})
	numbers = append(numbers, &model.Number{Number: 3, Score: 55})
	numbers = append(numbers, &model.Number{Number: 3, Score: 55})
	numbers = append(numbers, &model.Number{Number: 3, Score: 55})
	numbers = append(numbers, &model.Number{Number: 4, Score: 83})
	numbers = append(numbers, &model.Number{Number: 4, Score: 83})
	numbers = append(numbers, &model.Number{Number: 4, Score: 83})
	numbers = append(numbers, &model.Number{Number: 5, Score: 111})
	numbers = append(numbers, &model.Number{Number: 5, Score: 111})
	numbers = append(numbers, &model.Number{Number: 5, Score: 111})
	numbers = append(numbers, &model.Number{Number: 6, Score: 139})
	numbers = append(numbers, &model.Number{Number: 6, Score: 139})
	numbers = append(numbers, &model.Number{Number: 6, Score: 139})
	numbers = append(numbers, &model.Number{Number: 8, Score: 139})
	numbers = append(numbers, &model.Number{Number: 8, Score: 139})
	numbers = append(numbers, &model.Number{Number: 8, Score: 139})
	numbers = append(numbers, &model.Number{Number: 9, Score: 111})
	numbers = append(numbers, &model.Number{Number: 9, Score: 111})
	numbers = append(numbers, &model.Number{Number: 9, Score: 111})
	numbers = append(numbers, &model.Number{Number: 10, Score: 83})
	numbers = append(numbers, &model.Number{Number: 10, Score: 83})
	numbers = append(numbers, &model.Number{Number: 10, Score: 83})
	numbers = append(numbers, &model.Number{Number: 11, Score: 55})
	numbers = append(numbers, &model.Number{Number: 11, Score: 55})
	numbers = append(numbers, &model.Number{Number: 11, Score: 55})
	numbers = append(numbers, &model.Number{Number: 12, Score: 27})
	numbers = append(numbers, &model.Number{Number: 12, Score: 27})

	return numbers
}

func generateHarborSetLarge(numberOfHarbors int) []*model.Harbor {
	harbors := make([]*model.Harbor, 0, numberOfHarbors)
	harbors = append(harbors, &model.Harbor{Name: "2:1 Grain", Resource: model.Grain})
	harbors = append(harbors, &model.Harbor{Name: "2:1 Brick", Resource: model.Brick})
	harbors = append(harbors, &model.Harbor{Name: "2:1 Ore", Resource: model.Ore})
	harbors = append(harbors, &model.Harbor{Name: "2:1 Wool", Resource: model.Wool})
	harbors = append(harbors, &model.Harbor{Name: "2:1 Wool", Resource: model.Wool})
	harbors = append(harbors, &model.Harbor{Name: "2:1 Lumber", Resource: model.Lumber})
	harbors = append(harbors, &model.Harbor{Name: "3:1", Resource: model.All})
	harbors = append(harbors, &model.Harbor{Name: "3:1", Resource: model.All})
	harbors = append(harbors, &model.Harbor{Name: "3:1", Resource: model.All})
	harbors = append(harbors, &model.Harbor{Name: "3:1", Resource: model.All})
	harbors = append(harbors, &model.Harbor{Name: "3:1", Resource: model.All})
	return harbors
}

func generateHarborLayoutLarge() [][]string {
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

// createLargeGame creates a Large game for up to four players.
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

//............../10\.................a- b- c- d0 e- f- g-
//........../.3\\.4//.8\.............a- b- c0 d0 e0 f- g-
//....../11\\.2//.9\\.4//.8\.........a- b0 c0 d1 e0 f0 g-
//../.6\\.3//.3\\.4//11\\.2//.0\.....a0 b0 c1 d1 e1 f0 g0
//..\.1//.6\\.4//.3\\.4//.9\\.0/.....a0 b1 c1 d2 e1 f1 g0
//../.4\\.1//.9\\.3//.3\\.5//.8\.....a1 b1 c2 d2 e2 f1 g1
//..\.2//.4\\.3//12\\.2//10/\.5/.....a1 b2 c2 d3 e2 f2 g1
//../.5\\.1//.5\\.4//.2\\.3//.2\.....a2 b2 c3 d3 e3 f2 g2
//..\.3//12\\.1//.4\\.4//11\\.2/.....a2 b3 c3 d4 e3 f3 g2
//......\.4//10\\.2//.5\\.4/.........a- b3 c4 d4 e4 f3 g-
//..........\.5//.5\\.6/.............a- b- c4 d5 e4 f- g-
//..............\.3/.................a- b- c- d5 e- f- g-

func createLargeGame() GameType {
	game := GameType{
		Name:          "Normal",
		TilesCount:    30,
		DesertCount:   2,
		ForestCount:   6,
		PastureCount:  6,
		FieldCount:    6,
		RiverCount:    5,
		MountainCount: 5,
		HarborCount:   11,
		NumberSet:     generateNumberSetLarge(30),
		HarborSet:     generateHarborSetLarge(11),
		BoardLayout:   generateLargeGameLayout(),
		HarborLayout:  generateHarborLayoutLarge(),
		ToConsole:     PrintLargeGameToConsole,
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

const (
	line00TemplateNormal string = "............H%v...........\n"
	line01TemplateNormal string = ".........../%v\\...........\n"
	line02TemplateNormal string = ".....H%v/%v\\\\.%v//%v\\H%v.....\n"
	line03TemplateNormal string = ".../%v\\\\.%v//%v\\\\.%v//%v\\...\n"
	line04TemplateNormal string = "...\\.%v//%v\\\\.%v//%v\\\\.%v/...\n"
	line05TemplateNormal string = ".H%v/%v\\\\.%v//%v\\\\.%v//%v\\H%v.\n"
	line06TemplateNormal string = "...\\.%v//%v\\\\.%v//%v/\\.%v/...\n"
	line07TemplateNormal string = ".../%v\\\\.%v//%v\\\\.%v//%v\\...\n"
	line08TemplateNormal string = ".H%v\\.%v//%v\\\\.%v//%v\\\\.%v/H%v.\n"
	line09TemplateNormal string = ".......\\.%v//%v\\\\.%v/.......\n"
	line10TemplateNormal string = "........H%v.\\.%v/.H%v........\n"
	line11TemplateNormal string = "..........................\n"
)

func PrintNormalGameToConsole(b *Board) {

	h := b.Harbors

	// 5x10
	fmt.Printf(fmt.Sprintf(line00TemplateNormal, h["c0"].Resource)) // 0
	fmt.Printf(fmt.Sprintf(line01TemplateNormal, b.element("0cn"))) // 1 - 0cn
	fmt.Printf(fmt.Sprintf(line02TemplateNormal,
		h["a0,b0"].Resource,
		b.element("0bn"), b.element("0cl"), b.element("0dn"),
		h["e0,d0"].Resource)) // 2 - 0bn, 0cl, 0dn
	fmt.Printf(fmt.Sprintf(line03TemplateNormal,
		b.element("0an"),
		b.element("0bl"),
		b.element("1cn"),
		b.element("0dl"),
		b.element("0en"))) // 3 - 0an, 0bl, 1cn, 0dl, 0en
	fmt.Printf(fmt.Sprintf(line04TemplateNormal,
		b.element("0al"),
		b.element("1bn"),
		b.element("1cl"),
		b.element("1dn"),
		b.element("0el"))) // 4 - 0al, 1bn, 1cl, 1dn, 0el
	fmt.Printf(fmt.Sprintf(line05TemplateNormal,
		h["a1,a0"].Resource,
		b.element("1an"),
		b.element("1bl"),
		b.element("2cn"),
		b.element("1dl"),
		b.element("1en"),
		h["e1,e0"].Resource)) // 5 - 1an, 1bl, 2cn, 1dl, 1en
	fmt.Printf(fmt.Sprintf(line06TemplateNormal,
		b.element("1al"),
		b.element("2bn"),
		b.element("2cl"),
		b.element("2dn"),
		b.element("1el"))) // 6 - 1al, 2bn, 2cl, 2dn, 1el
	fmt.Printf(fmt.Sprintf(line07TemplateNormal,
		b.element("2an"),
		b.element("2bl"),
		b.element("3cn"),
		b.element("2dl"),
		b.element("2en"))) // 7 - 2an, 2bl, 3cn, 2dl, 2en
	fmt.Printf(fmt.Sprintf(line08TemplateNormal,
		h["a2"].Resource,
		b.element("2al"),
		b.element("3bn"),
		b.element("3cl"),
		b.element("3dn"),
		b.element("2el"),
		h["e2"].Resource)) // 8 - 2al, 3bn, 3cl, 3dn, 2el
	fmt.Printf(fmt.Sprintf(line09TemplateNormal,
		b.element("3bl"),
		b.element("4cn"),
		b.element("3dl"))) // 9 - 3bl, 4cn, 3dl
	fmt.Printf(fmt.Sprintf(line10TemplateNormal,
		h["b3,c4"].Resource,
		b.element("4cl"),
		h["d3,c4"].Resource)) // 10 - 4cl
	fmt.Printf(line11TemplateNormal) // 11
}

const (
	//                          a    b    c    d    e
	// a: [3], b: [4], c: [5], d: [4], e: [3]

	line00TemplateLarge string = "...................................\n"            //a- b- c- d0 e- f- g-
	line01TemplateLarge string = "............../%v\\.................\n"           //a- b- c- d0 e- f- g-
	line02TemplateLarge string = "........../%v\\\\.%v//%v\\.............\n"        //a- b- c0 d0 e0 f- g-
	line03TemplateLarge string = "....../%v\\\\.%v//%v\\\\%v//.%v\\..........\n"    //a- b0 c0 d1 e0 f0 g-
	line04TemplateLarge string = "../%v\\\\.%v//%v\\\\.%v//%v\\\\.%v//%v\\.....\n"  //a0 b0 c1 d1 e1 f0 g0
	line05TemplateLarge string = "..\\.%v//%v\\\\.%v//%v\\\\.%v//%v\\\\.%v/.....\n" //a0 b1 c1 d2 e1 f1 g0
	line06TemplateLarge string = "../%v\\\\.%v//%v\\\\.%v//%v\\\\.%v//%v\\.....\n"  //a1 b1 c2 d2 e2 f1 g1
	line07TemplateLarge string = "..\\.%v//%v\\\\.%v//%v\\\\.%v//%v\\\\.%v/.....\n" //a1 b2 c2 d3 e2 f2 g1
	line08TemplateLarge string = "../%v\\\\.%v//%v\\\\.%v//%v\\\\.%v//%v\\.....\n"  //a2 b2 c3 d3 e3 f2 g2
	line09TemplateLarge string = "..\\.%v//%v\\\\.%v//%v\\\\.%v//%v\\\\.%v/.....\n" //a2 b3 c3 d4 e3 f3 g2
	line10TemplateLarge string = "......\\.%v//%v\\\\.%v//%v\\\\.%v/.........\n"    //a- b3 c4 d4 e4 f3 g-
	line11TemplateLarge string = "..........\\.%v//%v\\\\.%v/............\n"        //a- b- c4 d5 e4 f- g-
	line12TemplateLarge string = "..............\\.%v/.................\n"          //a- b- c- d5 e- f- g-
	line13TemplateLarge string = "...................................\n"            //a- b- c- d- e- f- g-
)

//............../10\.................a- b- c- d0 e- f- g-
//........../.3\\.4//.8\.............a- b- c0 d0 e0 f- g-
//....../11\\.2//.9\\.4//.8\.........a- b0 c0 d1 e0 f0 g-
//../.6\\.3//.3\\.4//11\\.2//.0\.....a0 b0 c1 d1 e1 f0 g0
//..\.1//.6\\.4//.3\\.4//.9\\.0/.....a0 b1 c1 d2 e1 f1 g0
//../.4\\.1//.9\\.3//.3\\.5//.8\.....a1 b1 c2 d2 e2 f1 g1
//..\.2//.4\\.3//12\\.2//10/\.5/.....a1 b2 c2 d3 e2 f2 g1
//../.5\\.1//.5\\.4//.2\\.3//.2\.....a2 b2 c3 d3 e3 f2 g2
//..\.3//12\\.1//.4\\.4//11\\.2/.....a2 b3 c3 d4 e3 f3 g2
//......\.4//10\\.2//.5\\.4/.........a- b3 c4 d4 e4 f3 g-
//..........\.5//.5\\.6/.............a- b- c4 d5 e4 f- g-
//..............\.3/.................a- b- c- d5 e- f- g-
func PrintLargeGameToConsole(b *Board) {
	fmt.Printf(fmt.Sprintf(line00TemplateLarge)) //
	fmt.Printf(fmt.Sprintf(line01TemplateLarge,
		b.element("0dn"))) //           0dn
	fmt.Printf(fmt.Sprintf(line02TemplateLarge,
		b.element("0cn"),
		b.element("0dl"),
		b.element("0en"))) //           0cn, 0dl, 0en
	fmt.Printf(fmt.Sprintf(line03TemplateLarge,
		b.element("0bn"),
		b.element("0cl"),
		b.element("1dl"),
		b.element("0el"),
		b.element("0fn"))) //      0bn, 0cl, 1dn, 0el, 0fn,
	fmt.Printf(fmt.Sprintf(line04TemplateLarge,
		b.element("0an"),
		b.element("0bl"),
		b.element("1cn"),
		b.element("1dl"),
		b.element("1en"),
		b.element("0fl"),
		b.element("0gn"))) // 0an, 0bl, 1cn, 1dl, 1en, 0fl, 0gn
	fmt.Printf(fmt.Sprintf(line05TemplateLarge,
		b.element("0al"),
		b.element("1bn"),
		b.element("1cl"),
		b.element("2dn"),
		b.element("1el"),
		b.element("1fn"),
		b.element("0gl"))) // 0al, 1bn, 1cl, 2dn, 1el, 1fn, 0gl
	fmt.Printf(fmt.Sprintf(line06TemplateLarge,
		b.element("1an"),
		b.element("1bl"),
		b.element("2cn"),
		b.element("2dl"),
		b.element("2en"),
		b.element("1fl"),
		b.element("1gn"))) // 1an, 1bl, 2cn, 2dl, 2en, 1fl, 1gn
	fmt.Printf(fmt.Sprintf(line07TemplateLarge,
		b.element("1al"),
		b.element("2bn"),
		b.element("2cl"),
		b.element("3dn"),
		b.element("2el"),
		b.element("2fn"),
		b.element("1gl"))) // 1al, 2bn, 2cl, 3dn, 2el, 2fn, 1gl
	fmt.Printf(fmt.Sprintf(line08TemplateLarge,
		b.element("2an"),
		b.element("2bl"),
		b.element("3cn"),
		b.element("3dl"),
		b.element("3en"),
		b.element("2fl"),
		b.element("2gn"))) // 2an, 2bl, 3cn, 3dl, 3en, 2fl, 2gn
	fmt.Printf(fmt.Sprintf(line09TemplateLarge,
		b.element("2al"),
		b.element("3bn"),
		b.element("3cl"),
		b.element("4dn"),
		b.element("3el"),
		b.element("3fn"),
		b.element("2gl"))) // 2al, 3bn, 3cl, 4dn, 3el, 3fn, 2gl
	fmt.Printf(fmt.Sprintf(line10TemplateLarge,
		b.element("3bl"),
		b.element("4cn"),
		b.element("4dl"),
		b.element("4en"),
		b.element("3fl"))) //      3bl, 4cn, 4dl, 4en, 3fl
	fmt.Printf(fmt.Sprintf(line11TemplateLarge,
		b.element("4cn"),
		b.element("5dl"),
		b.element("4en"))) //           4cn, 5dl, 4en
	fmt.Printf(fmt.Sprintf(line12TemplateLarge,
		b.element("5dl"))) //                5dl
	fmt.Printf(fmt.Sprintf(line13TemplateLarge)) //
}
