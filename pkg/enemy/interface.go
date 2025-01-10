package enemy

import "jamegam/pkg/lib"

type EnemyType int

const (
	EnemyTypeBasic EnemyType = iota
	EnemyTypeFast
	EnemyTypeTank
)

type Enemy interface {
	GetHealth() int
	SetHealth(int)
	GetSpeed() int
	SetSpeed(int)
	GetPosition() lib.Vec2
}
