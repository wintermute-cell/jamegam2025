package towers

import (
	"jamegam/pkg/lib"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

var _ Tower = &TowerBasic{}

type TowerBasic struct {
	*Towercore
}

func NewTowerBasic(position lib.Vec2I) *TowerBasic {
	return &TowerBasic{
		Towercore: NewTowercore(1.0, spriteTowerBasic, position),
	}
}

// Update implements Tower.
func (t *TowerBasic) Update(EnemyManager) error {
	dt := 1.0 / ebiten.ActualTPS()

	if t.ShouldFire(dt) {
		// TODO:
		// panic("unimplemented")
		log.Println("TowerBasic fired!")
	}

	return nil
}
