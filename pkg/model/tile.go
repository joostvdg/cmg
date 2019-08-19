package model

// Tile the tiles that make up the Catan game board
// A tile consists out of a Landscape, a Number, and a Harbor
// The Resource is for easier processing/printing
type Tile struct {
	Landscape Landscape
	Number    Number
	Harbor    Harbor
}
