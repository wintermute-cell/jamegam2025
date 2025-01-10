package entity

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// Ensure EntityGrid implements Entity
var _ Entity = &EntityGrid{}

type EntityGrid struct {
	xTiles     int
	yTiles     int
	tilePixels int
}

func NewEntityGrid(
	xTiles int,
	yTiles int,
	tilePixels int,
) *EntityGrid {
	newEnt := &EntityGrid{
		xTiles:     xTiles,
		yTiles:     yTiles,
		tilePixels: tilePixels,
	}
	return newEnt
}

func (e *EntityGrid) Init(EntitySpawner) {

}

func (e *EntityGrid) Update(EntitySpawner) error {
	return nil
}

func (e *EntityGrid) Deinit(EntitySpawner) {

}

func (e *EntityGrid) Draw(screen *ebiten.Image) {
	//vector.StrokeLine(screen,
	// for x := range e.xTiles {
	// 	for y := range e.yTiles {
	for x := 0; x <= e.xTiles; x++ {
		for y := 0; y <= e.yTiles; y++ {
			vector.StrokeLine(screen,
				float32(x*e.tilePixels),
				float32(y*e.tilePixels),
				float32(x+1*e.tilePixels),
				float32(y*e.tilePixels),
				3,
				color.RGBA{255, 255, 255, 255},
				false)
			vector.StrokeLine(screen,
				float32(x*e.tilePixels),
				float32(y*e.tilePixels),
				float32(x*e.tilePixels),
				float32(y+1*e.tilePixels),
				3,
				color.RGBA{255, 255, 255, 255},
				false)
		}
	}
}
