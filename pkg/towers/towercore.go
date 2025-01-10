package towers

import (
	"image/color"
	"jamegam/pkg/lib"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Towercore struct {
	rof          float64
	radius       float32
	sprite       *ebiten.Image
	position     lib.Vec2I
	drawPosition lib.Vec2

	lastFiredAgo float64
}

func NewTowercore(rof float64, radius float32, sprite *ebiten.Image, position lib.Vec2I) *Towercore {
	ret := &Towercore{
		rof:      rof,
		radius:   radius,
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
	vector.StrokeCircle(
		screen,
		float32(tc.drawPosition.X+32),
		float32(tc.drawPosition.Y+32),
		float32(tc.radius)+64,
		1.0,
		color.RGBA{255, 0, 255, 255},
		false)
	vector.StrokeCircle(
		screen,
		float32(tc.drawPosition.X+32),
		float32(tc.drawPosition.Y+32),
		float32(tc.radius),
		1.0,
		color.RGBA{255, 0, 255, 255},
		false)
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
