package towers

import (
	"jamegam/pkg/audio"
	"jamegam/pkg/enemy"
	"jamegam/pkg/lib"
)

var _ Tower = &TowerTacks{}

type TowerTacks struct {
	*Towercore
}

func NewTowerTacks(position lib.Vec2I) *TowerTacks {
	return &TowerTacks{
		Towercore: NewTowercore(1.0, 90.0, spriteTowerTacks, position),
	}
}

func (t *TowerTacks) Price() int64 {
	return 100
}

// Update implements Tower.
func (t *TowerTacks) Update(em EnemyManager, pm ProjectileManager) error {
	enemies, _ := em.GetEnemies(t.position.ToVec2().Add(lib.NewVec2(32, 32)), t.radius)
	var furthestProgress float64 = -1
	var furthestEnemy *enemy.Enemy
	for _, e := range enemies {
		prog := e.GetNumPassedNodes() + e.GetPathProgress()
		if prog > furthestProgress {
			furthestProgress = prog
			furthestEnemy = e
		}
	}

	if t.ShouldFire(lib.Dt()) && furthestEnemy != nil {
		// Spawn projectiles in a circle around the tower
		for i := 0; i < 8; i++ {
			angle := float32(i) * 45
			dirToEnemy := lib.NewVec2(1, 0).Rotate(angle)
			prj := NewProjectileBasic(
				dirToEnemy,
				t.position.ToVec2().Add(lib.NewVec2(32, 32)),
				800.0,
				12.0,
				0.13,
				int(1*t.damageUpgrades+1),
			)
			idx := pm.AddProjectile(prj)
			prj.SelfIdx = idx
		}
		audio.Controller.Play("test_pew", 0)
	}

	return nil
}
