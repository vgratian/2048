package tui

import (
	"fmt"
	"strconv"
	"strings"
    "github.com/vgratian/2048/params"
)

// most if special characters, from:
// https://espterm.github.io/docs/VT100%20escape%20codes.html

const (
	jumpUp     = "\033[A"
	jumpDown  = "\033[B"
	jumpLeft   = "\033[D"
	jumpRight  = "\033[C"
    jumpHome   = "\033[H"
	clearLine  = "\033[K"
    clearAll   = "\033[J"
	resetColor = "\033[0m"
    resetTerm  = "\033[c"

    Bold = "\033[1m"
    End = "\033[0m"
    Red = "\033[31m"
    Pink = "\033[35m"
    Yellow = "\033[93m"
    Cyan = "\033[36m"
    Grey = "\033[90m"

)

var BgColors = []string{
        // colors that are close to GUI
		"\033[107m", //   0 - white
		"\033[103m", //   2 - yellow
		"\033[42m",  //   4 - green
		"\033[41m",  //   8 - orange? => red
		"\033[44m",  //  16 - blue
		"\033[105m", //  32 - pink
		"\033[40m",  //  64 - brown? => black
		"\033[45m",  // 128 - red
		"\033[46m",  // 256 - light brown? => light blue
		// random colors from here on
		"\033[47m",  //   512
		"\033[48m",  //  1024
		"\033[101m", //  2048
		"\033[102m", //  4096
		"\033[104m", //  8192
		"\033[100m", // 16384
		"\033[106m", // 32786
}

var (
	indentFmt  string = "%2s"
	tileDigit  string = "%s%6d" + resetColor
	tileString string = "%s%6s" + resetColor
    tileWidth int = 6
)

func SetIndent(size int) {
	indentFmt = "%" + strconv.Itoa(size) + "s"
}

func SetTileWidth(size int) {
    tileWidth = size
	tileDigit = "%s%" + strconv.Itoa(size) + "d" + resetColor
	tileString = "%s%" + strconv.Itoa(size) + "s" + resetColor
}

func IntDecor(x uint64) string {
    s := ""
    f := fmt.Sprintf("%d", x)
    i := len(f)
    for i>0 {
        k := i-3
        if k < 0 { k = 0 }
        s = f[k:i] + "," + s
        i = k
    }
    return s[:len(s)-1]
}

func RenderEval(eval [4]float32, move uint8) string {
    s := ""
    for i := range eval {
        if i == int(move) {
            s += fmt.Sprintf(" %s%6.3f%s", Bold, eval[i], End)
        } else {
            s += fmt.Sprintf(" %6.3f", eval[i])
        }
    }
    return s
}

/*
const (
    _width_step    = 8
    _width_score   = 10
    _width_depth   = 8
    _width_nodes   = 12
    _width_move    = 10
    _width_eval    = 35
    _width_time    = 20
)*/

func PrintScoreHeader() {
    fmt.Printf(" -------- ---------- -------- --------------- ---------- ---------------------------------- -------------------- \n")
    fmt.Printf("     step      score    depth           nodes       move          eval (up-down-left-right)                 time \n")
    fmt.Printf(" -------- ---------- -------- --------------- ---------- ---------------------------------- -------------------- \n")
}

func PrintScore0(step, score, depth uint32, nodes uint64, move uint8, eval [4]float32, t float64) {
    fmt.Printf(" %8d %10d %8d %15s %10s %35s %20.5f\n", step, score, depth, IntDecor(nodes), params.Moves[move], RenderEval(eval, move), t)
}

func PrintScore(step, score, depth uint32, nodes uint64, move uint8, eval [4]float32) {
    fmt.Printf(" %8d %10d %8d %15s %10s %35s\n", step, score, depth, IntDecor(nodes), params.Moves[move], RenderEval(eval, move))
}

func PrintScore2(step, score uint32, move uint8) {
    fmt.Printf(" %s%8d %10d %8s %15s %10d%s\n", Grey, step, score, "-", "-", move, End)
}

func PrintBoard(board []uint8, title string) {

    fmt.Print(jumpDown)
    fmt.Print(clearLine)
    fmt.Print(jumpUp)
	//fmt.Print(clearLine)
    //fmt.Print(jumpDown)
    //fmt.Print(clearLine)
	//fmt.Println()
	//fmt.Print(clearLine)
	//fmt.Printf(indentFmt, "")
	//fmt.Printf(title)

	for i, v := range board {

        tileColor := BgColors[v]

		if i % params.Size == 0 {
			fmt.Println(jumpDown)
			fmt.Print(clearLine)
			fmt.Printf(indentFmt, "")
			fmt.Printf(tileString, tileColor, "")
			fmt.Println()
			fmt.Printf(indentFmt, "")
		} else {
			fmt.Printf(tileString, tileColor, "")
			fmt.Print(jumpDown)
			fmt.Print(strings.Repeat(jumpLeft, tileWidth))
		}

        if v == 0 {
    		fmt.Printf(tileDigit, tileColor, 0)
        } else {
    		fmt.Printf(tileDigit, tileColor, 1 << v)
        }

		fmt.Print(jumpUp)
	}
}

func ClearScreen() {
    fmt.Print(jumpHome)
    fmt.Print(clearAll)
}

func ResetTerm() {
    fmt.Print(resetTerm)
}

func ScrollDown() {
	fmt.Print(strings.Repeat(jumpDown, params.Size*2+3))
	fmt.Println()
}

func ScrollBack() {
	//fmt.Print(strings.Repeat(jumpUp, params.Size*2+1))
	fmt.Print(strings.Repeat(jumpUp, params.Size*2))
	fmt.Print(strings.Repeat(jumpLeft, 100))
	fmt.Print(clearLine)
}
