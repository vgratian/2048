package engine

import (
    "github.com/vgratian/2048/board"
    "github.com/vgratian/2048/params"
)

func SearchAlphaBeta(root board.Board, depth uint32) (uint8, [4]float32) {
    var (
        maxIndex uint8
        maxValue float32
        node board.Board
        values [4]float32
    )

    NumNodes = 0

    for i := range params.Moves {
        if node = root.DoMove(uint8(i)); ! node.Equals(root) {
            //values[i] = searchAlphaBetaGreedy(node, depth, false, MIN_EVAL_F32, MAX_EVAL_F32)
            values[i] = searchAlphaBeta(node, depth, false, MIN_EVAL_F32, MAX_EVAL_F32)
            if values[i] > maxValue {
                maxValue = values[i]
                maxIndex = uint8(i)
            }
        }
    }
    return maxIndex, values
}

// hopefully a performance improvement over _sp_alpha_beta
// by turning the _sp_expand function into a "generator"
func searchAlphaBeta(node board.Board, depth uint32, playersTurn bool, alpha, beta float32) float32 {

    if depth == 0 {
        return eval(node)
    }

    var (
        child board.Board
        childIndex uint8
        childCount uint64
        bestEval, eval float32
    )

    // player's turn, search for max evaluation score
    if playersTurn {
        bestEval = MIN_EVAL_F32
        child, childIndex = nextChildForPlayer(node, childIndex)
        if child == nil { // Game over
            return 0
        }

        for child != nil {
            childCount++
            eval = searchAlphaBeta(child, depth-1, false, alpha, beta)
            if eval > bestEval {
                bestEval = eval
            }
            if bestEval > alpha {
                alpha = bestEval
            }
            if alpha >= beta {
                break
            }
            child, childIndex = nextChildForPlayer(node, childIndex)
        }
        NumNodes += childCount
        return bestEval
    }

    // opponent's turn, search for min evaluation score
    bestEval = MAX_EVAL_F32
    child, childIndex = nextChildForOpponent(node, childIndex)
    if child == nil { // Game over
        return 0
    }

    for child != nil {
        childCount++
        eval = searchAlphaBeta(child, depth-1, true, alpha, beta)
        if eval < bestEval {
            bestEval = eval
        }
        if bestEval < beta {
            beta = bestEval
        }
        if alpha >= beta {
            break
        }
        child, childIndex = nextChildForOpponent(node, childIndex)
    }
    NumNodes += childCount
    return bestEval
}

func nextChildForPlayer(b board.Board, index uint8) (board.Board, uint8) {
    for index < 4 {
        if child := b.DoMove(index); ! child.Equals(b) {
            return child, index+1
        }
        index++
    }
    return nil, 0
}

func nextChildForOpponent(b board.Board, index uint8) (board.Board, uint8) {
    // for index < params.N {
    // TODO: replace 16 with params.N
    for index < 16 {
        if b[index] == 0 {
            return b.Insert(index, 1), index+1
        }
        index++
    }
    return nil, 0
}
