package entity

import "github.com/hajimehoshi/ebiten/v2"

// Ensure EntityEnemy implements Entity
var _ Entity = &EntityEnemy{}

type EntityEnemy struct {
}

func NewEntityEnemy() *EntityEnemy {
	newEnt := &EntityEnemy{}
	return newEnt
}

func (e *EntityEnemy) Init(EntitySpawner) {

}

func (e *EntityEnemy) Update(EntitySpawner) error {
	return nil
}

func (e *EntityEnemy) Deinit(EntitySpawner) {

}

func (e *EntityEnemy) Draw(screen *ebiten.Image) {

}
