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

func generateMap(gameType game.GameType, verbose bool) game.Board {

	log.Info("Generating new Map")
	tiles := generateTiles(gameType)
	distributeNumbers(gameType, tiles)
	if verbose {
		for _, tile := range tiles {
			log.WithFields(log.Fields{
				"Landscape": tile.Landscape,
				"Number":    tile.Number,
				"Harbor":    tile.Harbor,
			}).Info("Tile:")
		}
	}
	boardMap := distributeTiles(gameType, tiles, verbose)
	board := &game.Board{
		Tiles:    tiles,
		Board:    boardMap,
		GameType: gameType,
	}
	log.Info("Created a new board")
	return *board
}

func generateTiles(gameType game.GameType) []*model.Tile {
	if gameType.Name != "Normal" {
		log.Fatal("Currently not supported")
	}

	tiles := make([]*model.Tile, 0, gameType.TilesCount)
	tiles = append(tiles, addTilesOfType(gameType.DesertCount, model.Desert)...)
	tiles = append(tiles, addTilesOfType(gameType.FieldCount, model.Field)...)
	tiles = append(tiles, addTilesOfType(gameType.ForestCount, model.Forest)...)
	tiles = append(tiles, addTilesOfType(gameType.MountainCount, model.Mountain)...)
	tiles = append(tiles, addTilesOfType(gameType.PastureCount, model.Pasture)...)
	tiles = append(tiles, addTilesOfType(gameType.RiverCount, model.River)...)
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

func distributeNumbers(game game.GameType, tileSet []*model.Tile) {
	numbersAllocated := make([]int, 0, game.TilesCount-1)
	randomRange := (game.TilesCount - 1) // desert tile doesn't get a number
	log.Info("Allocating numbers to Tiles")
	for i := 0; i < game.TilesCount; i++ {
		if tileSet[i].Landscape == model.Desert {
			continue
		}
		drawnNumber := drawTileNumber(randomRange, numbersAllocated)
		number := game.NumberSet[drawnNumber]
		numbersAllocated = append(numbersAllocated, drawnNumber)
		tileSet[i].Number = *number
	}
}

func distributeTiles(gameType game.GameType, tileSet []*model.Tile, verbose bool) map[string][]*model.Tile {
	var tilesOnBoard map[string][]*model.Tile
	tilesOnBoard = make(map[string][]*model.Tile)

	randomRange := gameType.TilesCount
	numbersAllocated := make([]int, 0, gameType.TilesCount)
	for gridLane, tileInLane := range gameType.BoardLayout {
		tilesLine := make([]*model.Tile, tileInLane, tileInLane)
		for i := 0; i < tileInLane; i++ {
			drawnTileNumber := drawTileNumber(randomRange, numbersAllocated)
			tile := tileSet[drawnTileNumber]
			numbersAllocated = append(numbersAllocated, drawnTileNumber)
			tilesLine[i] = tile
		}
		tilesOnBoard[gridLane] = tilesLine
		if verbose {
			log.Info("Current Tiles Allocated: ", len(numbersAllocated), " / 19")
			log.Info("Tiles Allocated: ", numbersAllocated)
		}
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
