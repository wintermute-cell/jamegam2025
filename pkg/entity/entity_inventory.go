package entity

import (
	"image/color"
	"jamegam/pkg/lib"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// Ensure EntityInventory implements Entity
var _ Entity = &EntityInventory{}

// TODO: Replace with real EntityItem later
type EntityItemPlaceholder struct {
}

type EntityInventory struct {
	inventory [4]EntityItemPlaceholder

	tilePixels int

	// Resources
	inventorySlotImage *ebiten.Image
}

func NewEntityInventory(tilePixels int) *EntityInventory {
	inventorySlotImage, _, err := ebitenutil.NewImageFromFile("test_inventoryslot.png")
	lib.Must(err)
	newEnt := &EntityInventory{
		tilePixels:         tilePixels,
		inventorySlotImage: inventorySlotImage,
		inventory:          [4]EntityItemPlaceholder{},
	}
	return newEnt
}

func (e *EntityInventory) Init(EntitySpawner) {

}

func (e *EntityInventory) Update(EntitySpawner) error {
	return nil
}

func (e *EntityInventory) Deinit(EntitySpawner) {

}

func (e *EntityInventory) Draw(screen *ebiten.Image) {
	border := ebiten.NewImage(16*e.tilePixels, 2*e.tilePixels)
	border.Fill(color.RGBA{120, 120, 120, 255})
	geomBord := ebiten.GeoM{}
	geomBord.Translate(0, float64(12*e.tilePixels))
	screen.DrawImage(border, &ebiten.DrawImageOptions{GeoM: geomBord})

	background := ebiten.NewImage(16*e.tilePixels-20, 2*e.tilePixels-20)
	background.Fill(color.RGBA{40, 40, 40, 255})
	geomBack := ebiten.GeoM{}
	geomBack.Translate(10, float64(12*e.tilePixels)+10)
	screen.DrawImage(background, &ebiten.DrawImageOptions{GeoM: geomBack})

	for index, _ := range e.inventory {
		geomItem := ebiten.GeoM{}
		geomItem.Scale(6, 6)
		geomItem.Translate(float64(index*e.tilePixels+index*e.tilePixels/2+e.tilePixels/4), float64(12*e.tilePixels+e.tilePixels/4))
		screen.DrawImage(e.inventorySlotImage, &ebiten.DrawImageOptions{GeoM: geomItem})
	}
}
