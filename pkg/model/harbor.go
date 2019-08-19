package model

type harbors map[string]Harbor

// Harbor is a Catan harbor, consisting out of a simple name, and the resource it has the trade benefit for
type Harbor struct {
	Name string
	Resource
}

var (
	HarborGrain  = &Harbor{Name: "2:1 Grain", Resource: *Grain}
	HarborBrick  = &Harbor{Name: "2:1 Brick", Resource: *Brick}
	HarborOre    = &Harbor{Name: "2:1 Ore", Resource: *Ore}
	HarborWool   = &Harbor{Name: "2:1 Wool", Resource: *Wool}
	HarborLumber = &Harbor{Name: "2:1 Lumber", Resource: *Lumber}
	HarborAll    = &Harbor{Name: "3:1", Resource: *All}
	HarborNone   = &Harbor{Name: "None", Resource: *None}

	Harbors = harbors{
		Grain.Code:  *HarborGrain,
		Lumber.Code: *HarborLumber,
		All.Code:    *HarborAll,
		None.Code:   *HarborNone,
		Wool.Code:   *HarborWool,
		Ore.Code:    *HarborOre,
		Brick.Code:  *HarborBrick,
	}
)
