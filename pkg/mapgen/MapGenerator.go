package mapgen

import (
	"fmt"
	"github.com/joostvdg/cmg/pkg/game"
	"github.com/joostvdg/cmg/pkg/model"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"time"
)

type Game int

const (
	Normal    Game = 0
	FiveOrSix Game = 1
)

const (
	NormalTilesCount    = 19
	NormalDesertCount   = 1
	NormalForestCount   = 4
	NormalPastureCount  = 4
	NormalFieldCount    = 4
	NormalRiverCount    = 3
	NormalMountainCount = 3
)

func GenerateMap(count int, loop bool, verbose bool, rules game.GameRules) {

	numberOfLoops := count
	if !loop {
		numberOfLoops = 1
	}

	var gameType game.GameType
	if rules.GameType == 0 {
		gameType = game.NormalGame
	} else if rules.GameType == 1 {
		gameType = game.LargeGame
	}

	failedGenerations := 0
	totalGenerations := 0
	board := generateMap(gameType, verbose)
	for i := 0; i < numberOfLoops; i++ {
		totalGenerations++
		for !board.IsValid(rules, gameType, verbose) {
			if totalGenerations > 1001 {
				log.Fatal("Can not generate a map... (1000+ runs)")
			}
			log.Info(fmt.Sprintf("Loop %v::%v", i, failedGenerations))
			totalGenerations++
			failedGenerations++
			board = generateMap(gameType, verbose)
		}
		board.PrintToConsole()
	}
	log.WithFields(log.Fields{
		"Map Generation Loops":    numberOfLoops,
		"Map Generation Failures": failedGenerations,
	}).Info("Finished generation loop:")
}

func generateMap(game game.GameType, verbose bool) model.Board {

	log.Info("Generating new Map")
	tiles := generateTiles(Normal)
	numbers := generateNumberSet()
	distributeNumbersNormalGame(tiles, numbers)
	if verbose {
		for _, tile := range tiles {
			log.WithFields(log.Fields{
				"Landscape": tile.Landscape,
				"Number":    tile.Number,
				"Harbor":    tile.Harbor,
			}).Info("Tile:")
		}
	}
	boardMap := distributeTilesNormalGame(tiles, verbose)
	board := &model.Board{
		Tiles: tiles,
		Board: boardMap,
	}
	log.Info("Created a new board")
	return *board
}

func generateTiles(game Game) []*model.Tile {
	if game != Normal {
		log.Fatal("Currently not supported")
	}

	tiles := make([]*model.Tile, 0, NormalTilesCount)
	tiles = append(tiles, addTilesOfType(NormalDesertCount, model.Desert)...)
	tiles = append(tiles, addTilesOfType(NormalFieldCount, model.Field)...)
	tiles = append(tiles, addTilesOfType(NormalForestCount, model.Forest)...)
	tiles = append(tiles, addTilesOfType(NormalMountainCount, model.Mountain)...)
	tiles = append(tiles, addTilesOfType(NormalPastureCount, model.Pasture)...)
	tiles = append(tiles, addTilesOfType(NormalRiverCount, model.River)...)
	return tiles
}

func addTilesOfType(number int, landscape model.LandscapeCode) []*model.Tile {
	tiles := make([]*model.Tile, number, number)
	for i := 0; i < number; i++ {
		tile := model.Tile{
			Landscape: landscape,
		}
		tiles[i] = &tile
	}
	return tiles
}

func generateNumberSet() []*model.Number {
	numbers := make([]*model.Number, 0, NormalTilesCount-1)

	numbers = append(numbers, &model.Number{Number: 2, Weight: 27})
	numbers = append(numbers, &model.Number{Number: 3, Weight: 55})
	numbers = append(numbers, &model.Number{Number: 3, Weight: 55})
	numbers = append(numbers, &model.Number{Number: 4, Weight: 83})
	numbers = append(numbers, &model.Number{Number: 4, Weight: 83})
	numbers = append(numbers, &model.Number{Number: 5, Weight: 111})
	numbers = append(numbers, &model.Number{Number: 5, Weight: 111})
	numbers = append(numbers, &model.Number{Number: 6, Weight: 139})
	numbers = append(numbers, &model.Number{Number: 6, Weight: 139})
	numbers = append(numbers, &model.Number{Number: 8, Weight: 139})
	numbers = append(numbers, &model.Number{Number: 8, Weight: 139})
	numbers = append(numbers, &model.Number{Number: 9, Weight: 111})
	numbers = append(numbers, &model.Number{Number: 9, Weight: 111})
	numbers = append(numbers, &model.Number{Number: 10, Weight: 83})
	numbers = append(numbers, &model.Number{Number: 10, Weight: 83})
	numbers = append(numbers, &model.Number{Number: 11, Weight: 55})
	numbers = append(numbers, &model.Number{Number: 11, Weight: 55})
	numbers = append(numbers, &model.Number{Number: 12, Weight: 27})

	return numbers
}

