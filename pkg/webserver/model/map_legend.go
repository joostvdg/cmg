package model

import (
	"github.com/joostvdg/cmg/pkg/model"
)

// MapLegend Legend for API uses, which allows use of codes (which can than be mapped via the Legend
type MapLegend struct {
	Harbors    []model.Harbor
	Landscapes []model.Landscape
}
