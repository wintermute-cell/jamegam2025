package enemy

import (
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
		ret.currentSpeed = 1
	case EnemyTypeFast:
		ret.currentHealth = 1
		ret.currentSpeed = 3
	case EnemyTypeTank:
		ret.currentHealth = 10
		ret.currentSpeed = 0.8
	default:
		panic("Unknown enemy type")
	}

	return ret
}

func (e *Enemy) GetSprite() *ebiten.Image {
	switch e.enemyType {
	case EnemyTypeBasic:
		return SpriteEnemyBasic
	case EnemyTypeFast:
		return SpriteEnemyFast
	case EnemyTypeTank:
		return SpriteEnemyTank
	}

	log.Fatal("Unknown enemy type")
	return nil
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
		audio.Controller.Play("test_deathsound")
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
