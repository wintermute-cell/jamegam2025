package towers

import (
	"jamegam/pkg/lib"
	"log"
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
	if t.ShouldFire(lib.Dt()) {
		// TODO:
		// panic("unimplemented")
		log.Println("TowerBasic fired!")
	}

	return nil
}
