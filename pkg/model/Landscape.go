package model

type LandscapeCode int

const (
	Desert   LandscapeCode = 0
	Forest   LandscapeCode = 1
	Pasture  LandscapeCode = 2
	Field    LandscapeCode = 3
	River    LandscapeCode = 4
	Mountain LandscapeCode = 5
)

type Landscape struct {
	Name     string
	Code     LandscapeCode
	Resource string
}
