package enemy

import "jamegam/pkg/lib"

type Enemy interface {
	GetHealth() int
	SetHealth(int)
	GetSpeed() int
	SetSpeed(int)
	GetPosition() lib.Vec2
}
