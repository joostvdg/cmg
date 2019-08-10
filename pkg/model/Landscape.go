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

// Landscape is the Catan landscape type, such as Forest, Mountains
// Each has a name (for readability), and a code for compact data transfer/processing
type Landscape struct {
	Name string
	Code LandscapeCode
}
