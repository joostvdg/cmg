package model

type Resource int

const (
	Grain  Resource = 0
	Lumber Resource = 1
	Wool   Resource = 2
	Brick  Resource = 3
	Ore    Resource = 4
	None   Resource = 5
)

type Harbor struct {
	Name     string
	Resource Resource
}
