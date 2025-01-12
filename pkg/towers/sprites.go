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

	SpritesheetTowerBasic *ebiten.Image
	SpritesheetTowerTacks *ebiten.Image
	SpritesheetTowerIce   *ebiten.Image
	SpritesheetTowerAoe   *ebiten.Image
	SpritesheetTowerCash  *ebiten.Image
	SpritesheetTowerSuper *ebiten.Image

	SpriteProjectileBasic *ebiten.Image
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

	SpritesheetTowerBasic, _, err = ebitenutil.NewImageFromFile("sheet_4_towerbasic.png")
	lib.Must(err)

	SpritesheetTowerTacks, _, err = ebitenutil.NewImageFromFile("sheet_4_towertacks.png")
	lib.Must(err)

	SpritesheetTowerIce, _, err = ebitenutil.NewImageFromFile("sheet_4_towerice.png")
	lib.Must(err)

	SpritesheetTowerAoe, _, err = ebitenutil.NewImageFromFile("sheet_5_toweraoe.png")
	lib.Must(err)

	SpritesheetTowerCash, _, err = ebitenutil.NewImageFromFile("sheet_8_towercash.png")
	lib.Must(err)

	SpritesheetTowerSuper, _, err = ebitenutil.NewImageFromFile("sheet_4_towersuper.png")
	lib.Must(err)

	SpriteProjectileBasic, _, err = ebitenutil.NewImageFromFile("projectile_basic.png")
	lib.Must(err)
}
