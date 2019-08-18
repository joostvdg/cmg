package game

// GameRules the rules for generating this Game's map
type GameRules struct {
	MaximumScore         int
	MinimumScore         int
	MaximumResourceScore int
	MinimumResourceScore int
	MaxOver300           int
	GameType             int
}
