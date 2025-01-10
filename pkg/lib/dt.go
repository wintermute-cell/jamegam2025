package lib

import "github.com/hajimehoshi/ebiten/v2"

func Dt() float64 {
	tps := ebiten.ActualTPS()
	if tps <= 1/2000 { // The first frame might have ActualTPS 0.
		return 0.0
	}

	return 1.0 / tps
}
