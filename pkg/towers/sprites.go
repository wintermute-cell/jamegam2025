package towers

import (
	"jamegam/pkg/lib"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var spriteTowerBasic *ebiten.Image
var spriteTowerSlowing *ebiten.Image
var spriteTowerAoe *ebiten.Image

func init() {
	var err error

	spriteTowerBasic, _, err = ebitenutil.NewImageFromFile("test_tower.png")
	lib.Must(err)

	spriteTowerSlowing, _, err = ebitenutil.NewImageFromFile("test_tower.png")
	lib.Must(err)

	spriteTowerAoe, _, err = ebitenutil.NewImageFromFile("test_tower.png")
	lib.Must(err)
}
