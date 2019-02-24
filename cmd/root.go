package cmd

import (
	"fmt"
	"github.com/joostvdg/cmg/cmd/webserver"
	"github.com/joostvdg/cmg/pkg/game"
	"github.com/joostvdg/cmg/pkg/mapgen"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
)

var GenCount int
var GenLoop bool
var Verbose bool
var MaxScore int
var MinScore int
var MaxResourceScore int
var MinResourceScore int
var MaxOver300 int
var GameType int

func init() {
	mapGenCmd.Flags().IntVar(&MaxScore, "max", 361, "Maximum Probability score of 3 adjacent tiles")
	mapGenCmd.Flags().IntVar(&MinScore, "min", 165, "Minimum Probability score of 3 adjacent tiles")
	mapGenCmd.Flags().IntVar(&MaxResourceScore, "maxResource", 130, "Maximum average Probability score for resources per tile")
	mapGenCmd.Flags().IntVar(&MinResourceScore, "minResource", 30, "Minimum average Probability score for resources per tile")
	mapGenCmd.Flags().IntVar(&GameType, "gameType", 0, "GameType, 0 = normal, 1 = large (5or6 players)")
	mapGenCmd.Flags().IntVar(&MaxOver300, "max300", 18, "Number times the probability score of 3 adjacent tiles can exceed 300")
	mapGenCmd.Flags().IntVar(&GenCount, "count", 0, "Number of times to generate a map, only for loop")
	mapGenCmd.Flags().BoolVar(&GenLoop, "loop", false, "Generate maps in a loop 'count' times, or just once")
	mapGenCmd.Flags().BoolVar(&Verbose, "verbose", false, "Verbose logging")

	rootCmd.AddCommand(mapGenCmd)
	rootCmd.AddCommand(webServerCmd)
}

var mapGenCmd = &cobra.Command{
	Use:   "mapgen",
	Short: "Will generate a map",
	Long:  `Anything to do with generating a Catan map`,
	Run: func(cmd *cobra.Command, args []string) {
		rules := game.GameRules{
			GameType:             GameType,
			MinimumScore:         MinScore,
			MaximumScore:         MaxScore,
			MaxOver300:           MaxOver300,
			MaximumResourceScore: MaxResourceScore,
			MinimumResourceScore: MinResourceScore,
		}
		mapgen.GenerateMap(GenCount, GenLoop, Verbose, rules)
	},
}

var webServerCmd = &cobra.Command{
	Use:   "serve",
	Short: "Starts an http server",
	Long:  `Starts an http server that allows you to retrieve a generated map as json `,
	Run: func(cmd *cobra.Command, args []string) {
		webserver.StartWebserver()
	},
}

var rootCmd = &cobra.Command{
	Use:   "cmg",
	Short: "CMG is a Catan Map Generator",
	Long: `Catan Map Generator (CMG) is a map generator for the board game Catan.
				Aiming to ease creating solid maps that satisfy specific constraints.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
		log.Info("Hello from CMG!")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
