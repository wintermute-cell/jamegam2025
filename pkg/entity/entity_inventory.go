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
	hoveredTileIsOnPath bool

	tilePixels int

	// Tower Buttons
	towerSelected         int
	basicTowerNumber      int
	basicTowerButton      lib.Vec2
	basicTowerButtonImage lib.Vec2

	// Resources
	inventorySlotImage *ebiten.Image
	basicTowerImage    *ebiten.Image
}

func isInBounds(vect lib.Vec2I) bool {
	return vect.Y < 12
}

func (e *EntityInventory) isOnPath(vect lib.Vec2I) bool {
	for _, vec := range e.grid.enemyPath {
		if vec.X == vect.X && vec.Y == vect.Y {
			return true
		}
	}
	return false
}

func NewEntityInventory(tilePixels int, grid *EntityGrid) *EntityInventory {
	inventorySlotImage, _, err := ebitenutil.NewImageFromFile("test_inventoryslot.png")
	lib.Must(err)
	basicTowerImage, _, err := ebitenutil.NewImageFromFile("test_tower.png")
	lib.Must(err)
	newEnt := &EntityInventory{
		tilePixels:          tilePixels,
		inventorySlotImage:  inventorySlotImage,
		basicTowerImage:     basicTowerImage,
		inventory:           [4]EntityItemPlaceholder{},
		grid:                grid,
		hoveredTileHasTower: false,
		hoveredTileIsOnPath: false,
		towerSelected:       0,
		basicTowerNumber:    1,
	}
	newEnt.basicTowerButton = newEnt.getTowerButtonPosition(newEnt.basicTowerNumber)
	newEnt.basicTowerButtonImage = newEnt.getTowerButtonIconPosition(newEnt.basicTowerNumber)
	return newEnt
}

func (e *EntityInventory) Init(EntitySpawner) {
}

func (e *EntityInventory) Update(EntitySpawner) error {
	mouseX, mouseY := ebiten.CursorPosition()

	// Tower Buttons
	if isInButton(mouseX, mouseY, e.basicTowerButton) && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		if e.towerSelected == e.basicTowerNumber {
			e.towerSelected = 0
		} else {
			e.towerSelected = e.basicTowerNumber
		}
	}

	// Tower Placement
	e.hoveredTile = lib.NewVec2I(mouseX/e.tilePixels, mouseY/e.tilePixels)
	e.hoveredTileIsOnPath = e.isOnPath(e.hoveredTile)
	_, e.hoveredTileHasTower = e.grid.towers[e.hoveredTile]
	if e.towerSelected != 0 && isInBounds(e.hoveredTile) && !e.hoveredTileIsOnPath && !e.hoveredTileHasTower {
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
	if e.hoveredTileHasTower || e.hoveredTileIsOnPath {
		outlineColor = color.RGBA{255, 100, 100, 255}
	}
	if e.towerSelected != 0 && isInBounds(e.hoveredTile) {
		vector.StrokeRect(screen,
			float32(e.hoveredTile.X*e.tilePixels),
			float32(e.hoveredTile.Y*e.tilePixels),
			float32(e.tilePixels),
			float32(e.tilePixels),
			3.0,
			outlineColor,
			false,
		)
	}

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

	// Items
	for index, _ := range e.inventory {
		geomItem := ebiten.GeoM{}
		geomItem.Scale(6, 6)
		geomItem.Translate(float64(index*e.tilePixels+index*e.tilePixels/2+e.tilePixels/4), float64(12*e.tilePixels+e.tilePixels/4))
		screen.DrawImage(e.inventorySlotImage, &ebiten.DrawImageOptions{GeoM: geomItem})
	}

	// Towers
	geomT1bg := ebiten.GeoM{}
	geomT1bg.Scale(6, 6)
	geomT1bg.Translate(float64(e.basicTowerButton.X), float64(e.basicTowerButton.Y))
	screen.DrawImage(e.inventorySlotImage, &ebiten.DrawImageOptions{GeoM: geomT1bg})
	geomT1im := ebiten.GeoM{}
	geomT1im.Scale(4, 4)
	geomT1im.Translate(float64(e.basicTowerButtonImage.X), float64(e.basicTowerButtonImage.Y))
	screen.DrawImage(e.basicTowerImage, &ebiten.DrawImageOptions{GeoM: geomT1im})

	// Select Tower
	buttonOutline := color.RGBA{100, 255, 100, 255}

	if e.towerSelected == e.basicTowerNumber {
		e.highlightButton(e.basicTowerButton, buttonOutline, screen)
	}
}

func (e *EntityInventory) getTowerButtonPosition(buttonNumber int) lib.Vec2 {
	return lib.NewVec2(float32(16*e.tilePixels-buttonNumber*(7*e.tilePixels/4)), float32(12*e.tilePixels+e.tilePixels/4))
}

func (e *EntityInventory) getTowerButtonIconPosition(buttonNumber int) lib.Vec2 {
	return lib.NewVec2(float32(16*e.tilePixels-(buttonNumber-1)*(7*e.tilePixels/4)-(6*e.tilePixels/4)), float32(12*e.tilePixels+e.tilePixels/2))
}

func isInButton(mouseX int, mouseY int, button lib.Vec2) bool {
	return mouseX >= int(button.X) && mouseX < int(button.X+96) && mouseY >= int(button.Y) && mouseY < int(button.Y+96)
}

func (e *EntityInventory) highlightButton(button lib.Vec2, col color.RGBA, screen *ebiten.Image) {
	vector.StrokeRect(screen,
		button.X,
		button.Y,
		float32(96),
		float32(96),
		3.0,
		col,
		false,
	)
}
