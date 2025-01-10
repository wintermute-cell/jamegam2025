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
		ret.currentHealth = 100
		ret.currentSpeed = 1
	case EnemyTypeFast:
		ret.currentHealth = 100
		ret.currentSpeed = 3
	case EnemyTypeTank:
		ret.currentHealth = 200
		ret.currentSpeed = 0.8
	default:
		panic("Unknown enemy type")
	}

	return ret
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
}

func (e *Enemy) GetSpeed() float32 {
	return e.currentSpeed
}

func (e *Enemy) SetSpeed(speed float32) {
	e.currentSpeed = speed
}