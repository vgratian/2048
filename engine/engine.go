package engine

import (
    "github.com/vgratian/2048/board"
    "github.com/vgratian/2048/params"
    "github.com/vgratian/2048/tui"
    "fmt"
    "math/rand"
    "time"
)

var (
    boardForBenchmark = board.Board{
        2,  2,  3,  3,
        3,  6, 10,  7,
        5,  5,  8,  2,
        0,  1,  5,  0,
    }
    NumNodes uint64
)

func SelfPlay() {

    tui.ClearScreen()

    var (
        score, step, depth uint32
        move uint8
        eval [4]float32
        runTime, totalTime time.Duration
        start time.Time
        //history []board.Board
    )
    brd := board.New()
    //history = make([]board.Board, 0, 500)
    ourTurn := false

    tui.PrintScoreHeader()
    tui.PrintScore2(0, 0, 0)
    tui.PrintBoard(brd, "initial board:")

    for {

        if int(step) > params.MaxMoves {
            break
        }

        tui.ScrollBack()

        if ourTurn {
            depth = AdjustDepth(brd)
            start = time.Now()
            //move, eval = SPminimax(brd, depth)
            //move, eval = SPminimax2(brd, depth)
            //move, eval = SPminimax3(brd, depth)
            move, eval = SearchAlphaBeta(brd, depth)
            runTime = time.Since(start)
            totalTime += runTime
            brd = brd.DoMove(move)
            //history = append(history, brd)
            score = brd.Score()
            tui.PrintScore0(step, score, depth, NumNodes, move, eval, runTime.Seconds())
            tui.PrintBoard(brd, "")
            step++
        } else {
            move, brd = playOpponent(brd)
            tui.PrintScore2(step, score, move)
            tui.PrintBoard(brd, "")
        }

        ourTurn = !ourTurn

        if brd.GameOver(ourTurn) {
            break
        }

        //time.Sleep(time.Duration(params.TurnDelay) * time.Millisecond)
    }

    tui.ScrollDown()

    //for _, b := range history {
    //    fmt.Println(b)
    //}

    fmt.Println(" ** GAME OVER ** ")
    fmt.Println(" ourTurn: ", ourTurn)
    fmt.Println(" maxMoves:", params.MaxMoves)
    fmt.Println(" num Moves:", step)
    fmt.Println(" default Depth:", params.MaxDepth)
    avgTime := totalTime.Seconds() / float64(step)
    fmt.Println(" average think time (s):", avgTime)
}

func SingleSearch() {
    singleSearch(boardForBenchmark.Copy(), uint32(params.MaxDepth))
}

func singleSearch(brd board.Board, depth uint32) {

    //playersTurn := true

    //fmt.Println("before:")
    //tui.PrintBoard(brd, "")
    //fmt.Println()
    //fmt.Println()

    start := time.Now()
    move, eval := SearchAlphaBeta(brd, depth)
    thinkTime := time.Since(start)

    tui.PrintScoreHeader()
    tui.PrintScore0(1, 0, depth, NumNodes, move, eval, thinkTime.Seconds())

    //brd, _ = brd.DoMove(move)
    //fmt.Println("after:")
    //tui.PrintBoard(brd, "")
    //fmt.Println()
    //fmt.Println()

}

func playOpponent(b board.Board) (uint8, board.Board) {
    indices := b.IndicesOfZeros()
    index := indices[rand.Intn(len(indices))]
    b.Insert(index, uint8(1))
    return index, b
}

func AdjustDepth(b board.Board) uint32 {
    zeros := b.CountOfZeros()
    depth := float32(params.MaxDepth)
    bonus := float32(params.N-int(zeros))
    alpha := bonus/float32((params.N*2))
    depth += bonus*alpha
    return uint32(depth)
}
