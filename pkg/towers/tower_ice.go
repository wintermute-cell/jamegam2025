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
	return &TowerIce{
		Towercore: NewTowercore(1.0, 90.0, spriteTowerIce, position),
	}
}

func (t *TowerIce) Price() int64 {
	return 100
}

// Update implements Tower.
func (t *TowerIce) Update(em EnemyManager, pm ProjectileManager) error {
	enemies, _ := em.GetEnemies(t.position.ToVec2(), t.radius)
	hitEnemies := []*enemy.Enemy{} // only at max 8 enemies can be hit
	for i, e := range enemies {
		if i >= 8 {
			break
		}
		hitEnemies = append(hitEnemies, e)
	}

	// TODO: if there is an ememy in range...
	if len(hitEnemies) > 0 && t.ShouldFire(lib.Dt()) {
		// Spawn projectiles in a circle around the tower
		speedMod := float32(0.5 - (0.05 * float64(t.speedUpgrades+t.damageUpgrades)))
		for _, e := range hitEnemies {
			e.SetSpeedMod(speedMod, 2)
		}
		audio.Controller.Play("test_pew")
	}

	return nil
}
