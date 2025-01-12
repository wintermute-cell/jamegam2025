package lib

import "github.com/hajimehoshi/ebiten/v2"

func Dt() float64 {
	tps := ebiten.ActualTPS()
	if tps <= 1/2000 { // The first frame might have ActualTPS 0.
		return 0.0
	}

	realdt := 1.0 / tps
	if realdt > 1 { // prevent huge spikes if the game is paused
		return 1.0 / 60.0
	}

	return realdt / 1.0
}
