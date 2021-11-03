package main

import (
    "github.com/vgratian/gnome-2048-player/board"
    "github.com/vgratian/gnome-2048-player/engine"
    "github.com/vgratian/gnome-2048-player/options"
    "github.com/vgratian/gnome-2048-player/gui"
    "fmt"
    "os"
)

func main() {

    var (
        opts *options.Options
        eng engine.Engine
        brd *board.Board
        step uint64
        err error
    )

    opts = options.ParseOrExit()

    fmt.Printf("Starting Game with options: %s\n", opts.String())

    if eng, err = engine.GetEngine(opts); err != nil {
        exitWithError("GetEngine", err)
    }

    if err = ui.Init(opts); err != nil {
        exitWithError("UI Init", err)
    }

    brd = board.New()

    if err = ui.UpdateBoard(brd); err != nil {
        exitWithError("UpdateBoard", err)
    }

    fmt.Printf(" -------- -------- -------- --------\n")
    fmt.Printf("     step    depth     move    score\n")
    fmt.Printf(" -------- -------- -------- --------\n")

    for ! brd.GameOver() {

        move, depth := eng.Search(brd)

        ui.DoMove(move)

        ui.UpdateBoard(brd)

        fmt.Printf(" %8d %8d %8s %8d\n", step++, depth, move, brd.Score())

    }

    fmt.Println()
    fmt.Println(" ** GAME OVER ** ")

}

func ExitWithError(s string, e error) {
    fmt.Printf("Err (%s): %s\n", s, e.Error())
    os.Exit(1)
}

