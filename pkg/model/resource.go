package model

type Resource struct {
	Name string
	Code string
}

var (
	All    = &Resource{"All", "0"}
	Lumber = &Resource{"Lumber", "1"}
	Wool   = &Resource{"Wool", "2"}
	Grain  = &Resource{"Grain", "3"}
	Brick  = &Resource{"Brick", "4"}
	Ore    = &Resource{"Ore", "5"}
	None   = &Resource{"None", "6"} // Desert
)
