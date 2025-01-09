package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Game struct {
	gopher             *ebiten.Image
	gopherGeoM         *ebiten.GeoM
	firstClickHappened bool
	audioContext       *audio.Context
	audioPlayer        *audio.Player
}

func (g *Game) Update() error {
	if !g.firstClickHappened {
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			g.firstClickHappened = true
			ebiten.SetCursorMode(ebiten.CursorModeCaptured)
		}
		return nil
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		g.audioPlayer.Rewind()
		g.audioPlayer.Play()
	}

	dt := 1.0 / ebiten.ActualTPS()
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		g.gopherGeoM.Translate(-200*dt, 0)
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		g.gopherGeoM.Translate(200*dt, 0)
	}
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		g.gopherGeoM.Translate(0, -200*dt)
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		g.gopherGeoM.Translate(0, 200*dt)
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.Black)
	ebitenutil.DebugPrint(screen, "Bye, World!")
	screen.DrawImage(g.gopher, &ebiten.DrawImageOptions{
		GeoM: *g.gopherGeoM,
	})
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
	game.gopherGeoM = &ebiten.GeoM{}

	game.audioContext = audio.NewContext(44100)

	audioFile, err := ebitenutil.OpenFile("audio.ogg")
	if err != nil {
		panic(err)
	}
	audioTrack, err := vorbis.DecodeWithSampleRate(44100, audioFile)

	game.audioPlayer, err = game.audioContext.NewPlayer(audioTrack)

	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
