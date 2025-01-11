package towers

import (
	"jamegam/pkg/enemy"
	"jamegam/pkg/lib"

	"github.com/hajimehoshi/ebiten/v2"
)

type Tower interface {
	Update(EnemyManager, ProjectileManager) error
	Draw(screen *ebiten.Image)
	Price() int64
	Radius() float32
}

type EnemyManager interface {
	GetEnemies(point lib.Vec2, radius float32) ([]*enemy.Enemy, []lib.Vec2I)
}

type ProjectileManager interface {
	AddProjectile(projectile Projectile) int
	RemoveProjectile(idx int)
}

type TowerType int

const (
	TowerTypeNone TowerType = iota
	TowerTypeBasic
	TowerTypeTacks
	TowerTypeIce
	TowerTypeAoe
)
