package model

// Tile the tiles that make up the Catan game board
// A tile consists out of a Landscape, a Number, and a Harbor
// The Resource is for easier processing/printing
type Tile struct {
	Landscape Landscape
	Number    Number
	Harbor    Harbor
}

// TileCode is an identifier for the Tile withing a Board
// Useful for creating mappings, such as for validating Adjacent Tiles
type TileCode struct {
	Column int
	Row    int
}
