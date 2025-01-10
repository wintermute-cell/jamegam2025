package entity

import (
	"jamegam/pkg/lib"

	"github.com/hajimehoshi/ebiten/v2"
)

// Ensure EntityTower implements Entity
var _ Entity = &EntityTower{}

type TowerStats struct {
	Range    int
	Damage   int
	FireRate int
}

type TowerType int

const (
	TowerTypeBasic TowerType = iota
	TowerTypeSlowdown
	TowerTypeAoe
)

var towerStats = map[TowerType]TowerStats{
	TowerTypeBasic: {5, 1, 1},
}

type EntityTower struct {
	towerType     TowerType
	towerImage    *ebiten.Image
	towerPosition lib.Vec2

	towerVisualPosition lib.Vec2
}

func NewEntityTower(towerType TowerType, position lib.Vec2, towerImage *ebiten.Image) *EntityTower {
	newEnt := &EntityTower{
		towerType:     towerType,
		towerPosition: position,
		towerImage:    towerImage,
	}
	newEnt.towerVisualPosition = position // TODO: for now, later maybe "fly in" animation
	return newEnt
}

func (e *EntityTower) Init(EntitySpawner) {}

func (e *EntityTower) Update(EntitySpawner) error {
	return nil
}

func (e *EntityTower) Deinit(EntitySpawner) {}

func (e *EntityTower) Draw(screen *ebiten.Image) {
	geom := ebiten.GeoM{}
	geom.Scale(4, 4)
	geom.Translate(float64(e.towerVisualPosition.X), float64(e.towerVisualPosition.Y))
	screen.DrawImage(e.towerImage, &ebiten.DrawImageOptions{GeoM: geom})
}
