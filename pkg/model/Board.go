package model

import (
	"fmt"
	"github.com/joostvdg/cmg/pkg/game"
	log "github.com/sirupsen/logrus"
	"strconv"
)

type Board struct {
	Tiles []*Tile
	Board map[string][]*Tile
}

// TODO: validate distribution: do we have to many > 300 spots or to many < 200 spots
func (b *Board) IsValid(rules game.GameRules, game game.GameType, verbose bool) bool {
	return len(b.Tiles) == game.TilesCount && b.validateAdjacentTiles(rules, verbose)
}

//........................
//........../.3\..........a- b- c0 d- e-
//....../11\\.2//.8\......a- b0 c0 d0 e-
//../.6\\.3//.3\\.2//.0\..a0 b0 c1 d0 e0
//..\.1//.6\\.4//.9\\.0/..a0 b1 c1 d1 e0
//../.4\\.1//.9\\.5//.8\..a1 b1 c2 d1 e1
//..\.2//.4\\.3//10/\.5/..a1 b2 c2 d2 e1
//../.5\\.1//.5\\.3//.2\..a2 b2 c3 d2 e2
//..\.3//12\\.1//11\\.2/..a2 b3 c3 d3 e2
//......\.4//10\\.4/......a- b3 c4 d3 e-
//..........\.5/..........a- b- c4 d- e-
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

func (b *Board) validateAdjacentTiles(rules game.GameRules, verbose bool) bool {

	// only have to process b, c, and d, as we need collections of three
	//    / \\ / - a0 - b0 - b1
	// / \\ // \ - a1 - a0 - b1
	// \ // \/ \ -

	tileGroups := [][]string{
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

	weights := make([]int, 0, len(tileGroups))
	for _, tileGroup := range tileGroups {
		valid, weightTotal := b.validateAdjectTileGroup(rules.MaximumScore, rules.MinimumScore, tileGroup[0], tileGroup[1], tileGroup[2])
		weights = append(weights, weightTotal)
		if verbose {
			tileGroupSet := fmt.Sprintf("[%s, %s, %s]", tileGroup[0], tileGroup[1], tileGroup[2])
			log.WithFields(log.Fields{
				"tileGroup": tileGroupSet,
				"weight":    weightTotal,
			}).Info("Validating Tile Group:")
		}
		if !valid {
			return false
		}
	}
	return true

}

func (b *Board) validateAdjectTileGroup(max int, min int, tileCodeA string, tileCodeB string, tileCodeC string) (bool, int) {
	weightTileA, _ := strconv.Atoi(b.element(tileCodeA))
	weightTileB, _ := strconv.Atoi(b.element(tileCodeB))
	weightTileC, _ := strconv.Atoi(b.element(tileCodeC))
	weightTotal := weightTileA + weightTileB + weightTileC
	if weightTotal > max || weightTotal < min {
		return false, weightTotal
	}
	return true, weightTotal
}

const (
	line00Template string = "........................\n"
	line01Template string = "........../%v\\..........\n"
	line02Template string = "....../%v\\\\.%v//%v\\......\n"
	line03Template string = "../%v\\\\.%v//%v\\\\.%v//%v\\..\n"
	line04Template string = "..\\.%v//%v\\\\.%v//%v\\\\.%v/..\n"
	line05Template string = "../%v\\\\.%v//%v\\\\.%v//%v\\..\n"
	line06Template string = "..\\.%v//%v\\\\.%v//%v/\\.%v/..\n"
	line07Template string = "../%v\\\\.%v//%v\\\\.%v//%v\\..\n"
	line08Template string = "..\\.%v//%v\\\\.%v//%v\\\\.%v/..\n"
	line09Template string = "......\\.%v//%v\\\\.%v/......\n"
	line10Template string = "..........\\.%v/..........\n"
	line11Template string = "........................\n"
	//                          a    b    c    d    e
	// a: [3], b: [4], c: [5], d: [4], e: [3]
)

func (b *Board) PrintToConsole() {

	// 5x10
	fmt.Printf(line00Template)                                // 0
	fmt.Printf(fmt.Sprintf(line01Template, b.element("0cn"))) // 1 - 0cn
	fmt.Printf(fmt.Sprintf(line02Template,
		b.element("0bn"), b.element("0cl"), b.element("0dn"))) // 2 - 0bn, 0cl, 0dn
	fmt.Printf(fmt.Sprintf(line03Template,
		b.element("0an"),
		b.element("0bl"),
		b.element("1cn"),
		b.element("0dl"),
		b.element("0en"))) // 3 - 0an, 0bl, 1cn, 0dl, 0en
	fmt.Printf(fmt.Sprintf(line04Template,
		b.element("0al"),
		b.element("1bn"),
		b.element("1cl"),
		b.element("1dn"),
		b.element("0el"))) // 4 - 0al, 1bn, 1cl, 1dn, 0el
	fmt.Printf(fmt.Sprintf(line05Template,
		b.element("1an"),
		b.element("1bl"),
		b.element("2cn"),
		b.element("1dl"),
		b.element("1en"))) // 5 - 1an, 1bl, 2cn, 1dl, 1en
	fmt.Printf(fmt.Sprintf(line06Template,
		b.element("1al"),
		b.element("2bn"),
		b.element("2cl"),
		b.element("2dn"),
		b.element("1el"))) // 6 - 1al, 2bn, 2cl, 2dn, 1el
	fmt.Printf(fmt.Sprintf(line07Template,
		b.element("2an"),
		b.element("2bl"),
		b.element("3cn"),
		b.element("2dl"),
		b.element("2en"))) // 7 - 2an, 2bl, 3cn, 2dl, 2en
	fmt.Printf(fmt.Sprintf(line08Template,
		b.element("2al"),
		b.element("3bn"),
		b.element("3cl"),
		b.element("3dn"),
		b.element("2el"))) // 8 - 2al, 3bn, 3cl, 3dn, 2el
	fmt.Printf(fmt.Sprintf(line09Template,
		b.element("3bl"),
		b.element("4cn"),
		b.element("3dl"))) // 9 - 3bl, 4cn, 3dl
	fmt.Printf(fmt.Sprintf(line10Template, b.element("4cl"))) // 10 - 4cl
	fmt.Printf(line11Template)                                // 11
}

func (board *Board) element(code string) string {
	runeCode := []rune(code)
	row, _ := strconv.Atoi(string(runeCode[0:1]))
	column := string(runeCode[1:2])
	elementType := string(runeCode[2:3])
	switch elementType {
	case "l":
		return fmt.Sprintf("%v", board.Board[column][row].Landscape)
	case "n":
		number := board.Board[column][row].Number.Number
		padding := ""
		if number < 10 {
			padding = "."
		}
		element := fmt.Sprintf("%s%v", padding, number)
		return element
	case "w": // weight of the number
		return fmt.Sprintf("%v", board.Board[column][row].Number.Weight)
	default: // todo: panic?
		return ""
	}

}
