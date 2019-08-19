package game

// GameRules the rules for generating this Game's map
type GameRules struct {
	MaximumScore         int
	MinimumScore         int
	MaximumResourceScore int
	MinimumResourceScore int
	MaxOver300           int
	GameType             int
	GameTypeString       string
	Generations          int
	Delimiter            string
}

var (
	DefaultGameRulesNormal = GameRules{
		MaximumScore:         361,
		MinimumScore:         165,
		MaximumResourceScore: 130,
		MinimumResourceScore: 30,
		MaxOver300:           10,
		GameType:             0,
		GameTypeString:       "Normal",
		Generations:          1500,
		Delimiter:            "_",
	}

	DefaultGameRulesLarge = GameRules{
		MaximumScore:         365,
		MinimumScore:         156,
		MaximumResourceScore: 130,
		MinimumResourceScore: 65,
		MaxOver300:           22,
		GameType:             1,
		GameTypeString:       "Large",
		Generations:          5000,
		Delimiter:            "_",
	}
)
