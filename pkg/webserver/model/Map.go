package model

import "github.com/joostvdg/cmg/pkg/model"

type Map struct {
	GameType string
	Board    map[string][]*model.Tile
}
