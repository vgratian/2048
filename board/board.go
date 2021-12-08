package board

import (
	"fmt"
	"math"
)

// Type board holds the N tiles of a 2048 Board.
// Rather than storing decimal values, it exploits the fact that
// these values are always equal to 2^k (except for the empty tile)
// and stores only the exponent.
type Board []uint8

var (
	// Default size of the board: 4x4
	Size int = 4
	// Number of tiles in the board (= Size^2)
	N int = 16
	// The four possible moves for the player
	Moves = [4]string{"Left", "Down", "Right", "Up"}
)

// New returns a new Board instance with all tile values initialized
// to "0"s
func New() Board {
	return make(Board, N, N)
}

// NewEncode returns a new Board with tiles values initialized from
// the given decimal values.
func NewEncode(values []int) Board {
	if len(values) != N {
		panic(fmt.Sprintf("invalid size: %d (expected %d)", len(values), N))
	}

	x := New()

	for i, v := range values {
		x[i] = Encode(v)
	}

	return x
}

// Encode returns the exponent (power of 2) of a decimal tile value.
func Encode(d int) uint8 {
	if d == 0 {
		return 0
	}
	if d%2 == 0 {
		return uint8(math.Log2(float64(d)))
        // this line can be optimized, but that's not so urgent
        // since the function is used infrequently.
	}
	panic(fmt.Sprintf("invalid number for a tile: %d", d))
}

// Decode converts a tile value back to its decimal notation.
func Decode(k uint8) int {
	if k == 0 {
		return 0
	}
	return 1 << k
}

// Decode returns human-friendly, decimal values of the board tiles.
func (x Board) Decode() []int {
	values := make([]int, N)
	for i, v := range x {
		values[i] = Decode(v)
	}
	return values
}

// Copy returns an exact copy of Board x
func (x Board) Copy() Board {
	y := New()
	copy(y, x)
	return y
}

// Equals compares if boards x, y have exactly same tile values
func (x Board) Equals(y Board) bool {
	for i := range x {
		if x[i] != y[i] {
			return false
		}
	}
	return true
}

// GameOver returns true if there are no more legal moves
// on ther board, otherwise false. This depends on who's
// move it is currently:
// Player's turn - game is over when none of the 4 Moves
//                 are possible.
// Opponent's turn - game is over when no empty tiles are
//                 left.
func (x Board) GameOver(playersTurn bool) bool {
	// if there is at least one empty tile, we know
	// that game is not over for both.
	for _, v := range x {
		if v == 0 {
			return false
		}
	}
	// even if there are no empty tiles, it is possible
	// that one of the Moves merges some of the tiles,
	// if so, game is not over yet.
	// @TODO: check moves in non-greedy way
	if playersTurn {
		var i uint8
		for i = 0; i < 4; i++ {
			if y := x.DoMove(i); !x.Equals(y) {
				return false
			}
		}
	}
	return true
}

// Return number of empty tiles
func (x Board) CountOfZeros() uint32 {
	var count uint32
	for _, v := range x {
		if v == 0 {
			count++
		}
	}
	return count
}

// Return array with indices of empty tiles
func (x Board) IndicesOfZeros() []uint8 {
	indices := make([]uint8, 0, N)
	for i, v := range x {
		if v == 0 {
			indices = append(indices, uint8(i))
		}
	}
	return indices
}
