package enemy

import (
	"image"
	"jamegam/pkg/audio"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type EnemyType int

const (
	EnemyTypeBasic EnemyType = iota
	EnemyTypeFast
	EnemyTypeTank
)

type Enemy struct {
	enemyType    EnemyType
	pathNodeLast int
	pathNodeNext int
	pathProgress float64

	destroyFunc func()

	numPassedNodes float64 // The number of path nodes already passed, can be combined with pathProgress to get the exact total path progress

	currentHealth   int
	currentSpeed    float32
	currentSpeedMod float32
	speedModEnd     time.Time

	wander         float32 // the sideways wander from the path line
	WanderVelocity float32
	bounce         float32

	spriteSheetTimer float32
	spriteSheetIndex int
}

func NewEnemy(enemyType EnemyType, pathNodeLast, pathNodeNext int, pathProgress float64) *Enemy {
	ret := &Enemy{
		enemyType:       enemyType,
		pathNodeLast:    pathNodeLast,
		pathNodeNext:    pathNodeNext,
		pathProgress:    pathProgress,
		currentSpeedMod: 1,
	}

	switch enemyType {
	case EnemyTypeBasic:
		ret.currentHealth = 1
		ret.currentSpeed = 1.6
	case EnemyTypeFast:
		ret.currentHealth = 1
		ret.currentSpeed = 4
	case EnemyTypeTank:
		ret.currentHealth = 5
		ret.currentSpeed = 1.1
	default:
		panic("Unknown enemy type")
	}

	return ret
}

func (e *Enemy) GetSprite() *ebiten.Image {
	// i just hope this is called every frame once
	fps := ebiten.ActualFPS()
	dt := 1.0 / 60.0
	if fps > 1/2000 {
		dt = 1.0 / fps
	}
	e.spriteSheetTimer += float32(dt) * (e.currentSpeed * 1.2)
	if e.spriteSheetTimer > 0.1 {
		e.spriteSheetTimer = 0
		mod := 4
		if e.enemyType == EnemyTypeFast {
			mod = 5
		}
		e.spriteSheetIndex = (e.spriteSheetIndex + 1) % mod
	}

	switch e.enemyType {
	case EnemyTypeBasic:
		return SpriteEnemyBasicSheet.SubImage(image.Rect(e.spriteSheetIndex*16, 0, (e.spriteSheetIndex+1)*16, 16)).(*ebiten.Image)
	case EnemyTypeFast:
		// return SpriteEnemyFast
		return SpriteEnemyFastSheet.SubImage(image.Rect(e.spriteSheetIndex*16, 0, (e.spriteSheetIndex+1)*16, 16)).(*ebiten.Image)
	case EnemyTypeTank:
		// return SpriteEnemyTank
		return SpriteEnemyTankSheet.SubImage(image.Rect(e.spriteSheetIndex*16, 0, (e.spriteSheetIndex+1)*16, 16)).(*ebiten.Image)
	}

	log.Fatal("Unknown enemy type")
	return nil
}

func (e *Enemy) GetWander() float32 {
	return e.wander
}

func (e *Enemy) SetWander(wander float32) {
	e.wander = wander
}

func (e *Enemy) GetBounce() float32 {
	return e.bounce
}

func (e *Enemy) SetBounce(bounce float32) {
	e.bounce = bounce
}

func (e *Enemy) SetDestroyFunc(f func()) {
	e.destroyFunc = f
}

func (e *Enemy) GetPathNodes() (last, next int) {
	return e.pathNodeLast, e.pathNodeNext
}

func (e *Enemy) SetPathNodes(last, next int) {
	e.pathNodeLast = last
	e.pathNodeNext = next
}

func (e *Enemy) GetPathProgress() float64 {
	return e.pathProgress
}

func (e *Enemy) SetPathProgress(pathProgress float64) {
	// we need to check time since we can't know when and how frequently
	// SetPathProgress is called. but its called sometimes, which is good...
	if e.speedModEnd.Before(time.Now()) {
		e.currentSpeedMod = 1
		e.speedModEnd = time.Time{}
	}
	e.pathProgress = pathProgress
}

func (e *Enemy) GetHealth() int {
	return e.currentHealth
}

func (e *Enemy) SetHealth(health int) {
	e.currentHealth = health
	if e.currentHealth <= 0 {
		audio.Controller.Play("test_deathsound", 0.2)
		e.destroyFunc()
	}
}

func (e *Enemy) GetSpeed() float32 {
	return e.currentSpeed * e.currentSpeedMod
}

func (e *Enemy) SetSpeedMod(speedMod float32, howLong float32) {
	// This is criminal...
	e.speedModEnd = time.Now().Add(time.Duration(howLong) * time.Second)
	e.currentSpeedMod = speedMod
}

func (e *Enemy) GetSpeedMod() float32 {
	return e.currentSpeedMod
}

func (e *Enemy) GetValue() int64 {
	switch e.enemyType {
	case EnemyTypeBasic:
		return 1
	case EnemyTypeFast:
		return 2
	case EnemyTypeTank:
		return 4
	}
	return 0
}

func (e *Enemy) GetNumPassedNodes() float64 {
	return e.numPassedNodes
}

func (e *Enemy) SetNumPassedNodes(numPassedNodes float64) {
	e.numPassedNodes = numPassedNodes
}
