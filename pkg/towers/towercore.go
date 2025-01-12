package towers

import (
	"image"
	"jamegam/pkg/lib"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type Towercore struct {
	rof            float64
	radius         float32
	sprite         *ebiten.Image
	position       lib.Vec2I
	drawPosition   lib.Vec2
	speedUpgrades  int32
	damageUpgrades int32

	tempSpeedBuff       float32
	tempDamageBuff      float32
	tempSpeedBuffTimer  time.Time
	tempDamageBuffTimer time.Time

	lastFiredAgo float64

	shotThisTick     bool
	spriteFrames     int
	spriteSheetIdx   int
	spriteSheetTimer float64
	isAnimating      bool

	lookAt lib.Vec2

	animSpeed float64

	settled    bool
	settleAnim float64
}

func NewTowercore(rof float64, radius float32, sprite *ebiten.Image, position lib.Vec2I) *Towercore {
	ret := &Towercore{
		rof:            rof,
		radius:         radius,
		sprite:         sprite,
		spriteFrames:   4,
		position:       position,
		speedUpgrades:  0,
		damageUpgrades: 0,
		animSpeed:      0.06,
		settleAnim:     -0.2,
		lastFiredAgo:   100,
		lookAt:         lib.Vec2{X: 0, Y: 1},
	}

	ret.drawPosition = position.ToVec2() // TODO: for now, later some animation

	return ret
}

func (tc *Towercore) Radius() float32 {
	return tc.radius
}

func (tc *Towercore) SetSpeedBuff(buff float32, duration float32) {
	tc.tempSpeedBuff = buff
	tc.tempSpeedBuffTimer = time.Now().Add(time.Duration(duration) * time.Second)
}

func (tc *Towercore) SetDamageBuff(buff float32, duration float32) {
	tc.tempDamageBuff = buff
	tc.tempDamageBuffTimer = time.Now().Add(time.Duration(duration) * time.Second)
}

func (tc *Towercore) GetTotalUpgrades() int32 {
	return tc.speedUpgrades + tc.damageUpgrades
}

func (tc *Towercore) GetSpeedUpgrades() int32 {
	return tc.speedUpgrades
}

func (tc *Towercore) GetDamageUpgrades() int32 {
	return tc.damageUpgrades
}

func (tc *Towercore) SpeedUpgrade() {
	tc.speedUpgrades++
}

func (tc *Towercore) DamageUpgrade() {
	tc.damageUpgrades++
}

func (tc *Towercore) Draw(screen *ebiten.Image) {
	geom := ebiten.GeoM{}
	geom.Translate(-8, -8)
	geom.Rotate(float64(-tc.lookAt.Angle()))
	geom.Translate(8, 8)
	geom.Scale(4, 4)
	geom.Translate(float64(tc.drawPosition.X), float64(tc.drawPosition.Y))
	// screen.DrawImage(tc.sprite, &ebiten.DrawImageOptions{GeoM: geom})

	fps := ebiten.ActualFPS()
	dt := 1.0 / 60.0
	if fps > 1/2000 {
		dt = 1.0 / fps
	}

	if !tc.settled {
		tc.settleAnim += 6.7 * dt
		geom.Translate(0, 12*math.Sin(tc.settleAnim))
		if tc.settleAnim > math.Pi {
			tc.settled = true
		}
	}

	if tc.shotThisTick {
		tc.shotThisTick = false
		tc.isAnimating = true
	}

	if tc.isAnimating {
		tc.spriteSheetTimer += dt
		if tc.spriteSheetTimer > 0.06 {
			tc.spriteSheetTimer = 0
			tc.spriteSheetIdx = (tc.spriteSheetIdx + 1) % tc.spriteFrames
		}

		if tc.spriteSheetIdx == tc.spriteFrames-1 {
			tc.spriteSheetIdx = 0
			tc.isAnimating = false
		}
	}
	subsprite := tc.sprite.SubImage(image.Rect(tc.spriteSheetIdx*16, 0, (tc.spriteSheetIdx+1)*16, 16)).(*ebiten.Image)
	screen.DrawImage(subsprite, &ebiten.DrawImageOptions{GeoM: geom})

	// vector.StrokeCircle(
	// 	screen,
	// 	float32(tc.drawPosition.X+32),
	// 	float32(tc.drawPosition.Y+32),
	// 	float32(tc.radius),
	// 	1.0,
	// 	color.RGBA{255, 0, 255, 255},
	// 	false)
}

// WARN: ShouldFire must be called every tick to determine if the tower should fire
func (tc *Towercore) ShouldFire(dt float64) bool {
	if tc.lastFiredAgo >= tc.rof*(math.Pow(0.9, float64(tc.speedUpgrades))) {
		tc.lastFiredAgo = 0
		return true
	}

	// check and reset temporary buffs
	if tc.tempSpeedBuffTimer.Before(time.Now()) {
		tc.tempSpeedBuff = 0
	}
	if tc.tempDamageBuffTimer.Before(time.Now()) {
		tc.tempDamageBuff = 0
	}

	tc.lastFiredAgo += dt
	return false
}
