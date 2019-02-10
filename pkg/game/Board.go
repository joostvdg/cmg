package game

import (
	"fmt"
	"github.com/joostvdg/cmg/pkg/model"
	log "github.com/sirupsen/logrus"
	"strconv"
	"strings"
)

type Board struct {
	Tiles    []*model.Tile
	Board    map[string][]*model.Tile
	GameType GameType
	Harbors  map[string]*model.Harbor
}

func (b *Board) IsValid(rules GameRules, game GameType, verbose bool) bool {
	return len(b.Tiles) == game.TilesCount &&
		b.validateAdjacentTiles(rules, verbose) &&
		b.validateResourceScores(rules, verbose) &&
		b.validateHarbors()
}

func (b *Board) validateResourceScores(rules GameRules, verbose bool) bool {
	isValid := true
	resourceCounts := make([]int, 6, 6)
	resourceScores := make([]int, 6, 6)
	resourceScores[model.Desert] = 0
	resourceScores[model.Forest] = 0
	resourceScores[model.Pasture] = 0
	resourceScores[model.Field] = 0
	resourceScores[model.River] = 0
	resourceScores[model.Mountain] = 0
	resourceCounts[model.Desert] = 0
	resourceCounts[model.Forest] = 0
	resourceCounts[model.Pasture] = 0
	resourceCounts[model.Field] = 0
	resourceCounts[model.River] = 0
	resourceCounts[model.Mountain] = 0

	for _, tile := range b.Tiles {
		switch tile.Landscape {
		case model.Forest:
			resourceScores[model.Forest] = resourceScores[model.Forest] + tile.Number.Score
			resourceCounts[model.Forest] = resourceCounts[model.Forest] + 1
		case model.Pasture:
			resourceScores[model.Pasture] = resourceScores[model.Pasture] + tile.Number.Score
			resourceCounts[model.Pasture] = resourceCounts[model.Pasture] + 1
		case model.Field:
			resourceScores[model.Field] = resourceScores[model.Field] + tile.Number.Score
			resourceCounts[model.Field] = resourceCounts[model.Field] + 1
		case model.River:
			resourceScores[model.River] = resourceScores[model.River] + tile.Number.Score
			resourceCounts[model.River] = resourceCounts[model.River] + 1
		case model.Mountain:
			resourceScores[model.Mountain] = resourceScores[model.Mountain] + tile.Number.Score
			resourceCounts[model.Mountain] = resourceCounts[model.Mountain] + 1
		}
	}

	for resourceId, score := range resourceScores {
		if resourceId == 0 {
			// skip Desert tiles
			continue
		}
		avgScore := score / resourceCounts[resourceId]
		if avgScore > rules.MaximumResourceScore || avgScore < rules.MinimumResourceScore {
			log.WithFields(log.Fields{
				"resourceId": resourceId,
				"avgScore":   avgScore,
			}).Warn("Invalid scoring for resource:")
			isValid = false
		}
	}

	return isValid
}

func (b *Board) validateAdjacentTiles(rules GameRules, verbose bool) bool {

	scoresOver300 := 0
	for _, tileGroup := range b.GameType.AdjacentTileGroups {
		valid, weightTotal := b.validateAdjectTileGroup(rules.MaximumScore, rules.MinimumScore, tileGroup[0], tileGroup[1], tileGroup[2])
		if weightTotal > 300 {
			scoresOver300++
		}
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

	if scoresOver300 > rules.MaxOver300 {
		return false
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

// validateHarbors validates whether or not a harbor is linked to a resource tile with the same resource as the harbor
func (b *Board) validateHarbors() bool {
	for k,v := range b.Harbors {
		harborResource := v.Resource
		tileCodeA := k
		tileCodeB := ""
		if strings.Contains(tileCodeA, ",") {
			tileCodes := strings.Split(tileCodeA, ",")
			tileCodeA = tileCodes[0]
			tileCodeB = tileCodes[1]
		}
		if sameResource(tileCodeA, harborResource, b.Board) || sameResource(tileCodeB, harborResource, b.Board) {
			return false
		}
	}
	return true
}

func sameResource(tileCode string, harborResource model.Resource, board map[string][]*model.Tile) bool {
	if tileCode == "" {
		return false
	}
	runeCode := []rune(tileCode)
	column := string(runeCode[0:1])
	row, _ := strconv.Atoi(string(runeCode[1:2]))
	if board[column][row].Resource == harborResource {
		return true
	}

	return false
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
	//                          a    b    c    d    e
	// a: [3], b: [4], c: [5], d: [4], e: [3]
)

func (b *Board) PrintToConsole() {

	h := b.Harbors

	// 5x10
	fmt.Printf(fmt.Sprintf(line00TemplateNormal, h["c0"].Resource))                                // 0
	fmt.Printf(fmt.Sprintf(line01TemplateNormal, b.element("0cn"))) // 1 - 0cn
	fmt.Printf(fmt.Sprintf(line02TemplateNormal,
		h["a0,b0"].Resource,
		b.element("0bn"), b.element("0cl"), b.element("0dn"),
		h["e0,d0"].Resource,)) // 2 - 0bn, 0cl, 0dn
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
		h["e1,e0"].Resource,)) // 5 - 1an, 1bl, 2cn, 1dl, 1en
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
		h["d3,c4"].Resource,)) // 10 - 4cl
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
		return fmt.Sprintf("%v", board.Board[column][row].Number.Score)
	default: // todo: panic?
		return ""
	}

}
