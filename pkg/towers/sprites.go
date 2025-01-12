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

	spritesheetTowerBasic *ebiten.Image
	spritesheetTowerTacks *ebiten.Image
	spritesheetTowerIce   *ebiten.Image
	spritesheetTowerAoe   *ebiten.Image
	spritesheetTowerCash  *ebiten.Image
	spritesheetTowerSuper *ebiten.Image
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

	spritesheetTowerBasic, _, err = ebitenutil.NewImageFromFile("sheet_4_towerbasic.png")
	lib.Must(err)

	spritesheetTowerTacks, _, err = ebitenutil.NewImageFromFile("sheet_4_towertacks.png")
	lib.Must(err)

	spritesheetTowerIce, _, err = ebitenutil.NewImageFromFile("sheet_4_towerice.png")
	lib.Must(err)

	spritesheetTowerAoe, _, err = ebitenutil.NewImageFromFile("sheet_5_toweraoe.png")
	lib.Must(err)

	spritesheetTowerCash, _, err = ebitenutil.NewImageFromFile("sheet_8_towercash.png")
	lib.Must(err)

	spritesheetTowerSuper, _, err = ebitenutil.NewImageFromFile("sheet_4_towersuper.png")
	lib.Must(err)
}
