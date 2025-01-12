package wavecontroller

import (
	"fmt"
	"jamegam/pkg/enemy"
	"math/rand"
)

type WaveController struct {
	resources int64
	peacetime bool
}

func NewWaveController(starting_resources int64) *WaveController {
	newEnt := &WaveController{
		resources: starting_resources,
	}
	newEnt.Init()
	return newEnt
}

func (e *WaveController) Init() {

}

func (e *WaveController) GetResources() int64 {
	return e.resources
}

func (e *WaveController) GenerateNextWave() []enemy.EnemyType {
	next_enemies := []enemy.EnemyType{}
	var currentCost int64
	for currentCost = 0; currentCost < e.resources; {
		budget := e.resources - currentCost
		random := rand.Intn(100)
		if random < 75 {
			// Add basic enemy
			next_enemies = append(next_enemies, enemy.EnemyTypeBasic)
			currentCost += 1
		} else if random < 90 {
			// Try to add fast enemy
			if budget >= 2 {
				next_enemies = append(next_enemies, enemy.EnemyTypeFast)
				currentCost += 2
			} else {
				next_enemies = append(next_enemies, enemy.EnemyTypeBasic)
				currentCost += 1
			}
		} else {
			// Try to add tank
			if budget >= 4 {
				next_enemies = append(next_enemies, enemy.EnemyTypeTank)
				currentCost += 4
			} else if budget >= 2 {
				next_enemies = append(next_enemies, enemy.EnemyTypeFast)
				currentCost += 2
			} else {
				next_enemies = append(next_enemies, enemy.EnemyTypeBasic)
				currentCost += 1
			}
		}
	}
	fmt.Print(next_enemies)
	return next_enemies
}

func (e *WaveController) IncreaseResources() {
	e.resources += int64(float64(e.resources) * 0.1)
	fmt.Printf("Next Wave Budget: %d", e.resources)
}

func (e *WaveController) Reset() {
	e.resources = 100
}

func (e *WaveController) Deinit() {

}
