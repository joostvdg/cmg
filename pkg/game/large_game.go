package game

import (
	"github.com/joostvdg/cmg/pkg/model"
)

func InflateLargeGameFromCode(code string) (Board, error) {
	gameLayout := generateLargeGameLayout()
	return inflateGameFromCode(code, gameLayout)
}

func generateLargeGameLayout() []int {
	return []int{3, 4, 5, 6, 5, 4, 3}
}

func generateNumberSetLarge(numberOfTiles int) []*model.Number {
	numbers := make([]*model.Number, 0, numberOfTiles-2) // two desert tiles
	numbers = append(numbers, model.Number2)
	numbers = append(numbers, model.Number2)
	numbers = append(numbers, model.Number3)
	numbers = append(numbers, model.Number3)
	numbers = append(numbers, model.Number3)
	numbers = append(numbers, model.Number4)
	numbers = append(numbers, model.Number4)
	numbers = append(numbers, model.Number4)
	numbers = append(numbers, model.Number5)
	numbers = append(numbers, model.Number5)
	numbers = append(numbers, model.Number5)
	numbers = append(numbers, model.Number6)
	numbers = append(numbers, model.Number6)
	numbers = append(numbers, model.Number6)
	numbers = append(numbers, model.Number8)
	numbers = append(numbers, model.Number8)
	numbers = append(numbers, model.Number8)
	numbers = append(numbers, model.Number9)
	numbers = append(numbers, model.Number9)
	numbers = append(numbers, model.Number9)
	numbers = append(numbers, model.Number10)
	numbers = append(numbers, model.Number10)
	numbers = append(numbers, model.Number10)
	numbers = append(numbers, model.Number11)
	numbers = append(numbers, model.Number11)
	numbers = append(numbers, model.Number11)
	numbers = append(numbers, model.Number12)
	numbers = append(numbers, model.Number12)
	return numbers
}

func generateHarborSetLarge(numberOfHarbors int) []*model.Harbor {
	harbors := make([]*model.Harbor, 0, numberOfHarbors)
	harbors = append(harbors, model.HarborGrain)
	harbors = append(harbors, model.HarborBrick)
	harbors = append(harbors, model.HarborOre)
	harbors = append(harbors, model.HarborWool)
	harbors = append(harbors, model.HarborWool)
	harbors = append(harbors, model.HarborLumber)
	harbors = append(harbors, model.HarborAll)
	harbors = append(harbors, model.HarborAll)
	harbors = append(harbors, model.HarborAll)
	harbors = append(harbors, model.HarborAll)
	harbors = append(harbors, model.HarborAll)

	return harbors
}

func generateHarborLayoutLarge() []*model.TileCode {
	return []*model.TileCode{
		TC(3, 0),
		TC(0, 0),
		TC(0, 1),
		TC(0, 2),
		TC(2, 4),
		TC(4, 4),
		TC(5, 3),
		TC(6, 2),
		TC(6, 1),
		TC(6, 0),
		TC(4, 0),
	}
	//return []string{"d0","a0","a1","a2","c4","e4","f3","g2","g1","g0","e0"}
}

