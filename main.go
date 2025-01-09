package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Game struct {
	gopher       *ebiten.Image
	gopherCursed *ebiten.Image
	goperOptions *ebiten.DrawImageOptions
	colorSet     bool
}

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		g.goperOptions.GeoM.Translate(0, -1)
	} else if ebiten.IsKeyPressed(ebiten.KeyA) {
		g.goperOptions.GeoM.Translate(-1, 0)
	} else if ebiten.IsKeyPressed(ebiten.KeyS) {
		g.goperOptions.GeoM.Translate(0, 1)
	} else if ebiten.IsKeyPressed(ebiten.KeyD) {
		g.goperOptions.GeoM.Translate(1, 0)
	}

	if ebiten.IsKeyPressed(ebiten.KeyQ) {
		g.colorSet = true
	} else {
		g.colorSet = false
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, "Bye, World!")
	screen.Fill(color.White)
	if g.colorSet {
		screen.DrawImage(g.gopherCursed, g.goperOptions)
	} else {
		screen.DrawImage(g.gopher, g.goperOptions)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Bye, World!")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeDisabled)
	game := &Game{}

	var err error
	game.gopher, _, err = ebitenutil.NewImageFromFile("gopher.png")
	if err != nil {
		panic(err)
	}

	game.gopherCursed, _, err = ebitenutil.NewImageFromFile("MyGopher.png")
	if err != nil {
		panic(err)
	}

	game.goperOptions = new(ebiten.DrawImageOptions)

	game.colorSet = false

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
