package towers

import (
	"jamegam/pkg/audio"
	"jamegam/pkg/enemy"
	"jamegam/pkg/lib"
)

var _ Tower = &TowerBasic{}

type TowerBasic struct {
	*Towercore
}

func NewTowerBasic(position lib.Vec2I) *TowerBasic {
	return &TowerBasic{
		Towercore: NewTowercore(1.0, 128.0, spriteTowerBasic, position),
	}
}

func (t *TowerBasic) Price() int64 {
	return 100
}

// Update implements Tower.
func (t *TowerBasic) Update(em EnemyManager, pm ProjectileManager) error {
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

	if t.ShouldFire(lib.Dt()) && furthestEnemy != nil {
		lastIdx, nextIdx := furthestEnemy.GetPathNodes()
		last := path[lastIdx].ToVec2().Mul(64)
		next := path[nextIdx].ToVec2().Mul(64)
		pos := last.Lerp(next, float32(furthestEnemy.GetPathProgress()))
		dirToEnemy := pos.Sub(t.position.ToVec2()).Normalize()
		prj := NewProjectileBasic(
			dirToEnemy,
			t.position.ToVec2().Add(lib.NewVec2(32, 32)),
			800.0,
			12.0,
			1,
			int(1*t.damageUpgrades+1),
		)
		idx := pm.AddProjectile(prj)
		prj.SelfIdx = idx
		audio.Controller.Play("test_pew", 0)
	}

	return nil
}
