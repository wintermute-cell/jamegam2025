package towers

import (
	"jamegam/pkg/lib"

	"github.com/hajimehoshi/ebiten/v2"
)

type Tower interface {
	Update(EnemyManager) error
	Draw(screen *ebiten.Image)
}

type EnemyManager interface {
	GetEnemies(point lib.Vec2, radius int) []Enemy
}

type Enemy interface {
	GetHealth() int
	SetHealth(int)
	GetSpeed() int
	SetSpeed(int)
	GetPosition() lib.Vec2
}
