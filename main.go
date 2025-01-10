package main

import (
	"jamegam/pkg/game"

	"github.com/hajimehoshi/ebiten/v2"
)

func configure() {
	ebiten.SetWindowSize(1024, 768)
	ebiten.SetWindowTitle("Bye, World!")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeDisabled)
}

func main() {
	configure()
	game := game.NewGame()
	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
