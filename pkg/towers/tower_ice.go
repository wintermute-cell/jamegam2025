package towers

import (
	"jamegam/pkg/audio"
	"jamegam/pkg/enemy"
	"jamegam/pkg/lib"
)

var _ Tower = &TowerIce{}

type TowerIce struct {
	*Towercore
}

func NewTowerIce(position lib.Vec2I) *TowerIce {
	tc := NewTowercore(2.0, 90.0, SpritesheetTowerIce, position)
	tc.animSpeed = 0.1
	return &TowerIce{
		Towercore: tc,
	}
}

func (t *TowerIce) Price() int64 {
	return 100
}

// Update implements Tower.
func (t *TowerIce) Update(em EnemyManager, pm ProjectileManager) error {
	enemies, _ := em.GetEnemies(t.position.ToVec2().Add(lib.NewVec2(32, 32)), t.radius)
	hitEnemies := []*enemy.Enemy{} // only at max 6 enemies can be hit
	for i, e := range enemies {
		if i >= 6 {
			break
		}
		hitEnemies = append(hitEnemies, e)
	}

	if t.ShouldFire(lib.Dt()) && len(hitEnemies) > 0 {
		// Spawn projectiles in a circle around the tower
		speedMod := float32(0.5 - (0.05 * float64(t.speedUpgrades+t.damageUpgrades)))
		for _, e := range hitEnemies {
			e.SetSpeedMod(speedMod, 2)
		}
		audio.Controller.Play("ice_tower_shoot", 0.05)
		t.shotThisTick = true
		// TODO: visual effect
	}

	return nil
}
