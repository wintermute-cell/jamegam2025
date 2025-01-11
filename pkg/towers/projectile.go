package towers

import (
	"image/color"
	"jamegam/pkg/lib"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Projectile interface {
	Update(em EnemyManager, pm ProjectileManager)
	Draw(screen *ebiten.Image)
}

// ========================================
// ProjectileBasic
// ========================================

var _ Projectile = &ProjectileBasic{}

type ProjectileBasic struct {
	direction lib.Vec2
	speed     float32
	position  lib.Vec2
	radius    float32
	SelfIdx   int
	lifetime  float32
	damage    int
}

func NewProjectileBasic(direction, position lib.Vec2, speed float32, radius float32) *ProjectileBasic {
	p := &ProjectileBasic{
		direction: direction,
		speed:     speed,
		position:  position,
		radius:    radius,
		damage:    1,
	}
	return p
}

func (p *ProjectileBasic) Update(em EnemyManager, pm ProjectileManager) {
	dt := float32(lib.Dt())
	offset := p.direction.Mul(p.speed * dt)

	// This might cause discrepancies in the future, but I hope that they're
	// small enough to be negligible in our case.
	p.position = p.position.Add(offset)

	p.lifetime = p.lifetime + dt
	if p.lifetime > 3 {
		pm.RemoveProjectile(p.SelfIdx)
	}

	// Check for collision with enemies
	enemies, _ := em.GetEnemies(p.position, p.radius)
	for _, e := range enemies {
		newHealth := e.GetHealth() - p.damage
		if newHealth <= 0 {
			pm.RemoveProjectile(p.SelfIdx)
		}
		e.SetHealth(newHealth)
	}

}

func (p *ProjectileBasic) Draw(screen *ebiten.Image) {
	// vector.DrawFilledCircle(screen, float32(p.shBounds.Mx), float32(p.shBounds.My), 20, color.RGBA{255, 255, 0, 255}, false)
	geom := ebiten.GeoM{}
	geom.Translate(float64(p.position.X), float64(p.position.Y))
	vector.DrawFilledCircle(screen, float32(p.position.X), float32(p.position.Y), 5, color.RGBA{255, 255, 0, 255}, false)
	// TODO: draw shbounds box
}

// ========================================
// Projectile...
// ========================================
