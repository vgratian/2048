package engine

import (
    "github.com/vgratian/2048/board"
)

func searchAlphaBetaGreedy(b board.Board, depth uint32, playersTurn bool, alpha, beta float32) float32 {

    if depth == 0 {
        return eval(b)
    }

    children := expandGreedy(b, playersTurn)

    if len(children) == 0 {  // means Game Over
        return 0
    }

    NumNodes += uint64(len(children))

    if playersTurn {
        maxEval := MIN_EVAL_F32
        for _, child := range children {
            value := searchAlphaBetaGreedy(child, depth-1, false, alpha, beta)
            if value > maxEval {
                maxEval = value
            }
            if maxEval > alpha {
                alpha = maxEval
            }
            if alpha >= beta {
                break
            }
        }
        return maxEval
    }

    minEval := MAX_EVAL_F32
    for _, child := range children {
        value := searchAlphaBetaGreedy(child, depth-1, true, alpha, beta)
        if value < minEval {
            minEval = value
        }
        if minEval < beta {
            beta = minEval
        }
        if beta <= alpha {
            break
        }
    }
    return minEval
}

func expandGreedy(b board.Board, playersTurn bool) []board.Board {

    var i uint8

    if playersTurn {
        nodes := make([]board.Board, 0, 4)
        for i=0; i<4; i++ {
            if ch := b.DoMove(i); ! ch.Equals(b) {
                nodes = append(nodes, ch)
            }
        }
        return nodes
    }

    nodes := make([]board.Board, 0, 15)
    for i=0; i<16; i++ {
        if b[i] == 0 {
            nodes = append(nodes, b.Insert(i, 1))
        }
    }
    return nodes
}
