/* Parse command-line arguments and initialize parameters.
These parameters define how the engine will be running.

NOTE: all commands and options for playing GUI, currently disabled */
package params

import (
    "flag"
)

var (
    /* Options for playing against GUI */
    //GameName string = "gnome-2048"
    //GamePath string = "/usr/bin/gnome-2048"
    //GameArg string = "--size=4"
    //GameXWindowName = "GNOME 2048"

    /* Xlib options when playing against GUI */
    // X Key event delay (milliseconds)
    //KeyDelay int = 50

    /* Options when playing against GUI or TUI */
    // Delay between each turn played (milliseconds)
    TurnDelay int = 250
    // Delay before starting game (seconds)
    StartDelay int = 5

    /* Board settings */

    // Length of rows/columns
    Size int = 4
    // Number of tiles (always Size^2)
    N int = 16

    /* Engine modes and options */
    // Play against simulated opponent (random moves)
    SelfPlay bool = false
    // Search best move for a given board (for debugging)
    SingleSearch bool = false

    // Engine to use for playing
    Engine string = "minmax"
    // Evaluation function
    Eval string = "score"
    // Maximum search depth (for MinMax)
    MaxDepth int = 7
    // Limit number of moves to play
    MaxMoves int = 1000000
    // In SelfPlay mode: number of games to play
    NumGames int = 1
    // Print more details during the game
    Verbose bool = false
)

// probably doesn't belong here
var Moves = []string{"Up", "Down", "Left", "Right"}

func Init() error {
    flag.IntVar(&Size, "bs", Size, "Board size")
    //flag.StringVar(&Engine, "engine", Engine, "Engine that will play the game")
    //flag.StringVar(&Eval, "eval", Eval, "Engine evaluation function")
    flag.IntVar(&MaxDepth, "sd", MaxDepth, "Default search depth")
    flag.IntVar(&MaxMoves, "mm", MaxMoves, "Limit number of moves to play")
    flag.IntVar(&TurnDelay, "td", TurnDelay, "Delay between each turn played (milliseconds)")
    flag.IntVar(&StartDelay, "sd", StartDelay, "Delay before starting to play (seconds)")
    flag.BoolVar(&Verbose, "v", Verbose, "Run in verbose mode")
    flag.BoolVar(&SelfPlay, "sp", SelfPlay, "Play with simulated opponent (rather than GUI)")
    flag.BoolVar(&SingleSearch, "ss", SingleSearch, "Run a single engine search")
    //flag.BoolVar(&EngineTest, "e", EngineTest, "Only test engine (no game)")

    flag.Parse()

    N = Size*Size

    /*
    if Size != 4 {
        // Match flag for gnome-2048
        GameArg = "--size="+strconv.Itoa(Size)
    }

    cmd := exec.Command("which", GameName)
    if out, err := cmd.Output(); err == nil {
        GamePath = strings.TrimSuffix(string(out), "\n")
        return nil
    }
    return errors.New("game ["+GameName+"] not installed")
    */
    return nil
}

