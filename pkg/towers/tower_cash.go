package towers

import (
	"jamegam/pkg/audio"
	"jamegam/pkg/enemy"
	"jamegam/pkg/lib"
)

var _ Tower = &TowerCash{}

type TowerCash struct {
	*Towercore
}

func NewTowerCash(position lib.Vec2I) *TowerCash {
	tc := NewTowercore(3.0, 90.0, SpritesheetTowerCash, position)
	tc.spriteFrames = 8
	return &TowerCash{
		Towercore: tc,
	}
}

func (t *TowerCash) Price() int64 {
	return 100
}

// Update implements Tower.
func (t *TowerCash) Update(em EnemyManager, pm ProjectileManager) error {
	enemies, _ := em.GetEnemies(t.position.ToVec2().Add(lib.NewVec2(32, 32)), t.radius)
	hitEnemies := []*enemy.Enemy{} // only at max 8 enemies can be hit
	for i, e := range enemies {
		if i >= 8 {
			break
		}
		hitEnemies = append(hitEnemies, e)
	}

	var baseMana int32 = 1
	if t.ShouldFire(lib.Dt()) && len(hitEnemies) > 0 {
		mana := baseMana * (t.damageUpgrades + 1) * int32(len(hitEnemies))
		em.AddMana(int64(mana))
		// TODO: play sound
		audio.Controller.Play("tower_cash_shot", 0)
		t.shotThisTick = true
	}

	return nil
}
