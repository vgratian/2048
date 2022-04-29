
## 2048 Engine

This is a 2048-playing program for educational purposes. Currently it implements:
* the minimax algorithm
* alpha-beta pruning

It plays a simulated-randomized opponent and draws colorized board on your console.

Work-in-progress:
* re-inforcment learning
* play agains a GUI-program ([gnome-2048](https://github.com/GNOME/gnome-2048))

### Requirements
**Required:**
* Go >= 1.17
* Terminal emulator with colors

**Go libraries:**
* none

**Optional:**
* Playing against GUI requires [libX11](https://gitlab.freedesktop.org/xorg/lib/libx11).

### Build and usage
Build with go:
```sh
go build main.go -o 2048
```

Run `./2048` to start the engine with default parameters or run `./2048 -help` for available options.
