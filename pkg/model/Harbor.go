package model

type Resource int

const (
	All    Resource = 0
	Lumber Resource = 1
	Wool   Resource = 2
	Grain  Resource = 3
	Brick  Resource = 4
	Ore    Resource = 5
	None   Resource = 6 // Desert
)

// Harbor is a Catan harbor, consisting out of a simple name, and the resource it has the trade benefit for
type Harbor struct {
	Name     string
	Resource Resource
}
