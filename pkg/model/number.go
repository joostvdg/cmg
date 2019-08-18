package model

// Number the number fiche that comes on top the Tile
// Also has a Score, which is the probability score of the number being rolled with two dices
type Number struct {
	Number int
	Score  int
	Code   string
}
