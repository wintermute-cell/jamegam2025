package towers

import (
	"jamegam/pkg/lib"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var (
	spriteTowerBasic *ebiten.Image
	spriteTowerTacks *ebiten.Image
	spriteTowerIce   *ebiten.Image
	spriteTowerAoe   *ebiten.Image
	spriteTowerCash  *ebiten.Image
)

func init() {
	var err error

	spriteTowerBasic, _, err = ebitenutil.NewImageFromFile("test_tower.png")
	lib.Must(err)

	spriteTowerTacks, _, err = ebitenutil.NewImageFromFile("test_towertacks.png")
	lib.Must(err)

	spriteTowerIce, _, err = ebitenutil.NewImageFromFile("test_towerice.png")
	lib.Must(err)

	spriteTowerAoe, _, err = ebitenutil.NewImageFromFile("test_toweraoe.png")
	lib.Must(err)

	spriteTowerCash, _, err = ebitenutil.NewImageFromFile("test_towercash.png")
	lib.Must(err)
}
