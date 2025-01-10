package towers

import (
	"jamegam/pkg/enemy"
	"jamegam/pkg/lib"
	"log"
)

var _ Tower = &TowerBasic{}

type TowerBasic struct {
	*Towercore
}

func NewTowerBasic(position lib.Vec2I) *TowerBasic {
	return &TowerBasic{
		Towercore: NewTowercore(1.0, 100.0, spriteTowerBasic, position),
	}
}

// Update implements Tower.
func (t *TowerBasic) Update(em EnemyManager) error {
	enemies := em.GetEnemies(t.position.ToVec2(), t.radius)

	log.Printf("enemies: %v", len(enemies))
	var furthestProgress float64 = -1
	var furthestEnemy *enemy.Enemy
	for _, e := range enemies {
		prog := e.GetNumPassedNodes() + e.GetPathProgress()
		if prog > furthestProgress {
			furthestProgress = prog
			furthestEnemy = e
		}
	}

	if t.ShouldFire(lib.Dt()) && false { // TODO: remove false
		// TODO:
		// panic("unimplemented")
		log.Println("Firing at enemy %v", furthestEnemy)
	}

	return nil
}
