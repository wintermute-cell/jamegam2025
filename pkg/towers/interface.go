package towers

import (
	"jamegam/pkg/enemy"
	"jamegam/pkg/lib"

	"github.com/hajimehoshi/ebiten/v2"
)

type Tower interface {
	Update(EnemyManager) error
	Draw(screen *ebiten.Image)
}

type EnemyManager interface {
	GetEnemies(point lib.Vec2, radius float32) []*enemy.Enemy
}
