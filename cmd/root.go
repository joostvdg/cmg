package cmd

import (
	"fmt"
	"github.com/joostvdg/cmg/pkg/mapgen"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
)

var GenCount int
var GenLoop bool
var Verbose bool

func init() {
	mapGenCmd.Flags().IntVar(&GenCount,"count", 0, "Number of times to generate a map, only for loop")
	mapGenCmd.Flags().BoolVar(&GenLoop, "loop", false, "Generate maps in a loop 'count' times, or just once")
	mapGenCmd.Flags().BoolVar(&Verbose, "verbose", false, "Verbose logging")
	rootCmd.AddCommand(mapGenCmd)
}

var mapGenCmd = &cobra.Command{
	Use:   "mapgen",
	Short: "Will generate a map",
	Long:  `Anything to do with generating a Catan map`,
	Run: func(cmd *cobra.Command, args []string) {
		mapgen.GenerateMap(GenCount, GenLoop, Verbose)
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