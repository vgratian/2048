package engine

import "github.com/vgratian/2048/board"

const (
    MIN_EVAL_F32 float32 = 0.0
    MAX_EVAL_F32 float32 = 999999.0
    MIN_EVAL_U8 uint8 = 0
    MAX_EVAL_U8 uint8 = 255
)

// evaluation metric: sum of values
// penalized by the number of non-empty tiles
func eval(b board.Board) float32 {
    var sum, count uint32
    for _, v := range b {
        if v != 0 {
            sum += 1 << v
            count++
        }
    }
    return float32(sum) / float32(count)
}

func eval_u8_1(b board.Board) uint8 {
    var count uint8
    for _, v := range b {
        if v == 0 {
            count++
        }
    }
    return count
}

func eval_u8_2(b board.Board) uint8 {
    var sum, count uint8
    for _, v := range b {
        if v != 0 {
            sum += v
            count++
        }
    }
    return sum - count
}
