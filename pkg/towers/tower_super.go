package towers

import (
	"jamegam/pkg/audio"
	"jamegam/pkg/enemy"
	"jamegam/pkg/lib"
)

var _ Tower = &TowerSuper{}

type TowerSuper struct {
	*Towercore
}

func NewTowerSuper(position lib.Vec2I) *TowerSuper {
	return &TowerSuper{
		Towercore: NewTowercore(0.2, 128.0, SpritesheetTowerSuper, position),
	}
}

func (t *TowerSuper) Price() int64 {
	return 100
}

// Update implements Tower.
func (t *TowerSuper) Update(em EnemyManager, pm ProjectileManager) error {
	enemies, path := em.GetEnemies(t.position.ToVec2().Add(lib.NewVec2(32, 32)), t.radius)
	var furthestProgress float64 = -1
	var furthestEnemy *enemy.Enemy
	for _, e := range enemies {
		prog := e.GetNumPassedNodes() + e.GetPathProgress()
		if prog > furthestProgress {
			furthestProgress = prog
			furthestEnemy = e
		}
	}

	dirToEnemy := lib.NewVec2(0, 0)
	if furthestEnemy != nil {
		lastIdx, nextIdx := furthestEnemy.GetPathNodes()
		last := path[lastIdx].ToVec2().Mul(64)
		next := path[nextIdx].ToVec2().Mul(64)
		pos := last.Lerp(next, float32(furthestEnemy.GetPathProgress()))
		dirToEnemy = pos.Sub(t.position.ToVec2()).Normalize()
		t.lookAt = dirToEnemy
	}

	if t.ShouldFire(lib.Dt()) && furthestEnemy != nil {
		prj := NewProjectileBasic(
			dirToEnemy,
			t.position.ToVec2().Add(lib.NewVec2(32, 32)),
			800.0,
			12.0,
			0.3,
			int(1*t.damageUpgrades+1),
		)
		idx := pm.AddProjectile(prj)
		prj.SelfIdx = idx
		audio.Controller.Play("basic_tower_shoot", 0.05)
		t.shotThisTick = true
	}

	return nil
}
