package model

import "github.com/joostvdg/cmg/pkg/model"

// Map the Catan Map, a wrapper around the Game Board
type Map struct {
	GameType string
	Board    map[string][]*model.Tile
	GameCode string
	Error    string
}
