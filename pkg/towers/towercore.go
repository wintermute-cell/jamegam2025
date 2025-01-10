package towers

import (
	"jamegam/pkg/lib"

	"github.com/hajimehoshi/ebiten/v2"
)

type Towercore struct {
	rof          float64
	sprite       *ebiten.Image
	position     lib.Vec2I
	drawPosition lib.Vec2

	lastFiredAgo float64
}

func NewTowercore(rof float64, sprite *ebiten.Image, position lib.Vec2I) *Towercore {
	ret := &Towercore{
		rof:      rof,
		sprite:   sprite,
		position: position,
	}

	ret.drawPosition = position.ToVec2() // TODO: for now, later some animation

	return ret
}

func (tc *Towercore) Draw(screen *ebiten.Image) {
	geom := ebiten.GeoM{}
	geom.Scale(4, 4)
	geom.Translate(float64(tc.drawPosition.X), float64(tc.drawPosition.Y))
	screen.DrawImage(tc.sprite, &ebiten.DrawImageOptions{GeoM: geom})
}

// ShouldFire must be called every tick to determine if the tower should fire
func (tc *Towercore) ShouldFire(dt float64) bool {
	if tc.lastFiredAgo >= tc.rof {
		tc.lastFiredAgo = 0
		return true
	}
	tc.lastFiredAgo += dt
	return false
}
