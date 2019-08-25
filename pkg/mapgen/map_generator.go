package mapgen

import (
	"fmt"
	"github.com/getsentry/sentry-go"
	"github.com/joostvdg/cmg/pkg/game"
	"github.com/joostvdg/cmg/pkg/model"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"strconv"
	"time"
)

type Game int

func GenerateMap(count int, loop bool, verbose bool, rules game.GameRules) {

	maxGenerationAttempts := 1500
	numberOfLoops := count
	if !loop {
		numberOfLoops = 1
	}

	var gameType game.GameType
	if rules.GameType == 0 {
		gameType = game.NormalGame
	} else if rules.GameType == 1 {
		gameType = game.LargeGame
		maxGenerationAttempts = 5000 // it's more difficult
	}

	failedGenerations := 0
	totalGenerations := 0
	board := MapGenerationAttempt(gameType, verbose)
	for i := 0; i < numberOfLoops; i++ {
		totalGenerations++
		for !board.IsValid(rules, gameType) {
			if totalGenerations > maxGenerationAttempts {
				sentry.CaptureMessage(fmt.Sprintf("Could not generate map of type %v, tried %v times", gameType, totalGenerations))
				sentry.Flush(time.Second * 5)
				log.Fatal("Can not generate a map... (1000+ runs)")
			}
			log.Debug(fmt.Sprintf("Loop %v::%v", i, failedGenerations))
			totalGenerations++
			failedGenerations++
			board = MapGenerationAttempt(gameType, verbose)
		}
		board.PrintToConsole()
	}
	log.WithFields(log.Fields{
		"Map Generation Loops":    numberOfLoops,
		"Map Generation Failures": failedGenerations,
	}).Debug("Finished generation loop:")
}

func debugLogDuration(start time.Time, logMessage string) {
	if log.IsLevelEnabled(log.DebugLevel) {
		t := time.Now()
		elapsed := t.Sub(start)
		log.WithFields(log.Fields{
			"Duration": elapsed,
		}).Debug(logMessage)
	}
}

// MapGenerationAttempt attempts to generate a map for the specified game type
// It is regarded as an attempt, as the randomization can produce maps that are not valid and thus discarded
func MapGenerationAttempt(gameType game.GameType, verbose bool) game.Board {
	start := time.Now()
	log.Debug(" > Created a new board start")
	tiles := generateTiles(gameType)
	distributeNumbers(gameType, tiles)
	if verbose {
		for _, tile := range tiles {
			log.WithFields(log.Fields{
				"Landscape": tile.Landscape,
				"Number":    tile.Number,
				"Harbor":    tile.Harbor,
			}).Debug("Tile:")
		}
	}
	boardMap := distributeTiles(gameType, tiles, verbose)
	harborMap := distributeHarbors(gameType)
	updateTilesWithHarbors(boardMap, harborMap)

	board := &game.Board{
		Tiles:    tiles,
		Board:    boardMap,
		GameType: gameType,
		Harbors:  harborMap,
	}

	debugLogDuration(start, " < Created a new board finish")
	return *board
}

func updateTilesWithHarbors(tiles map[string][]*model.Tile, harbors map[string]*model.Harbor) {
	start := time.Now()
	for location, harbor := range harbors {
		column := location[0:1]
		indexString := location[1:2]
		index, _ := strconv.Atoi(indexString)
		tile := tiles[column][index]
		tile.Harbor = *harbor
	}

	debugLogDuration(start, " - updateTilesWithHarbors")
}

func generateTiles(gameType game.GameType) []*model.Tile {
	start := time.Now()

	tiles := make([]*model.Tile, 0, gameType.TilesCount)
	tiles = append(tiles, addTilesOfType(gameType.DesertCount, *model.Desert)...)
	tiles = append(tiles, addTilesOfType(gameType.FieldCount, *model.Field)...)
	tiles = append(tiles, addTilesOfType(gameType.ForestCount, *model.Forest)...)
	tiles = append(tiles, addTilesOfType(gameType.MountainCount, *model.Mountain)...)
	tiles = append(tiles, addTilesOfType(gameType.PastureCount, *model.Pasture)...)
	tiles = append(tiles, addTilesOfType(gameType.RiverCount, *model.Hill)...)

	debugLogDuration(start, " - generateTiles")
	return tiles
}

func addTilesOfType(numberOfTiles int, landscape model.Landscape) []*model.Tile {
	start := time.Now()
	tiles := make([]*model.Tile, numberOfTiles, numberOfTiles)
	for i := 0; i < numberOfTiles; i++ {
		tile := model.Tile{
			Landscape: landscape,
			Harbor:    *model.HarborNone,
		}
		log.WithFields(log.Fields{
			"Tile": tile,
		}).Debug("Created new tile")
		tiles[i] = &tile
	}

	debugLogDuration(start, " - addTilesOfType")
	return tiles
}

func distributeNumbers(game game.GameType, tileSet []*model.Tile) {
	start := time.Now()
	numbersAllocated := make([]int, 0, game.TilesCount-game.DesertCount)
	randomRange := game.TilesCount - game.DesertCount // desert tile doesn't get a number
	log.Debug("Allocating numbers to Tiles")
	for i := 0; i < game.TilesCount; i++ {
		if tileSet[i].Landscape == *model.Desert {
			tileSet[i].Number = *model.NumberEmpty
			continue
		}
		drawnNumber := drawTileNumber(randomRange, numbersAllocated)
		number := game.NumberSet[drawnNumber]
		numbersAllocated = append(numbersAllocated, drawnNumber)
		tileSet[i].Number = *number
	}

	debugLogDuration(start, " - distributeNumbers")
}

func distributeTiles(gameType game.GameType, tileSet []*model.Tile, verbose bool) map[string][]*model.Tile {
	start := time.Now()
	var tilesOnBoard map[string][]*model.Tile
	tilesOnBoard = make(map[string][]*model.Tile)

	randomRange := gameType.TilesCount
	numbersAllocated := make([]int, 0, gameType.TilesCount)
	for gridLane, tilesInLane := range gameType.BoardLayout {
		tilesLine := make([]*model.Tile, tilesInLane, tilesInLane)
		for i := 0; i < tilesInLane; i++ {
			drawnTileNumber := drawTileNumber(randomRange, numbersAllocated)
			tile := tileSet[drawnTileNumber]
			numbersAllocated = append(numbersAllocated, drawnTileNumber)
			tilesLine[i] = tile
		}
		tilesOnBoard[gridLane] = tilesLine
		if verbose {
			log.Debug("Current Tiles Allocated: ", len(numbersAllocated), " / 19")
			log.Debug("Tiles Allocated: ", numbersAllocated)
		}
	}

	debugLogDuration(start, " - distributeTiles")
	return tilesOnBoard
}

func distributeHarbors(gameType game.GameType) map[string]*model.Harbor {
	start := time.Now()
	var harborsOnBoard map[string]*model.Harbor
	harborsOnBoard = make(map[string]*model.Harbor)

	randomRange := gameType.HarborCount
	numbersAllocated := make([]int, 0, gameType.HarborCount)
	for _, positions := range gameType.HarborLayout {
		drawnNumber := drawTileNumber(randomRange, numbersAllocated)
		harbor := gameType.HarborSet[drawnNumber]
		numbersAllocated = append(numbersAllocated, drawnNumber)
		harborsOnBoard[positions] = harbor
	}

	debugLogDuration(start, " - distributeHarbors")
	return harborsOnBoard
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