// CreateLargeGame creates a Large game for up to four players.
// Will create a board layout as shown below.
func CreateLargeGame() *GameType {
	game := GameType{
		Name:          "Large",
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
		//ToConsole:     printLargeGameToConsole,
	}
	//game.AdjacentTileGroups = [][]string{
	//	{"0aw", "1aw", "1bw"}, {"0aw", "0bw", "1bw"},
	//	{"1aw", "2aw", "1bw"}, {"1aw", "1bw", "2bw"},
	//	{"2aw", "2bw", "3bw"},
	//	{"0bw", "0cw", "1cw"}, {"0bw", "1bw", "0cw"},
	//	{"1bw", "1cw", "2cw"}, {"1bw", "2bw", "1cw"},
	//	{"2bw", "2cw", "3cw"}, {"2bw", "3bw", "2cw"},
	//	{"3bw", "3cw", "4cw"},
	//	{"0cw", "1cw", "0dw"},
	//	{"1cw", "2cw", "1dw"}, {"1cw", "0dw", "1dw"},
	//	{"2cw", "3cw", "2dw"}, {"2cw", "2dw", "3dw"},
	//	{"3cw", "4cw", "3dw"}, {"3cw", "2dw", "3dw"},
	//	{"4cw", "4dw", "5dw"},
	//	{"0dw", "0ew", "1ew"}, {"0dw", "1dw", "0ew"},
	//	{"1dw", "2dw", "1ew"},
	//	{"2dw", "1ew", "2ew"}, {"2dw", "3dw", "2ew"},
	//	{"3dw", "2ew", "3ew"}, {"3dw", "4dw", "3ew"},
	//	{"4dw", "3ew", "4ew"}, {"4dw", "5dw", "3ew"},
	//	{"0ew", "1ew", "0fw"},
	//	{"1ew", "2ew", "1fw"}, {"1ew", "0fw", "1fw"},
	//	{"2ew", "3ew", "2fw"}, {"2ew", "1fw", "2fw"},
	//	{"3ew", "4ew", "3fw"}, {"3ew", "2fw", "3fw"},
	//	{"0fw", "1fw", "0gw"},
	//	{"1fw", "2fw", "1gw"}, {"1fw", "0gw", "1gw"},
	//	{"2fw", "3fw", "2gw"}, {"2fw", "1gw", "2gw"},
	//}
	game.AdjacentTileGroups = [][]*model.TileCode{
		{TC(0, 0), TC(0, 1), TC(1, 0)},
		{TC(0, 0), TC(1, 0), TC(1, 1)},
		{TC(0, 1), TC(0, 2), TC(1, 1)},
		{TC(0, 1), TC(1, 1), TC(1, 2)},
		{TC(0, 2), TC(1, 2), TC(1, 3)},
		{TC(1, 0), TC(2, 0), TC(2, 1)},
		{TC(1, 0), TC(1, 1), TC(2, 1)},
		{TC(1, 1), TC(2, 1), TC(2, 2)},
		{TC(1, 1), TC(1, 2), TC(2, 2)},
		{TC(1, 2), TC(2, 2), TC(2, 3)},
		{TC(1, 2), TC(1, 3), TC(2, 3)},
		{TC(1, 3), TC(2, 3), TC(2, 4)},
		{TC(2, 0), TC(3, 0), TC(3, 1)},
		{TC(2, 0), TC(2, 1), TC(3, 1)},
		{TC(2, 1), TC(3, 1), TC(3, 2)},
		{TC(2, 1), TC(2, 2), TC(3, 2)},
		{TC(2, 2), TC(3, 2), TC(3, 3)},
		{TC(2, 2), TC(2, 3), TC(3, 3)},
		{TC(2, 3), TC(3, 3), TC(3, 4)},
		{TC(2, 3), TC(2, 4), TC(3, 4)},
		{TC(2, 4), TC(3, 4), TC(3, 5)},
		{TC(3, 0), TC(3, 1), TC(4, 0)},
		{TC(3, 1), TC(4, 0), TC(4, 1)},
		{TC(3, 1), TC(3, 2), TC(4, 1)},
		{TC(3, 2), TC(4, 1), TC(4, 2)},
		{TC(3, 2), TC(3, 3), TC(4, 2)},
		{TC(3, 3), TC(4, 2), TC(4, 3)},
		{TC(3, 3), TC(3, 4), TC(4, 3)},
		{TC(3, 4), TC(4, 3), TC(4, 4)},
		{TC(3, 4), TC(3, 5), TC(4, 4)},
		{TC(4, 0), TC(4, 1), TC(5, 0)},
		{TC(4, 1), TC(5, 0), TC(5, 1)},
		{TC(4, 1), TC(4, 2), TC(5, 1)},
		{TC(4, 2), TC(5, 1), TC(5, 2)},
		{TC(4, 2), TC(4, 3), TC(5, 2)},
		{TC(4, 3), TC(5, 2), TC(5, 3)},
		{TC(4, 3), TC(4, 4), TC(5, 3)},
		{TC(5, 0), TC(5, 1), TC(6, 0)},
		{TC(5, 1), TC(6, 0), TC(6, 2)},
		{TC(5, 1), TC(5, 2), TC(6, 1)},
		{TC(5, 2), TC(6, 1), TC(6, 2)},
		{TC(5, 2), TC(5, 3), TC(6, 2)},
	}

	return &game
}

