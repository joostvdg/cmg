package model

// numbers map meant to container all the numbers
type numbers map[string]Number

// Number the number fiche that comes on top the Tile
// Also has a Score, which is the probability score of the number being rolled with two dices
type Number struct {
	Number int
	Score  int
	Code   string
}

var (
	Number1  = &Number{Number: 0, Score: 1, Code: "z"} // Desert or Sea tiles
	Number2  = &Number{Number: 2, Score: 27, Code: "a"}
	Number3  = &Number{Number: 3, Score: 55, Code: "b"}
	Number4  = &Number{Number: 4, Score: 83, Code: "c"}
	Number5  = &Number{Number: 5, Score: 111, Code: "d"}
	Number6  = &Number{Number: 6, Score: 139, Code: "e"}
	Number8  = &Number{Number: 8, Score: 139, Code: "f"}
	Number9  = &Number{Number: 9, Score: 111, Code: "g"}
	Number10 = &Number{Number: 10, Score: 83, Code: "h"}
	Number11 = &Number{Number: 11, Score: 55, Code: "i"}
	Number12 = &Number{Number: 12, Score: 27, Code: "j"}

	// Numbers map container all the numbers for searching by code
	Numbers = numbers{
		Number1.Code:  *Number1,
		Number2.Code:  *Number2,
		Number3.Code:  *Number3,
		Number4.Code:  *Number4,
		Number5.Code:  *Number5,
		Number6.Code:  *Number6,
		Number8.Code:  *Number8,
		Number9.Code:  *Number9,
		Number10.Code: *Number10,
		Number11.Code: *Number11,
		Number12.Code: *Number12,
	}
)
