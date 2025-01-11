package entity

import (
	"fmt"
	"image/color"
	"jamegam/pkg/enemy"
	"jamegam/pkg/lib"
	"jamegam/pkg/towers"
	"jamegam/pkg/wave_controller"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// Ensure EntityInventory implements Entity
var _ Entity = &EntityInventory{}

// TODO: Replace with real EntityItem later
type EntityItemPlaceholder struct {
}

type EntityInventory struct {
	inventory       [4]EntityItemPlaceholder
	grid            *EntityGrid
	waveController  *wavecontroller.WaveController
	currentWave     []enemy.EnemyType
	peace           bool
	enemySpawnTimer float64

	hoveredTile         lib.Vec2I
	hoveredTileHasTower bool
	hoveredTileIsOnPath bool

	tilePixels int

	currentMana int64
	maximumMana int64

	// Currency
	currentCurrency int64

	// Tower Buttons
	towerSelected         towers.TowerType
	basicTowerButton      lib.Vec2
	basicTowerButtonImage lib.Vec2

	// Resources
	inventorySlotImage *ebiten.Image
	basicTowerImage    *ebiten.Image
	hatImage           *ebiten.Image
	textFace           *text.GoTextFace
}

func isInBounds(vect lib.Vec2I) bool {
	return vect.Y < 12 && vect.X < 16
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
	hatImage, _, err := ebitenutil.NewImageFromFile("test_hat.png")
	lib.Must(err)
	arialFile, err := ebitenutil.OpenFile("Arial.ttf")
	lib.Must(err)
	textFaceSource, err := text.NewGoTextFaceSource(arialFile)
	lib.Must(err)

	newEnt := &EntityInventory{
		tilePixels:          tilePixels,
		inventorySlotImage:  inventorySlotImage,
		basicTowerImage:     basicTowerImage,
		hatImage:            hatImage,
		inventory:           [4]EntityItemPlaceholder{},
		grid:                grid,
		hoveredTileHasTower: false,
		hoveredTileIsOnPath: false,
		towerSelected:       0,
		currentMana:         0,
		maximumMana:         500,
		textFace:            &text.GoTextFace{Source: textFaceSource, Size: 24},
		waveController:      wavecontroller.NewWaveController(100),
		peace:               true,
		enemySpawnTimer:     0.0,
		currentCurrency:     1_000, // TODO: remove this
	}
	newEnt.basicTowerButton = newEnt.getTowerButtonPosition(int(towers.TowerTypeBasic))
	newEnt.basicTowerButtonImage = newEnt.getTowerButtonIconPosition(int(towers.TowerTypeBasic))
	return newEnt
}

func (e *EntityInventory) Init(EntitySpawner) {
}

func (e *EntityInventory) Update(EntitySpawner) error {
	mouseX, mouseY := ebiten.CursorPosition()

	if e.peace {
		e.enemySpawnTimer = 0.0
	} else {
		dt := lib.Dt()
		e.enemySpawnTimer += dt
		if len(e.currentWave) > 0 {
			if e.enemySpawnTimer > 0.5 {
				e.enemySpawnTimer = 0
				e.grid.SpawnEnemy(e.currentWave[0])
				e.currentWave = e.currentWave[1:]
			}
		} else {
			e.peace = true
			e.waveController.IncreaseResources()
		}
	}

	e.currentMana += e.grid.droppedMana
	e.grid.droppedMana = 0

	// Start Wave
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		e.currentWave = append(e.currentWave, e.waveController.GenerateNextWave()...)
		e.peace = false
	}

	// Hat Button
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) && e.isInHatButton(mouseX, mouseY) {
		manaPercentage := int(float32(e.currentMana) / float32(e.maximumMana) * 100)
		if manaPercentage < 50 {
			e.currentCurrency += int64(float64(e.currentMana) * 5.0 * 1.0)
		} else if manaPercentage < 75 {
			e.currentCurrency += int64(float64(e.currentMana) * 5.0 * 1.5)
		} else {
			e.currentCurrency += int64(float64(e.currentMana) * 5.0 * 2.0)
		}
		e.currentMana = 0
	}

	// Tower Buttons
	if isInButton(mouseX, mouseY, e.basicTowerButton) && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		if e.towerSelected == towers.TowerTypeBasic {
			e.towerSelected = towers.TowerTypeNone
		} else {
			e.towerSelected = towers.TowerTypeBasic
		}
	}

	// Tower Placement
	e.hoveredTile = lib.NewVec2I(mouseX/e.tilePixels, mouseY/e.tilePixels)
	e.hoveredTileIsOnPath = e.isOnPath(e.hoveredTile)
	_, e.hoveredTileHasTower = e.grid.towers[e.hoveredTile]
	if e.towerSelected != towers.TowerTypeNone && isInBounds(e.hoveredTile) && !e.hoveredTileIsOnPath && !e.hoveredTileHasTower {
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			var tower towers.Tower = nil
			switch e.towerSelected {
			case towers.TowerTypeBasic:
				tower = towers.NewTowerBasic(e.hoveredTile.Mul(e.tilePixels))
				//case towers.TowerType...:
			}
			if tower != nil {
				if e.currentCurrency >= tower.Price() {
					e.currentCurrency -= tower.Price()
					e.grid.towers[e.hoveredTile] = tower
				} else {
					e.grid.ShowMessage(fmt.Sprintf("Not enough currency to place tower. Need %d", tower.Price()))
				}
			}
		}
	} else if e.hoveredTileHasTower && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) && isInBounds(e.hoveredTile) {
		e.grid.ShowMessage("Can't place tower here, there is already a tower.")
	} else if e.hoveredTileIsOnPath && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) && isInBounds(e.hoveredTile) {
		e.grid.ShowMessage("Can't place tower on the path.")
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
	if e.towerSelected != towers.TowerTypeNone && isInBounds(e.hoveredTile) {
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

	// Hat
	geomHat := ebiten.GeoM{}
	geomHat.Scale(4, 4)
	geomHat.Translate(float64(7*e.tilePixels+e.tilePixels/2), float64(12*e.tilePixels+e.tilePixels/4))
	screen.DrawImage(e.hatImage, &ebiten.DrawImageOptions{GeoM: geomHat})

	// Hat Percentage
	manaPercentage := int(float32(e.currentMana) / float32(e.maximumMana) * 100.0)
	hatTextOptions := &text.DrawOptions{}
	hatTextOptions.GeoM.Translate(float64(7*e.tilePixels+5*e.tilePixels/8), float64(13*e.tilePixels+e.tilePixels/8))

	hatTextOptions.ColorScale.Scale(0.5, 1.0, 0.5, 1.0)
	if manaPercentage >= 50 && manaPercentage < 75 {
		hatTextOptions.ColorScale.Reset()
		hatTextOptions.ColorScale.Scale(1.0, 1.0, 0.0, 1.0)
	} else if manaPercentage >= 75 {
		hatTextOptions.ColorScale.Reset()
		hatTextOptions.ColorScale.Scale(1.0, 0.0, 0.0, 1.0)
	}

	text.Draw(screen, fmt.Sprintf("%03d%%", manaPercentage), e.textFace, hatTextOptions)

	geom := ebiten.GeoM{}
	geom.Translate(10, 10)
	text.Draw(screen, fmt.Sprintf("Currency: %d", e.currentCurrency), e.textFace, &text.DrawOptions{
		DrawImageOptions: ebiten.DrawImageOptions{GeoM: geom},
	})

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

	if e.towerSelected == towers.TowerTypeBasic {
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

func (e *EntityInventory) isInHatButton(mouseX int, mouseY int) bool {
	hatButtonX := int(7*e.tilePixels + e.tilePixels/2)
	hatButtonY := int(12*e.tilePixels + e.tilePixels/4)
	return mouseX >= hatButtonX && mouseX < hatButtonX+e.tilePixels && mouseY >= hatButtonY && mouseY < hatButtonY+5*e.tilePixels/4
}
