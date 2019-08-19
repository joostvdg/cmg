package model

type landscapes map[string]Landscape

// Landscape is the Catan landscape type, such as Forest, Mountain
// Each has a name (for readability), and a code for compact data transfer/processing
type Landscape struct {
	Name string
	Resource
}

var (
	Desert   = &Landscape{Name: "Desert", Resource: *None}
	Field    = &Landscape{Name: "Field", Resource: *Grain}
	Forest   = &Landscape{Name: "Forest", Resource: *Lumber}
	Pasture  = &Landscape{Name: "Pasture", Resource: *Wool}
	Mountain = &Landscape{Name: "Mountain", Resource: *Ore}
	Hill     = &Landscape{Name: "Hill", Resource: *Brick}

	Landscapes = landscapes{
		Grain.Code:  *Field,
		Lumber.Code: *Forest,
		None.Code:   *Desert,
		Wool.Code:   *Pasture,
		Ore.Code:    *Mountain,
		Brick.Code:  *Hill,
	}
)
