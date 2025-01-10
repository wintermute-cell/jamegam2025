package game

import (
	"jamegam/pkg/enemy"
	"jamegam/pkg/entity"
)

type WaveController struct {
	resources    int64
	grid         *entity.EntityGrid
	next_enemies []enemy.EnemyType
}

func NewWaveController(starting_resources int64, grid *entity.EntityGrid) *WaveController {
	newEnt := &WaveController{
		resources: starting_resources,
		grid:      grid,
	}
	newEnt.Init()
	return newEnt
}

func (e *WaveController) Init() {

}

func (e *WaveController) generateNextWave() {

}

func (e *WaveController) spawnNextWave() {

}

func (e *WaveController) increaseResources(value int64) {
	e.resources += value
}

func (e *WaveController) Deinit() {

}