//const (
//	//                          a    b    c    d    e
//	// a: [3], b: [4], c: [5], d: [4], e: [3]
//	// X = 7, Y = 6
//	// Z = (X - 1) / 2 (center)
//	// range (-Z -> Z) -> Y = (X - 1) + Z
//	// i0 = 0 + -3(Y) + 6(X-1) = 3 | 3
//	// i1 = 1 + -2(Y) + 6(X-1) = 3 | 4
//	// i2 = 2 + -1(Y) + 6(X-1) = 3 | 5
//	// i3 = 3 + -0(Y) + 6(X-1) = 3 | 6
//	// i4 = 4 +  1(Y) + 6(X-1) = 7 | 5
//	// i5 = 5 +  2(Y) + 6(X-1) = 8 | 4
//	// i6 = 6 +  3(Y) + 6(X-1) = 9 | 3
//	line00TemplateLarge string = "...................H%v..................\n"             //a- b- c- d- e- f- g-
//	line01TemplateLarge string = "................../%v\\.H%v..............\n"            //a- b- c- d0 e- f- g-
//	line02TemplateLarge string = "............../%v\\\\.%v//%v\\.............\n"          //a- b- c0 d0 e0 f- g-
//	line03TemplateLarge string = ".......H%v./%v\\\\.%v//%v\\\\%v//%v\\.H%v.......\n"     //a- b0 c0 d1 e0 f0 g-
//	line04TemplateLarge string = "....../%v\\\\.%v//%v\\\\.%v//%v\\\\.%v//%v\\.....\n"    //a0 b0 c1 d1 e1 f0 g0
//	line05TemplateLarge string = "......\\.%v//%v\\\\.%v//%v\\\\.%v//%v\\\\.%v/.....\n"   //a0 b1 c1 d2 e1 f1 g0
//	line06TemplateLarge string = "...H%v./%v\\\\.%v//%v\\\\.%v//%v\\\\.%v//%v\\.H%v..\n"  //a1 b1 c2 d2 e2 f1 g1
//	line07TemplateLarge string = "......\\.%v//%v\\\\.%v//%v\\\\.%v//%v\\\\.%v/.....\n"   //a1 b2 c2 d3 e2 f2 g1
//	line08TemplateLarge string = "....../%v\\\\.%v//%v\\\\.%v//%v\\\\.%v//%v\\.....\n"    //a2 b2 c3 d3 e3 f2 g2
//	line09TemplateLarge string = "...H%v.\\.%v//%v\\\\.%v//%v\\\\.%v//%v\\\\.%v/.H%v..\n" //a2 b3 c3 d4 e3 f3 g2
//	line10TemplateLarge string = "..........\\.%v//%v\\\\.%v//%v\\\\.%v/.........\n"      //a- b3 c4 d4 e4 f3 g-
//	line11TemplateLarge string = "..............\\.%v//%v\\\\.%v/.H%v..........\n"        //a- b- c4 d5 e4 f- g-
//	line12TemplateLarge string = "..............H%v.\\.%v/.H%v...............\n"          //a- b- c- d5 e- f- g-
//	line13TemplateLarge string = ".......................................\n"              //a- b- c- d- e- f- g-
//)

