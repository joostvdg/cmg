package model

// Number the number fiche that comes on top the Tile, has Name for printing and a Number for calculating
// Also has a Score, which is the probability score of the number being rolled with two dices
type Number struct {
	Name   string
	Number int
	Score  int
}
