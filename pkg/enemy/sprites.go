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

func init() {
	var err error

	SpriteEnemyBasic, _, err = ebitenutil.NewImageFromFile("test_enemy.png")
	lib.Must(err)

	SpriteEnemyFast, _, err = ebitenutil.NewImageFromFile("test_enemyfast.png")
	lib.Must(err)

	SpriteEnemyTank, _, err = ebitenutil.NewImageFromFile("test_enemytank.png")
	lib.Must(err)
}
