package enemy

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

	currentHealth int
	currentSpeed  float32
}

func NewEnemy(enemyType EnemyType, pathNodeLast, pathNodeNext int, pathProgress float64) *Enemy {
	ret := &Enemy{
		enemyType:    enemyType,
		pathNodeLast: pathNodeLast,
		pathNodeNext: pathNodeNext,
		pathProgress: pathProgress,
	}

	switch enemyType {
	case EnemyTypeBasic:
		ret.currentHealth = 1
		ret.currentSpeed = 1
	case EnemyTypeFast:
		ret.currentHealth = 1
		ret.currentSpeed = 3
	case EnemyTypeTank:
		ret.currentHealth = 4
		ret.currentSpeed = 0.8
	default:
		panic("Unknown enemy type")
	}

	return ret
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
	e.pathProgress = pathProgress
}

func (e *Enemy) GetHealth() int {
	return e.currentHealth
}

func (e *Enemy) SetHealth(health int) {
	e.currentHealth = health
	if e.currentHealth <= 0 {
		e.destroyFunc()
	}
}

func (e *Enemy) GetSpeed() float32 {
	return e.currentSpeed
}

func (e *Enemy) SetSpeed(speed float32) {
	e.currentSpeed = speed
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
