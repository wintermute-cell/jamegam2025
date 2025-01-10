package entity

import "github.com/hajimehoshi/ebiten/v2"

// Entity is an interface for game entities.
type Entity interface {
	Init(EntitySpawner)
	Update(EntitySpawner) error
	Deinit(EntitySpawner)
	Draw(screen *ebiten.Image)
}

// EntitySpawner is an interface for adding and removing entities to/from the
// game.
type EntitySpawner interface {
	AddEntity(e Entity)
	RemoveEntity(e Entity)
}
