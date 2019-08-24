package mapgen

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/go-errors/errors"
	"github.com/joostvdg/cmg/pkg/game"
	"github.com/joostvdg/cmg/pkg/model"
	log "github.com/sirupsen/logrus"
)

type Game int

func GenerateMap(count int, loop bool, verbose bool, rules game.GameRules) {

	maxGenerationAttempts := 1500
	numberOfLoops := count
	if !loop {
		numberOfLoops = 1
	}

	var gameType *game.GameType
	if rules.GameType == 0 {
		gameType = game.NormalGame
	} else if rules.GameType == 1 {
		gameType = game.LargeGame
		maxGenerationAttempts = 5000 // it's more difficult
	}

	failedGenerations := 0
	totalGenerations := 0
	board, err := MapGenerationAttempt(gameType, verbose)
	for i := 0; i < numberOfLoops; i++ {
		totalGenerations++
		for !board.IsValid(&rules, gameType) && err != nil {
			if totalGenerations > maxGenerationAttempts {
				sentry.CaptureMessage(fmt.Sprintf("Could not generate map of type %v, tried %v times", gameType, totalGenerations))
				sentry.Flush(time.Second * 5)
				log.Fatal("Can not generate a map... (1000+ runs)")
			}
			log.Debug(fmt.Sprintf("Loop %v::%v", i, failedGenerations))
			totalGenerations++
			failedGenerations++
			board, err = MapGenerationAttempt(gameType, verbose)
		}
		//board.PrintToConsole()
	}
	log.WithFields(log.Fields{
		"Map Generation Loops":    numberOfLoops,
		"Map Generation Failures": failedGenerations,
	}).Debug("Finished generation loop:")
}

// MapGenerationAttempt attempts to generate a map for the specified game type
// It is regarded as an attempt, as the randomization can produce maps that are not valid and thus discarded
func MapGenerationAttempt(gameType *game.GameType, verbose bool) (game.Board, error) {
	start := time.Now()
	log.Debug("Generating new Map")

	tiles := generateTiles(gameType)
	if log.IsLevelEnabled(log.DebugLevel) {
		for _, tile := range tiles {
			log.WithFields(log.Fields{
				"Landscape": tile.Landscape,
				"Number":    tile.Number,
				"Harbor":    tile.Harbor,
			}).Debug(" - Tile:")
		}
	}
	distributeNumbers(gameType, tiles)
	if log.IsLevelEnabled(log.DebugLevel) {
		for _, tile := range tiles {
			log.WithFields(log.Fields{
				"Landscape": tile.Landscape,
				"Number":    tile.Number,
				"Harbor":    tile.Harbor,
			}).Debug(" - Tile:")
		}
	}
	boardMap := distributeTiles(gameType, tiles, verbose)
	err := updateTilesWithHarbors(boardMap, gameType)

	board := &game.Board{
		Tiles:    tiles,
		Board:    boardMap,
		GameType: gameType,
	}
	t := time.Now()
	elapsed := t.Sub(start)
	log.WithFields(log.Fields{
		"Duration": elapsed,
	}).Debug("Created a new board")

	return *board, err
}

func updateTilesWithHarbors(tiles [][]*model.Tile, gameType *game.GameType) error {
	log.Debug("Allocating harbors to tiles")
	randomRange := gameType.HarborCount
	numbersAllocated := make([]int, 0, gameType.HarborCount)
	for _, tileCode := range gameType.HarborLayout {
		drawnNumber := drawTileNumber(randomRange, numbersAllocated)
		harbor := gameType.HarborSet[drawnNumber]
		tile := tiles[tileCode.Column][tileCode.Row]
		for i := 0; i < 3; i++ {
			if !(tile.Landscape.Resource == harbor.Resource) {
				numbersAllocated = append(numbersAllocated, drawnNumber)
				tile.Harbor = *harbor
				continue
			}
		}
		// shit, no suitable harbor
		return errors.New("Could not assign a harbor of different type, map is invalid")
	}
	return nil
}

func generateTiles(gameType *game.GameType) []*model.Tile {
	tiles := make([]*model.Tile, 0, gameType.TilesCount)
	tiles = append(tiles, addTilesOfType(gameType.DesertCount, *model.Desert)...)
	tiles = append(tiles, addTilesOfType(gameType.FieldCount, *model.Field)...)
	tiles = append(tiles, addTilesOfType(gameType.ForestCount, *model.Forest)...)
	tiles = append(tiles, addTilesOfType(gameType.MountainCount, *model.Mountain)...)
	tiles = append(tiles, addTilesOfType(gameType.PastureCount, *model.Pasture)...)
	tiles = append(tiles, addTilesOfType(gameType.RiverCount, *model.Hill)...)
	return tiles
}

func addTilesOfType(numberOfTiles int, landscape model.Landscape) []*model.Tile {
	tiles := make([]*model.Tile, numberOfTiles, numberOfTiles)
	for i := 0; i < numberOfTiles; i++ {
		tile := model.Tile{
			Landscape: landscape,
			Harbor:    *model.HarborNone,
		}
		log.WithFields(log.Fields{
			"Tile": tile,
		}).Debug(" - Created new tile")
		tiles[i] = &tile
	}
	return tiles
}

func distributeNumbers(game *game.GameType, tileSet []*model.Tile) {
	numbersAllocated := make([]int, 0, game.TilesCount-game.DesertCount)
	randomRange := game.TilesCount - game.DesertCount // desert tile doesn't get a number
	log.Debug(" > Allocating numbers to Tiles")
	for i := 0; i < game.TilesCount; i++ {
		log.Debugf("  - Tile (%d/%d): %v", i +1, game.TilesCount, *tileSet[i])
		if tileSet[i].Landscape == *model.Desert {
			tileSet[i].Number = *model.NumberEmpty
			continue
		}
		drawnNumber := drawTileNumber(randomRange, numbersAllocated)
		number := game.NumberSet[drawnNumber]
		numbersAllocated = append(numbersAllocated, drawnNumber)
		tileSet[i].Number = *number
	}
}

func distributeTiles(gameType *game.GameType, tileSet []*model.Tile, verbose bool) [][]*model.Tile {
	var tilesOnBoard [][]*model.Tile
	tilesOnBoard = make([][]*model.Tile, len(gameType.BoardLayout), len(gameType.BoardLayout))
	log.Debug("Distributing Tiles")

	randomRange := gameType.TilesCount
	numbersAllocated := make([]int, 0, gameType.TilesCount)
	for column, rowLength := range gameType.BoardLayout {

		log.WithFields(log.Fields{
			"Column":     column,
			"Row Length": rowLength,
		}).Debug(" - Generating Row of Tiles:")

		for i := 0; i < rowLength; i++ {
			drawnTileNumber := drawTileNumber(randomRange, numbersAllocated)
			tile := tileSet[drawnTileNumber]
			numbersAllocated = append(numbersAllocated, drawnTileNumber)
			tilesOnBoard[column] = append(tilesOnBoard[column], tile)
		}

		if verbose {
			log.Debug("Current Tiles Allocated: ", len(numbersAllocated), " / 19")
			log.Debug("Tiles Allocated: ", numbersAllocated)
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
