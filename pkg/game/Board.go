package game

import (
	"fmt"
	"github.com/joostvdg/cmg/pkg/model"
	log "github.com/sirupsen/logrus"
	"strconv"
)

type Board struct {
	Tiles    []*model.Tile
	Board    map[string][]*model.Tile
	GameType GameType
}

// TODO: validate distribution: do we have to many > 300 spots or to many < 200 spots
func (b *Board) IsValid(rules GameRules, game GameType, verbose bool) bool {
	return len(b.Tiles) == game.TilesCount && b.validateAdjacentTiles(rules, verbose)
}

func (b *Board) validateAdjacentTiles(rules GameRules, verbose bool) bool {

	weights := make([]int, 0, len(b.GameType.AdjacentTileGroups))
	for _, tileGroup := range b.GameType.AdjacentTileGroups {
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
		log.WithFields(log.Fields{
			"Score":       weightTotal,
			"Max allowed": max,
			"Min allowed": min,
		}).Warn("Invalid tile group")
		return false, weightTotal
	}
	return true, weightTotal
}

const (
	line00TemplateNormal string = "........................\n"
	line01TemplateNormal string = "........../%v\\..........\n"
	line02TemplateNormal string = "....../%v\\\\.%v//%v\\......\n"
	line03TemplateNormal string = "../%v\\\\.%v//%v\\\\.%v//%v\\..\n"
	line04TemplateNormal string = "..\\.%v//%v\\\\.%v//%v\\\\.%v/..\n"
	line05TemplateNormal string = "../%v\\\\.%v//%v\\\\.%v//%v\\..\n"
	line06TemplateNormal string = "..\\.%v//%v\\\\.%v//%v/\\.%v/..\n"
	line07TemplateNormal string = "../%v\\\\.%v//%v\\\\.%v//%v\\..\n"
	line08TemplateNormal string = "..\\.%v//%v\\\\.%v//%v\\\\.%v/..\n"
	line09TemplateNormal string = "......\\.%v//%v\\\\.%v/......\n"
	line10TemplateNormal string = "..........\\.%v/..........\n"
	line11TemplateNormal string = "........................\n"
	//                          a    b    c    d    e
	// a: [3], b: [4], c: [5], d: [4], e: [3]
)

func (b *Board) PrintToConsole() {

	// 5x10
	fmt.Printf(line00TemplateNormal)                                // 0
	fmt.Printf(fmt.Sprintf(line01TemplateNormal, b.element("0cn"))) // 1 - 0cn
	fmt.Printf(fmt.Sprintf(line02TemplateNormal,
		b.element("0bn"), b.element("0cl"), b.element("0dn"))) // 2 - 0bn, 0cl, 0dn
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
		b.element("1an"),
		b.element("1bl"),
		b.element("2cn"),
		b.element("1dl"),
		b.element("1en"))) // 5 - 1an, 1bl, 2cn, 1dl, 1en
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
		b.element("2al"),
		b.element("3bn"),
		b.element("3cl"),
		b.element("3dn"),
		b.element("2el"))) // 8 - 2al, 3bn, 3cl, 3dn, 2el
	fmt.Printf(fmt.Sprintf(line09TemplateNormal,
		b.element("3bl"),
		b.element("4cn"),
		b.element("3dl"))) // 9 - 3bl, 4cn, 3dl
	fmt.Printf(fmt.Sprintf(line10TemplateNormal, b.element("4cl"))) // 10 - 4cl
	fmt.Printf(line11TemplateNormal)                                // 11
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
