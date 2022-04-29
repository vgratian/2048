package engine

import (
    "testing"
    "github.com/vgratian/2048/board"
)

var (
    boardPlayerLeft = board.Board{
        3,  4,  0,  0,
        3,  6, 10,  7,
        6,  8,  2,  0,
        1,  5,  0,  0,
    }
    boardPlayerRight = board.Board{
        0,  0,  3,  4,
        3,  6, 10,  7,
        0,  6,  8,  2,
        0,  0,  1,  5,
    }
    boardPlayerDown = board.Board{
        0,  2,  3,  0,
        2,  6, 10,  3,
        3,  5,  8,  7,
        5,  1,  5,  2,
    }
    boardOpponent1 = board.Board{
        2,  2,  3,  3,
        3,  6, 10,  7,
        5,  5,  8,  2,
        1,  1,  5,  0,
    }
    boardOpponent2 = board.Board{
        2,  2,  3,  3,
        3,  6, 10,  7,
        5,  5,  8,  2,
        0,  1,  5,  1,
    } 
)

func TestNextChildForOpponent(t *testing.T) {
    var (
        child board.Board
        nextIndex uint8
    )

    // first non-empty cell is at index 12, so first time
    // we call the function, this cell should get "1" and
    // next index should be 13
    t.Log("1st invocation: expecting child (inserted value at index 12) and nextIndex [13]")
    child, nextIndex = nextChildForOpponent(boardForBenchmark, nextIndex)
    if nextIndex != 13 {
        t.Errorf("nextIndex = %d", nextIndex)
    }
    if child == nil || ! child.Equals(boardOpponent1) {
        t.Errorf("child = %v\nexpected %v", child, boardOpponent1)
    }

    t.Log("2nd invocation: expecting child (inserted value at index 15) and nextIndex [16]")
    child, nextIndex = nextChildForOpponent(boardForBenchmark, nextIndex)
    if nextIndex != 16 {
        t.Errorf("nextIndex = %d", nextIndex)
    }
    if child == nil || ! child.Equals(boardOpponent2) {
        t.Errorf("child = %v\nexpected %v", child, boardOpponent2)
    }
 
    t.Log("3rd invocation: expecting child (nil) and nextIndex [0]")
    child, nextIndex = nextChildForOpponent(boardForBenchmark, nextIndex)
    if nextIndex != 0 {
        t.Errorf("nextIndex = %d", nextIndex)
    }
    if child != nil {
        t.Errorf("child = %v\nexpected [nil]", child)
    }
}

func TestNextChildForPlayer(t *testing.T) {

    var (
        child board.Board
        nextIndex uint8
    )

    // note: function should return moves in following order:
    // 0 = up,  1 = down, 2 = left, 3 = right

    // should skip "up" (impossible move)
    // so first child should be "down"
    t.Log("1st invocation: expecting child [Down] and nextIndex [2]")
    child, nextIndex = nextChildForPlayer(boardForBenchmark, nextIndex)
    if nextIndex != 2 {
        t.Errorf("nextIndex = %d", nextIndex)
    }
    if child == nil || ! child.Equals(boardPlayerDown) {
        t.Errorf("child = %v\nexpected %v", child, boardPlayerDown)
    }

    t.Log("2nd invocation: expecting child [Left] and nextIndex [3]")
    child, nextIndex = nextChildForPlayer(boardForBenchmark, nextIndex)
    if nextIndex != 3 {
        t.Errorf("nextIndex = %d", nextIndex)
    }
    if child == nil || ! child.Equals(boardPlayerLeft) {
        t.Errorf("child = %v\nexpected %v", child, boardPlayerLeft)
    }

    t.Log("3rd invocation: expecting child [Right] and nextIndex [4]")
    child, nextIndex = nextChildForPlayer(boardForBenchmark, nextIndex)
    if nextIndex != 4 {
        t.Errorf("nextIndex = %d", nextIndex)
    }
    if child == nil || ! child.Equals(boardPlayerRight) {
        t.Errorf("child = %v\nexpected %v", child, boardPlayerRight)
    }

    // we should get no more children
    t.Log("4th invocation: expecting child [nil] and nextIndex [0]")
    child, nextIndex = nextChildForPlayer(boardForBenchmark, nextIndex)
    if nextIndex != 0 {
        t.Errorf("nextIndex = %d", nextIndex)
    }
    if child != nil {
        t.Errorf("child = %v\nexpected [nil]", child)
    }
}

func TestExpandGreedy(t *testing.T) {

    t.Log("testing for Player: expecting 3 results")
    children := expandGreedy(boardForBenchmark, true)
    if len(children) != 3 {
        t.Fatalf("got %d children", len(children))
    }
    if ! children[0].Equals(boardPlayerDown) {
        t.Errorf("1st child = %v, expected %v", children[0], boardPlayerDown)
    }
    if ! children[1].Equals(boardPlayerLeft) {
        t.Errorf("2nd child = %v, expected %v", children[1], boardPlayerLeft)
    }
    if ! children[2].Equals(boardPlayerRight) {
        t.Errorf("3rd child = %v, expected %v", children[2], boardPlayerRight)
    }

    t.Log("testing for Opponent: expecting 2 results")
    children = expandGreedy(boardForBenchmark, false)
    if len(children) != 2 {
        t.Fatalf("got %d children", len(children))
    }
    if ! children[0].Equals(boardOpponent1) {
        t.Errorf("1st child = %v, expected %v", children[0], boardOpponent1)
    }
    if ! children[1].Equals(boardOpponent2) {
        t.Errorf("2nd child = %v, expected %v", children[1], boardOpponent2)
    }
}

func BenchmarkSearchAlphaBeta(b *testing.B) {

    var depth uint32 = 10

    singleSearch(boardForBenchmark.Copy(), depth)
    b.Logf("N=%d", b.N)

    b.ReportAllocs()
    b.ResetTimer()

    for i := 0; i < b.N; i++ {
        SearchAlphaBeta(boardForBenchmark.Copy(), depth)
    }
}
