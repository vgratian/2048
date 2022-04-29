package board

import (
	"testing"
)

// For testing board state and score change after 4 moves
type testBoard struct {
    board Board
    // expected encoding of the board
    encoded []uint8
    // expected results of the 4 moves
    moves [4][]int
    // expected score change
    scores [4]uint32
}

// For testing if game is over at current state
type testGameOver struct {
    board Board
    forPlayer bool
    forOpponent bool
}

// boards for testing, each board has:
// - initial decimal values representing the tiles
// - binary encoding that we expect
// - board state after the 4 moves (decimal)
// - score change after each if the 4 moves
var (
    a = NewEncode([]int{
		0, 0, 2, 2,
		2, 0, 0, 2,
		8, 2, 4, 0,
		2, 8, 8, 8,
	})

	aEncoded = []uint8 {
		0, 0, 1, 1,
		1, 0, 0, 1,
		3, 1, 2, 0,
		1, 3, 3, 3,
	}

	aMovedLeft = []int{
		4, 0, 0, 0,
		4, 0, 0, 0,
		8, 2, 4, 0,
		2, 16,8, 0,
	}
    aLeftScore uint32 = 24

    aMovedRight = []int{
		0, 0, 0, 4,
		0, 0, 0, 4,
		0, 8, 2, 4,
		0, 2, 8, 16,
	}
	aRightScore uint32 = 24

	aMovedUp = []int{
		2, 2, 2, 4,
		8, 8, 4, 8,
		2, 0, 8, 0,
		0, 0, 0, 0,
	}
	aUpScore uint32 = 4

	aMovedDown = []int{
		0, 0, 0, 0,
		2, 0, 2, 0,
		8, 2, 4, 4,
		2, 8, 8, 8,
	}
    aDownScore uint32 = 4

    b = NewEncode([]int{
		0, 2, 4, 4,
		2, 4, 2, 2,
		2, 4, 8, 2,
		16,16,0, 0,
    })

    bEncoded = []uint8{
		0, 1, 2, 2,
		1, 2, 1, 1,
		1, 2, 3, 1,
		4, 4, 0, 0,
    }
    
    bMovedLeft = []int{
		2, 8, 0, 0,
		2, 4, 4, 0,
		2, 4, 8, 2,
		32,0, 0, 0,
	}
	bLeftScore uint32 = 44

	bMovedRight = []int{
		0, 0, 2, 8,
		0, 2, 4, 4,
		2, 4, 8, 2,
		0, 0, 0, 32,
	}
	bRightScore uint32 = 44

	bMovedUp = []int{
		4, 2, 4, 4,
		16,8, 2, 4,
		0, 16,8, 0,
		0, 0, 0, 0,
	}
	bUpScore uint32 = 16

	bMovedDown = []int{
		0, 0, 0, 0,
		0, 2, 4, 0,
		4, 8, 2, 4,
		16,16,8, 4,
	}
	bDownScore uint32 = 16

    c = NewEncode([]int{
        2, 2, 2, 2,
        4, 4, 8, 8,
        2, 2, 2, 8,
        0, 0, 8, 8,
    })

	cEncoded = []uint8 {
        1, 1, 1, 1,
        2, 2, 3, 3,
        1, 1, 1, 3,
        0, 0, 3, 3,
    }

    cMovedLeft = []int{
        4, 4, 0, 0,
        8, 16,0, 0,
        4, 2, 8, 0,
        16,0, 0, 0,
    }
    cLeftScore uint32 = 52

    cMovedRight = []int{
        0, 0, 4, 4,
        0, 0, 8,16,
        0, 2, 4, 8,
        0, 0, 0, 16,
    }
    cRightScore uint32 = 52

    cMovedUp = []int{
        2, 2, 2, 2,
        4, 4, 8,16,
        2, 2, 2, 8,
        0, 0, 8, 0,
    }
    cUpScore uint32 = 16

    cMovedDown = []int{
        0, 0, 2, 0,
        2, 2, 8, 2,
        4, 4, 2, 8,
        2, 2, 8,16,
    }
    cDownScore uint32 = 16
    
    // these ones are only for testing GameOver()
    d = NewEncode([]int{
        16,16,16,16,
         8, 8, 8, 8,
         4, 4, 4, 4,
         2, 2, 2, 2,
     })

    dGameOverPlayer = false
    dGameOverOpponent = true

    e = NewEncode([]int{
        4, 2, 4, 2,
        2, 4, 2, 4,
        4, 2, 4, 2,
        2, 4, 2, 4,
    })
    eGameOverPlayer = true
    eGameOverOpponent = true

    aTestBoard = testBoard{a, aEncoded,
        [4][]int{aMovedLeft, aMovedDown, aMovedRight, aMovedUp},
        [4]uint32{aLeftScore, aDownScore, aRightScore, aUpScore},
    }
    bTestBoard = testBoard{b, bEncoded,
        [4][]int{bMovedLeft, bMovedDown, bMovedRight, bMovedUp},
        [4]uint32{bLeftScore, bDownScore, bRightScore, bUpScore},
    }
    cTestBoard = testBoard{c, cEncoded,
        [4][]int{cMovedLeft, cMovedDown, cMovedRight, cMovedUp},
        [4]uint32{cLeftScore, cDownScore, cRightScore, cUpScore},
    }

    testBoards = []testBoard{aTestBoard, bTestBoard, cTestBoard}


    dGameOverCase = testGameOver{d, dGameOverPlayer, dGameOverOpponent}
    eGameOverCase = testGameOver{e, eGameOverPlayer, eGameOverOpponent}

    testCases = []testGameOver{dGameOverCase, eGameOverCase}
)

func TestEncoding(t *testing.T) {
    t.Log("Testing encoding from decimal (x) to binarry representation (n^x)")
    for i, x := range testBoards {
        if ! x.board.Equals(x.encoded) {
            t.Errorf("Encoding board #%d: FAIL", i)
            t.Logf("  -> Expected: %v", x.encoded)
            t.Logf("  -> Actual:   %v", x.board)
        }
    }
}

func TestMoves(t *testing.T) {
    t.Logf("Testing state and score change after %d moves on %d boards", len(Moves), len(testBoards))
    for i, move := range Moves {
        t.Logf("=> Testing Move: %s", move)

        for j, x := range testBoards {

            t.Logf("test Board #%d: %v", j, x.board.Decode())

            movedBoard := x.board.DoMove(uint8(i))
            score := movedBoard.Score() - x.board.Score()

            t.Logf(" -> Result: [+%d]: %v", score, movedBoard.Decode())

            if score != x.scores[i] {
                t.Errorf("    FAIL -> expected score: %d", x.scores[i])
            }

            if ! movedBoard.Equals(NewEncode(x.moves[i])) {
                t.Errorf("    FAIL -> expected board: %v", x.moves[i])
            }

        }
    }
}

func TestGameOver(t *testing.T) {
    t.Logf("Testing if checking for GameOver is correct on %d boards", len(testCases))
    for i, x := range testCases {

        t.Logf("=> test Board #%d: %v", i, x.board.Decode())

        // check for opponent
        forO := x.board.GameOver(false)
        if forO == x.forOpponent {
            t.Logf("OK -- GameOver for Opponent = %t", forO)
        } else {
            t.Errorf("FAIL -- GameOver for Opponent = %t", forO)
        }

        // check for player
        forP := x.board.GameOver(true)
        if forP == x.forPlayer {
            t.Logf("OK -- GameOver for Player = %t", forP)
        } else {
            t.Errorf("Fail -- GameOver for Player = %t", forP)
        }
    }
}
