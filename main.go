package main

import (
	"github.com/OscarSierra24/Earthquake-Simulator/game"
	"github.com/OscarSierra24/Earthquake-Simulator/ui2d"
)

func main() {
	ui := &ui2d.UI2d{}
	game.Run(ui)
}