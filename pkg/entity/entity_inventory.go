package entity

import "github.com/hajimehoshi/ebiten/v2"

// Ensure EntityInventory implements Entity
var _ Entity = &EntityInventory{}

type EntityInventory struct {
}

func NewEntityInventory() *EntityInventory {
	newEnt := &EntityInventory{}
	return newEnt
}

func (e *EntityInventory) Init(EntitySpawner) {

}

func (e *EntityInventory) Update(EntitySpawner) error {
	return nil
}

func (e *EntityInventory) Deinit(EntitySpawner) {

}

func (e *EntityInventory) Draw(screen *ebiten.Image) {

}
