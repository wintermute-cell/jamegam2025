package entity

import (
	"image/color"
	"jamegam/pkg/lib"
	"jamegam/pkg/towers"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// Ensure EntityInventory implements Entity
var _ Entity = &EntityInventory{}

// TODO: Replace with real EntityItem later
type EntityItemPlaceholder struct {
}

type EntityInventory struct {
	inventory [4]EntityItemPlaceholder
	grid      *EntityGrid

	hoveredTile         lib.Vec2I
	hoveredTileHasTower bool

	tilePixels int

	// Resources
	inventorySlotImage *ebiten.Image
}

func NewEntityInventory(tilePixels int, grid *EntityGrid) *EntityInventory {
	inventorySlotImage, _, err := ebitenutil.NewImageFromFile("test_inventoryslot.png")
	lib.Must(err)
	newEnt := &EntityInventory{
		tilePixels:          tilePixels,
		inventorySlotImage:  inventorySlotImage,
		inventory:           [4]EntityItemPlaceholder{},
		grid:                grid,
		hoveredTileHasTower: false,
	}
	return newEnt
}

func (e *EntityInventory) Init(EntitySpawner) {

}

func (e *EntityInventory) Update(EntitySpawner) error {
	// Tower Placement
	mouseX, mouseY := ebiten.CursorPosition()
	e.hoveredTile = lib.NewVec2I(mouseX/e.tilePixels, mouseY/e.tilePixels)
	_, e.hoveredTileHasTower = e.grid.towers[e.hoveredTile]
	if !e.hoveredTileHasTower {
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			tower := towers.NewTowerBasic(e.hoveredTile.Mul(e.tilePixels))
			e.grid.towers[e.hoveredTile] = tower
		}
	}

	return nil
}

func (e *EntityInventory) Deinit(EntitySpawner) {

}

func (e *EntityInventory) Draw(screen *ebiten.Image) {

	// Tower Placement
	outlineColor := color.RGBA{100, 255, 100, 255}
	if e.hoveredTileHasTower {
		outlineColor = color.RGBA{255, 100, 100, 255}
	}
	vector.StrokeRect(screen,
		float32(e.hoveredTile.X*e.tilePixels),
		float32(e.hoveredTile.Y*e.tilePixels),
		float32(e.tilePixels),
		float32(e.tilePixels),
		3.0,
		outlineColor,
		false,
	)

	// Inventory Bar
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
