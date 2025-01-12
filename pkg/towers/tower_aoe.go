package towers

import (
	"jamegam/pkg/audio"
	"jamegam/pkg/enemy"
	"jamegam/pkg/lib"
)

var _ Tower = &TowerAoe{}

type TowerAoe struct {
	*Towercore
}

func NewTowerAoe(position lib.Vec2I) *TowerAoe {
	tc := NewTowercore(3.0, 195.0, SpritesheetTowerAoe, position)
	tc.animSpeed = 0.20
	tc.spriteFrames = 5
	return &TowerAoe{
		Towercore: tc,
	}
}

func (t *TowerAoe) Price() int64 {
	return 250
}

// Update implements Tower.
func (t *TowerAoe) Update(em EnemyManager, pm ProjectileManager) error {
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
		// t.lookAt = dirToEnemy
	}

	// TODO: if there is an ememy in range...
	if t.ShouldFire(lib.Dt()) && furthestEnemy != nil {
		lastIdx, nextIdx := furthestEnemy.GetPathNodes()
		last := path[lastIdx].ToVec2().Mul(64)
		next := path[nextIdx].ToVec2().Mul(64)
		pos := last.Lerp(next, float32(furthestEnemy.GetPathProgress()))
		dirToEnemy = pos.Sub(t.position.ToVec2()).Normalize()
		prj := NewProjectileExplosive(
			dirToEnemy,
			t.position.ToVec2().Add(lib.NewVec2(32, 32)),
			550.0,
			12.0,
			0.45,
			50,
			int(1*t.damageUpgrades+1),
		)
		idx := pm.AddProjectile(prj)
		prj.SelfIdx = idx
		audio.Controller.Play("aoe_tower_shoot", 0)
		t.shotThisTick = true
	}

	return nil
}
