package main

import (
    "github.com/vgratian/2048/params"
    "github.com/vgratian/2048/engine"
)

func main() {

    params.Init()

    if params.SingleSearch() {
        engine.SingleSearch()
    } else { // default mode
        engine.SelfPlay()
    }
}
