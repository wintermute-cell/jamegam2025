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
	direction   lib.Vec2
	speed       float32
	position    lib.Vec2
	radius      float32
	SelfIdx     int
	lifetime    float32
	maxLifetime float32
	damage      int
}

func NewProjectileBasic(direction, position lib.Vec2, speed float32, radius float32, maxLifetime float32) *ProjectileBasic {
	p := &ProjectileBasic{
		direction:   direction,
		speed:       speed,
		position:    position,
		radius:      radius,
		maxLifetime: maxLifetime,
		damage:      1,
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
	if p.lifetime > p.maxLifetime {
		pm.RemoveProjectile(p.SelfIdx)
	}

	// Check for collision with enemies
	enemies, _ := em.GetEnemies(p.position, p.radius)
	for _, e := range enemies {
		newHealth := e.GetHealth() - p.damage
		e.SetHealth(newHealth)
		pm.RemoveProjectile(p.SelfIdx)
		return
		// if newHealth <= 0 {
		// }
	}

}

func (p *ProjectileBasic) Draw(screen *ebiten.Image) {
	geom := ebiten.GeoM{}
	geom.Translate(float64(p.position.X), float64(p.position.Y))
	vector.DrawFilledCircle(screen, float32(p.position.X), float32(p.position.Y), 5, color.RGBA{255, 255, 0, 255}, false)
	// TODO: draw shbounds box
}

// ========================================
// ProjectileExplosive
// ========================================

var _ Projectile = &ProjectileExplosive{}

type ProjectileExplosive struct {
	direction       lib.Vec2
	speed           float32
	position        lib.Vec2
	radius          float32
	SelfIdx         int
	lifetime        float32
	maxLifetime     float32
	explosionRadius float32
	damage          int
	exploding       bool
	explodingTimer  float32
}

func NewProjectileExplosive(direction, position lib.Vec2, speed float32, radius float32, maxLifetime float32, explosionRadius float32) *ProjectileExplosive {
	p := &ProjectileExplosive{
		direction:       direction,
		speed:           speed,
		position:        position,
		radius:          radius,
		maxLifetime:     maxLifetime,
		explosionRadius: explosionRadius,
		damage:          1,
	}
	return p
}

func (p *ProjectileExplosive) Update(em EnemyManager, pm ProjectileManager) {

	if p.exploding {
		p.explodingTimer += float32(lib.Dt())
		if p.explodingTimer > 0.3 {
			pm.RemoveProjectile(p.SelfIdx)
		}
		return
	}

	dt := float32(lib.Dt())
	offset := p.direction.Mul(p.speed * dt)

	// This might cause discrepancies in the future, but I hope that they're
	// small enough to be negligible in our case.
	p.position = p.position.Add(offset)

	p.lifetime = p.lifetime + dt
	if p.lifetime > p.maxLifetime {
		pm.RemoveProjectile(p.SelfIdx)
	}

	// Check for collision with enemies
	enemies, _ := em.GetEnemies(p.position, p.radius)
	if len(enemies) == 0 {
		return
	}

	explodedEnemies, _ := em.GetEnemies(p.position, p.explosionRadius)
	for _, e := range explodedEnemies {
		newHealth := e.GetHealth() - p.damage
		e.SetHealth(newHealth)
	}
	p.exploding = true

}

func (p *ProjectileExplosive) Draw(screen *ebiten.Image) {
	geom := ebiten.GeoM{}
	geom.Translate(float64(p.position.X), float64(p.position.Y))

	if p.exploding {
		vector.DrawFilledCircle(screen, float32(p.position.X), float32(p.position.Y), p.explosionRadius, color.RGBA{255, 0, 0, 150}, false)
		return
	}
	vector.DrawFilledCircle(screen, float32(p.position.X), float32(p.position.Y), 10, color.RGBA{50, 50, 50, 255}, false)
	// TODO: draw shbounds box
}
