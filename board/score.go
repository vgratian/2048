package board

// Instead of keeping track of score after each move
// it is easier to calculate it only once when needed.
//
// Note: this method assumes that only "0" and "2" are
// atomic tile values (i.e. not resulting from a merge).
// This is not true for all 2048: in some of them "4"
// can be an atomic value as well. See explanation below.
func (b Board) Score() uint32 {
	var totalScore uint32
	for _, k := range b {
		totalScore += tileScore(k)
	}
	return totalScore
}

// Calculate the score contribution of a single tile.
// This is the sum of all merges that resulted in the
// current tile value k.
//
// Examples:
//
// Tile value is "0" or "2": score is 0 (these are
// atomic values, not resulting from a merge).
//
// Tile value is "8": core is 16, since "8" resulted
// from the merge of two "4"s (4+4=8), which in turn
// resulted from the merge of four "2"s ((2+2)+(2+2)=(4+4)).
// So the score of this tile is 8 + (4+4) = 16.
func tileScore(k uint8) uint32 {
	if k < 2 {
		return 0
	}
	return 2*tileScore(k-1) + (1 << k)
	// the shift converts exponent to decimal value
}

// ScoreLazy does the same as Score(), but uses pre-
// calculated tile scores, and is about 7x faster.
func (b Board) ScoreLazy() uint32 {
	var totalScore uint32
	for _, tile := range b {
		totalScore += tileScores[tile]
	}
	return totalScore
}

// Pre-calculated values for each tile.
// A 4x4 board needs only first 16, if board size
// is larger than that, ScoreLazy() might panic if
// an extremely high value is requested.
var tileScores = []uint32{
	0, 0, 4, 16, 48, 128, 320, 768, 1792, 4096,
	9216, 20480, 45056, 98304, 212992, 458752,
	983040, 2097152, 4456448, 9437184,
}