func distributeNumbersNormalGame(tileSet []*model.Tile, numbers []*model.Number) {
	numbersAllocated := make([]int, 0, NormalTilesCount-1)
	randomRange := (NormalTilesCount - 1) // desert tile doesn't get a number
	log.Info("Allocating numbers to Tiles")
	for i := 0; i < NormalTilesCount; i++ {
		if tileSet[i].Landscape == model.Desert {
			continue
		}
		drawnNumber := drawTileNumber(randomRange, numbersAllocated)
		number := numbers[drawnNumber]
		numbersAllocated = append(numbersAllocated, drawnNumber)
		tileSet[i].Number = *number
	}
}

func distributeTilesNormalGame(tileSet []*model.Tile, verbose bool) map[string][]*model.Tile {
	// a: [3], b: [4], c: [5], d: [4], e: [3]
	var tilesOnBoard map[string][]*model.Tile
	tilesOnBoard = make(map[string][]*model.Tile)

	// TODO: do something with game type
	randomRange := NormalTilesCount
	numbersAllocated := make([]int, 0, NormalTilesCount)

	// Fill line A
	tilesLineA := make([]*model.Tile, 3, 3)
	for i := 0; i < 3; i++ {
		drawnTileNumber := drawTileNumber(randomRange, numbersAllocated)
		tile := tileSet[drawnTileNumber]
		numbersAllocated = append(numbersAllocated, drawnTileNumber)
		tilesLineA[i] = tile
	}
	tilesOnBoard["a"] = tilesLineA
	if verbose {
		log.Info("Current Tiles Allocated: ", len(numbersAllocated), " / 19")
		log.Info("Tiles Allocated: ", numbersAllocated)
	}

	// Fill line B
	tilesLineB := make([]*model.Tile, 4, 4)
	for i := 0; i < 4; i++ {
		drawnTileNumber := drawTileNumber(randomRange, numbersAllocated)
		tile := tileSet[drawnTileNumber]
		numbersAllocated = append(numbersAllocated, drawnTileNumber)
		tilesLineB[i] = tile
	}
	tilesOnBoard["b"] = tilesLineB
	if verbose {
		log.Info("Current Tiles Allocated: ", len(numbersAllocated), " / 19")
		log.Info("Tiles Allocated: ", numbersAllocated)
	}

	// Fill line C
	tilesLineC := make([]*model.Tile, 5, 5)
	for i := 0; i < 5; i++ {
		drawnTileNumber := drawTileNumber(randomRange, numbersAllocated)
		tile := tileSet[drawnTileNumber]
		numbersAllocated = append(numbersAllocated, drawnTileNumber)
		tilesLineC[i] = tile
	}
	tilesOnBoard["c"] = tilesLineC
	if verbose {
		log.Info("Current Tiles Allocated: ", len(numbersAllocated), " / 19")
		log.Info("Tiles Allocated: ", numbersAllocated)
	}

	// Fill line D
	tilesLineD := make([]*model.Tile, 4, 4)
	for i := 0; i < 4; i++ {
		drawnTileNumber := drawTileNumber(randomRange, numbersAllocated)
		tile := tileSet[drawnTileNumber]
		numbersAllocated = append(numbersAllocated, drawnTileNumber)
		tilesLineD[i] = tile
	}
	tilesOnBoard["d"] = tilesLineD
	if verbose {
		log.Info("Current Tiles Allocated: ", len(numbersAllocated), " / 19")
		log.Info("Tiles Allocated: ", numbersAllocated)
	}

	// Fill line E
	tilesLineE := make([]*model.Tile, 3, 3)
	for i := 0; i < 3; i++ {
		drawnTileNumber := drawTileNumber(randomRange, numbersAllocated)
		tile := tileSet[drawnTileNumber]
		numbersAllocated = append(numbersAllocated, drawnTileNumber)
		tilesLineE[i] = tile
	}
	tilesOnBoard["e"] = tilesLineE
	if verbose {
		log.Info("Current Tiles Allocated: ", len(numbersAllocated), " / 19")
		log.Info("Tiles Allocated: ", numbersAllocated)
	}

	return tilesOnBoard
}

func drawTileNumber(randomRange int, numbersAllocated []int) int {
	rand.Seed(time.Now().UnixNano())
	number := rand.Intn(randomRange)
	for numberIsAllocated(number, numbersAllocated) {
		number = rand.Intn(randomRange)
	}
	return number
}

func numberIsAllocated(number int, numbersAllocated []int) bool {
	isAllocated := false
	for _, numberAllocated := range numbersAllocated {
		if numberAllocated == number {
			isAllocated = true
		}
	}
	return isAllocated
}
