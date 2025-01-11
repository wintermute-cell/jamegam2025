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

	SpriteEnemyFast, _, err = ebitenutil.NewImageFromFile("test_enemyfast.png")
	lib.Must(err)

	SpriteEnemyTank, _, err = ebitenutil.NewImageFromFile("test_enemytank.png")
	lib.Must(err)

	// EFFECTS

	SpriteSlowEffect, _, err = ebitenutil.NewImageFromFile("test_effectslow.png")
	lib.Must(err)

	SpriteSpeedEffect, _, err = ebitenutil.NewImageFromFile("test_effectspeed.png")
	lib.Must(err)
}
