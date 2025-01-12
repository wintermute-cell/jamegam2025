package enemy

import (
	"jamegam/pkg/lib"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var (
	SpriteEnemyBasic *ebiten.Image
	SpriteEnemyFast  *ebiten.Image
	SpriteEnemyTank  *ebiten.Image

	SpriteEnemyBasicSheet *ebiten.Image
	SpriteEnemyFastSheet  *ebiten.Image
	SpriteEnemyTankSheet  *ebiten.Image

	SpriteEnemyPoofSheet *ebiten.Image
)

var (
	SpriteSlowEffect  *ebiten.Image
	SpriteSpeedEffect *ebiten.Image
)

func init() {
	var err error

	// ENEMIES

	SpriteEnemyBasic, _, err = ebitenutil.NewImageFromFile("test_enemy.png")
	lib.Must(err)
	SpriteEnemyBasicSheet, _, err = ebitenutil.NewImageFromFile("sheet_4_rat.png")
	lib.Must(err)

	SpriteEnemyFast, _, err = ebitenutil.NewImageFromFile("test_enemyfast.png")
	lib.Must(err)
	SpriteEnemyFastSheet, _, err = ebitenutil.NewImageFromFile("sheet_5_bat.png")
	lib.Must(err)

	SpriteEnemyTank, _, err = ebitenutil.NewImageFromFile("test_enemytank.png")
	lib.Must(err)
	SpriteEnemyTankSheet, _, err = ebitenutil.NewImageFromFile("sheet_4_zombie.png")
	lib.Must(err)

	SpriteEnemyPoofSheet, _, err = ebitenutil.NewImageFromFile("sheet_3_poof.png")
	lib.Must(err)

	// EFFECTS

	SpriteSlowEffect, _, err = ebitenutil.NewImageFromFile("test_effectslow.png")
	lib.Must(err)

	SpriteSpeedEffect, _, err = ebitenutil.NewImageFromFile("test_effectspeed.png")
	lib.Must(err)
}
