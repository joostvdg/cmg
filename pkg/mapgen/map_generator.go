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
		for !board.IsValid(rules, gameType, verbose) {
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

// MapGenerationAttempt attempts to generate a map for the specified game type
// It is regarded as an attempt, as the randomization can produce maps that are not valid and thus discarded
func MapGenerationAttempt(gameType game.GameType, verbose bool) game.Board {

	log.Debug("Generating new Map")
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
	harborMap := distributeHarbors(gameType)
	updateTilesWithHarbors(boardMap, harborMap)

	board := &game.Board{
		Tiles:    tiles,
		Board:    boardMap,
		GameType: gameType,
		Harbors:  harborMap,
	}
	log.Debug("Created a new board")
	return *board
}

func updateTilesWithHarbors(tiles map[string][]*model.Tile, harbors map[string]*model.Harbor) {
	for location, harbor := range harbors {
		column := location[0:1]
		indexString := location[1:2]
		index, _ := strconv.Atoi(indexString)
		tile := tiles[column][index]
		tile.Harbor = *harbor
	}
}

func generateTiles(gameType game.GameType) []*model.Tile {

	tiles := make([]*model.Tile, 0, gameType.TilesCount)
	tiles = append(tiles, addTilesOfType(gameType.DesertCount, model.Desert, model.None)...)
	tiles = append(tiles, addTilesOfType(gameType.FieldCount, model.Field, model.Grain)...)
	tiles = append(tiles, addTilesOfType(gameType.ForestCount, model.Forest, model.Lumber)...)
	tiles = append(tiles, addTilesOfType(gameType.MountainCount, model.Mountain, model.Ore)...)
	tiles = append(tiles, addTilesOfType(gameType.PastureCount, model.Pasture, model.Wool)...)
	tiles = append(tiles, addTilesOfType(gameType.RiverCount, model.River, model.Brick)...)
	return tiles
}

func addTilesOfType(number int, landscape model.LandscapeCode, resource model.Resource) []*model.Tile {
	tiles := make([]*model.Tile, number, number)
	for i := 0; i < number; i++ {
		tile := model.Tile{
			Landscape: landscape,
			Resource:  resource,
			Harbor:    model.Harbor{Resource: model.None, Name: ""},
		}
		tiles[i] = &tile
	}
	return tiles
}

func distributeNumbers(game game.GameType, tileSet []*model.Tile) {
	numbersAllocated := make([]int, 0, game.TilesCount-game.DesertCount)
	randomRange := (game.TilesCount - game.DesertCount) // desert tile doesn't get a number
	log.Debug("Allocating numbers to Tiles")
	for i := 0; i < game.TilesCount; i++ {
		if tileSet[i].Landscape == model.Desert {
			tileSet[i].Number = model.Number{Number: 0, Score: 0, Code: "z"}
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
	return tilesOnBoard
}

func distributeHarbors(gameType game.GameType) map[string]*model.Harbor {
	var harborsOnBoard map[string]*model.Harbor
	harborsOnBoard = make(map[string]*model.Harbor)

	randomRange := gameType.HarborCount
	numbersAllocated := make([]int, 0, gameType.HarborCount)
	for _, positions := range gameType.HarborLayout {
		drawnNumber := drawTileNumber(randomRange, numbersAllocated)
		harbor := gameType.HarborSet[drawnNumber]
		numbersAllocated = append(numbersAllocated, drawnNumber)
		harborPosition := positions[0]
		if len(positions) > 1 {
			harborPosition = fmt.Sprintf("%s,%s", positions[0], positions[1])
		}
		harborsOnBoard[harborPosition] = harbor
	}

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