// PrintLargeGameToConsole prints the game board to the console
//................H..................
//............../10\...H.............a- b- c- d0 e- f- g-
//........../.3\\.4//.8\.............a- b- c0 d0 e0 f- g-
//....H./11\\.2//.9\\.4//.8\.H.......a- b0 c0 d1 e0 f0 g-
//../.6\\.3//.3\\.4//11\\.2//.0\.....a0 b0 c1 d1 e1 f0 g0
//..\.1//.6\\.4//.3\\.4//.9\\.0/.....a0 b1 c1 d2 e1 f1 g0
//H./.4\\.1//.9\\.3//.3\\.5//.8\.H...a1 b1 c2 d2 e2 f1 g1
//..\.2//.4\\.3//12\\.2//10/\.5/.....a1 b2 c2 d3 e2 f2 g1
//../.5\\.1//.5\\.4//.2\\.3//.2\.....a2 b2 c3 d3 e3 f2 g2
//H.\.3//12\\.1//.4\\.4//11\\.2/.H...a2 b3 c3 d4 e3 f3 g2
//......\.4//10\\.2//.5\\.4/.........a- b3 c4 d4 e4 f3 g-
//..........\.5//.5\\.6/.H...........a- b- c4 d5 e4 f- g-
//...........H..\.3/.H...............a- b- c- d5 e- f- g-
//func printLargeGameToConsole(b *Board) {
//	// 3 "a0", "g0"
//	// 6 "a1", "g1"
//	// 9 "a2", "g2"
//	// 11 "f3"
//	// 12 "c4, "e4"
//	h := b.Harbors
//	fmt.Printf(fmt.Sprintf(line00TemplateLarge,
//		h["d0"].Resource)) //
//	fmt.Printf(fmt.Sprintf(line01TemplateLarge,
//		b.element("0dn"),
//		h["e0"].Resource)) //           0dn
//	fmt.Printf(fmt.Sprintf(line02TemplateLarge,
//		b.element("0cn"),
//		b.element("0dl"),
//		b.element("0en"))) //           0cn, 0dl, 0en
//	fmt.Printf(fmt.Sprintf(line03TemplateLarge,
//		h["a0"].Resource,
//		b.element("0bn"),
//		b.element("0cl"),
//		b.element("1dn"),
//		b.element("0el"),
//		b.element("0fn"),
//		h["g0"].Resource)) //      0bn, 0cl, 1dn, 0el, 0fn,
//	fmt.Printf(fmt.Sprintf(line04TemplateLarge,
//		b.element("0an"),
//		b.element("0bl"),
//		b.element("1cn"),
//		b.element("1dl"),
//		b.element("1en"),
//		b.element("0fl"),
//		b.element("0gn"))) // 0an, 0bl, 1cn, 1dl, 1en, 0fl, 0gn
//	fmt.Printf(fmt.Sprintf(line05TemplateLarge,
//		b.element("0al"),
//		b.element("1bn"),
//		b.element("1cl"),
//		b.element("2dn"),
//		b.element("1el"),
//		b.element("1fn"),
//		b.element("0gl"))) // 0al, 1bn, 1cl, 2dn, 1el, 1fn, 0gl
//	fmt.Printf(fmt.Sprintf(line06TemplateLarge,
//		h["a1"].Resource,
//		b.element("1an"),
//		b.element("1bl"),
//		b.element("2cn"),
//		b.element("2dl"),
//		b.element("2en"),
//		b.element("1fl"),
//		b.element("1gn"),
//		h["g1"].Resource)) // 1an, 1bl, 2cn, 2dl, 2en, 1fl, 1gn
//	fmt.Printf(fmt.Sprintf(line07TemplateLarge,
//		b.element("1al"),
//		b.element("2bn"),
//		b.element("2cl"),
//		b.element("3dn"),
//		b.element("2el"),
//		b.element("2fn"),
//		b.element("1gl"))) // 1al, 2bn, 2cl, 3dn, 2el, 2fn, 1gl
//	fmt.Printf(fmt.Sprintf(line08TemplateLarge,
//		b.element("2an"),
//		b.element("2bl"),
//		b.element("3cn"),
//		b.element("3dl"),
//		b.element("3en"),
//		b.element("2fl"),
//		b.element("2gn"))) // 2an, 2bl, 3cn, 3dl, 3en, 2fl, 2gn
//	fmt.Printf(fmt.Sprintf(line09TemplateLarge,
//		h["a2"].Resource,
//		b.element("2al"),
//		b.element("3bn"),
//		b.element("3cl"),
//		b.element("4dn"),
//		b.element("3el"),
//		b.element("3fn"),
//		b.element("2gl"),
//		h["g2"].Resource)) // 2al, 3bn, 3cl, 4dn, 3el, 3fn, 2gl
//	fmt.Printf(fmt.Sprintf(line10TemplateLarge,
//		b.element("3bl"),
//		b.element("4cn"),
//		b.element("4dl"),
//		b.element("4en"),
//		b.element("3fl"))) //      3bl, 4cn, 4dl, 4en, 3fl
//	fmt.Printf(fmt.Sprintf(line11TemplateLarge,
//		b.element("4cl"),
//		b.element("5dn"),
//		b.element("4el"),
//		h["f3"].Resource)) //           4cn, 5dl, 4en
//	fmt.Printf(fmt.Sprintf(line12TemplateLarge,
//		h["c4"].Resource,
//		b.element("5dl"),
//		h["e4"].Resource)) //                5dl
//	fmt.Printf(fmt.Sprintf(line13TemplateLarge)) //
//}
