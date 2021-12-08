package board

// Returns the board after the move corresponding to moveIndex
// by honoring the order of the array Moves.
// This is an optimized version of an earlier function:
// - omit calculating score during moves
// - re-order Moves, such that first two are orthogonally opposites.
func (x Board) DoMove(moveIndex uint8) Board {
	switch moveIndex {
	case 0:
		return x.Left()
	case 1:
		return x.Down()
	case 2:
		return x.Right()
	default:
		return x.Up()
	}
}

func (x Board) Left() Board {
	var (
		y       Board
		i, j, k int // i, j: indices pointing to the row and column of the grid
		// (as if it was a 4x4 matrix, although we store it as a 16-array)
		// k is sort of pointer to the top of the stack
		cache uint8
	)
	y = New()
	// first iterate over the rows of the "matrix"
	for i = 0; i < N; i += Size {
		// stack pointer should point to the first element of the row
		k = i
		cache = 0
		// second iterate over columns of the row
		for j = i; j < i+Size; j++ {
			// skip: current element of the input is 0
			if x[j] == 0 {
				continue
			}
			// merge: current element in input matches last element in ouput
			if x[j] == cache {
				y[k] = cache + 1
				cache = 0
				k++
				continue
			}

			if cache != 0 {
				y[k] = cache
				cache = 0
				k++
			}

			cache = x[j]

		}
		if cache != 0 {
			y[k] = cache
		}
	}
	return y
}

func (x Board) Up() Board {
	var (
		y       Board
		i, j, k int
		cache   uint8
	)
	y = New()
	for i = 0; i < Size; i++ {
		k = i
		cache = 0
		for j = i; j < i+N-1; j += Size {
			// skip
			if x[j] == 0 {
				continue
			}
			// merge
			if x[j] == cache {
				cache++
				y[k] = cache
				cache = 0
				k += Size
				continue
			}
			// push
			if cache != 0 {
				y[k] = cache
				cache = 0
				k += Size
			}
			cache = x[j]
		}
		if cache != 0 {
			y[k] = cache
		}
	}
	return y
}

func (x Board) Down() Board {
	var (
		y       Board
		i, j, k int
		cache   uint8
	)
	y = New()
	for i = (N - Size); i < N; i++ {
		k = i
		cache = 0
		for j = i; j >= i-(N-Size); j -= Size {
			// skip
			if x[j] == 0 {
				continue
			}
			// merge
			if x[j] == cache {
				y[k] = cache + 1
				cache = 0
				k -= Size
				continue
			}
			// push
			if cache != 0 {
				y[k] = cache
				cache = 0
				k -= Size
			}

			cache = x[j]
		}
		if cache != 0 {
			y[k] = cache
		}
	}
	return y
}

func (x Board) Right() Board {
	var (
		y       Board
		i, j, k int
		cache   uint8
	)
	y = New()
	for i = 0; i < N; i += Size {
		k = i + Size - 1
		cache = 0
		for j = i + Size - 1; j >= i; j-- {
			// skip
			if x[j] == 0 {
				continue
			}
			// merge
			if x[j] == cache {
				y[k] = cache + 1
				cache = 0
				k--
				continue
			}
			// push
			if cache != 0 {
				y[k] = cache
				cache = 0
				k--
			}
			cache = x[j]
		}
		if cache != 0 {
			y[k] = cache
		}
	}
	return y
}
