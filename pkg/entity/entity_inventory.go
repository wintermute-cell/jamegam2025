package entity

import (
	"fmt"
	"image/color"
	"jamegam/pkg/enemy"
	"jamegam/pkg/lib"
	"jamegam/pkg/towers"
	"jamegam/pkg/wave_controller"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// Ensure EntityInventory implements Entity
var _ Entity = &EntityInventory{}

type Item int64

const (
	NoItem Item = iota
	BasicTower
	TackTower
	IceTower
	AoeTower
	ManaTower
	FreeUpgrade
	MaxUpgrade
	CurrencyGift
	SpikeTrap
	ClearEnemies
	DamageBuff
	SpeedBuff
)

type EntityInventory struct {
	inventory           [4]Item
	selectedItem        int
	grid                *EntityGrid
	waveController      *wavecontroller.WaveController
	currentWave         []enemy.EnemyType
	peace               bool
	enemySpawnTimer     float64
	waveCounter         int64
	freeTurretSelected  towers.TowerType
	freeUpgradeSelected bool
	maxUpgradeSelected  bool

	hoveredTile          lib.Vec2I
	hoveredTileHasTower  bool
	hoveredTileIsOnPath  bool
	turretRangeIndicator bool

	tilePixels   int
	buttonPixels int

	currentMana int64
	maximumMana int64

	// Currency
	currentCurrency int64

	// Tower Buttons
	basicTowerButton  lib.Vec2I
	tackTowerButton   lib.Vec2I
	iceTowerButton    lib.Vec2I
	aoeTowerButton    lib.Vec2I
	cashTowerButton   lib.Vec2I
	blueprintSelected towers.TowerType

	// Menu Buttons
	playButton     lib.Vec2I
	removeButton   lib.Vec2I
	damageButton   lib.Vec2I
	firerateButton lib.Vec2I

	// Resources
	inventorySlotImage    *ebiten.Image
	basicTowerImage       *ebiten.Image
	tackTowerImage        *ebiten.Image
	iceTowerImage         *ebiten.Image
	aoeTowerImage         *ebiten.Image
	cashTowerImage        *ebiten.Image
	hatImage              *ebiten.Image
	textFace              *text.GoTextFace
	playButtonImage       *ebiten.Image
	removeButtonImage     *ebiten.Image
	damageButtonImage     *ebiten.Image
	firerateButtonImage   *ebiten.Image
	inventoryBarImage     *ebiten.Image
	upgradeIndicatorImage *ebiten.Image
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
	inventorySlotImage, _, err := ebitenutil.NewImageFromFile("inventory_slot_24x24.png")
	lib.Must(err)
	basicTowerImage, _, err := ebitenutil.NewImageFromFile("test_tower.png")
	lib.Must(err)
	tackTowerImage, _, err := ebitenutil.NewImageFromFile("test_towertacks.png")
	lib.Must(err)
	iceTowerImage, _, err := ebitenutil.NewImageFromFile("test_towerice.png")
	lib.Must(err)
	aoeTowerImage, _, err := ebitenutil.NewImageFromFile("test_toweraoe.png")
	lib.Must(err)
	hatImage, _, err := ebitenutil.NewImageFromFile("test_hat.png")
	lib.Must(err)
	arialFile, err := ebitenutil.OpenFile("font.ttf")
	lib.Must(err)
	textFaceSource, err := text.NewGoTextFaceSource(arialFile)
	lib.Must(err)
	inventoryBarImage, _, err := ebitenutil.NewImageFromFile("menu_bar_1024x246.png")
	lib.Must(err)
	playButtonImage, _, err := ebitenutil.NewImageFromFile("test_playbutton.png")
	lib.Must(err)
	removeButtonImage, _, err := ebitenutil.NewImageFromFile("test_removebutton.png")
	lib.Must(err)
	damageButtonImage, _, err := ebitenutil.NewImageFromFile("test_damagebutton.png")
	lib.Must(err)
	firerateButtonImage, _, err := ebitenutil.NewImageFromFile("test_ammobutton.png")
	lib.Must(err)
	upgradeIndicatorImage, _, err := ebitenutil.NewImageFromFile("upgradeindicator.png")
	lib.Must(err)
	cashTowerImage, _, err := ebitenutil.NewImageFromFile("test_towercash.png")
	lib.Must(err)

	newEnt := &EntityInventory{
		tilePixels:            tilePixels,
		buttonPixels:          96,
		inventorySlotImage:    inventorySlotImage,
		basicTowerImage:       basicTowerImage,
		tackTowerImage:        tackTowerImage,
		iceTowerImage:         iceTowerImage,
		aoeTowerImage:         aoeTowerImage,
		cashTowerImage:        cashTowerImage,
		hatImage:              hatImage,
		inventory:             [4]Item{BasicTower, IceTower, FreeUpgrade, CurrencyGift},
		selectedItem:          -1,
		grid:                  grid,
		hoveredTileHasTower:   false,
		hoveredTileIsOnPath:   false,
		blueprintSelected:     0,
		currentMana:           0,
		maximumMana:           500,
		textFace:              &text.GoTextFace{Source: textFaceSource, Size: 20},
		waveController:        wavecontroller.NewWaveController(100),
		peace:                 true,
		enemySpawnTimer:       0.0,
		currentCurrency:       1_000, // TODO: remove this
		waveCounter:           0,
		turretRangeIndicator:  true,
		freeTurretSelected:    towers.TowerTypeNone,
		freeUpgradeSelected:   false,
		maxUpgradeSelected:    false,
		inventoryBarImage:     inventoryBarImage,
		playButtonImage:       playButtonImage,
		removeButtonImage:     removeButtonImage,
		damageButtonImage:     damageButtonImage,
		firerateButtonImage:   firerateButtonImage,
		upgradeIndicatorImage: upgradeIndicatorImage,
		basicTowerButton:      lib.NewVec2I(2, 1),
		tackTowerButton:       lib.NewVec2I(3, 1),
		iceTowerButton:        lib.NewVec2I(4, 1),
		aoeTowerButton:        lib.NewVec2I(5, 1),
		cashTowerButton:       lib.NewVec2I(6, 1),
		playButton:            lib.NewVec2I(0, 0),
		removeButton:          lib.NewVec2I(1, 0),
		damageButton:          lib.NewVec2I(2, 0),
		firerateButton:        lib.NewVec2I(3, 0),
	}
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
			if e.enemySpawnTimer > 0.8 {
				e.enemySpawnTimer = (rand.Float64() - 0.5) * 0.7
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

	// Start Wave Button and Hotkey
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) && isInButton(mouseX, mouseY, e.getButtonPosition(e.playButton)) {
		e.StartWave()
	} else if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		e.StartWave()
	}

	// Toggle Turret Range Indicators
	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		e.ToggleTowerIndicator()
	}

	// Hat Button
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) && e.isInHatButton(mouseX, mouseY) {
		manaPercentage := int(float32(e.currentMana) / float32(e.maximumMana) * 100)
		var newCurrency int64 = 0
		if manaPercentage < 50 {
			newCurrency += int64(float64(e.currentMana) * 5.0 * 1.0)
		} else if manaPercentage < 75 {
			newCurrency += int64(float64(e.currentMana) * 5.0 * 1.5)
		} else {
			newCurrency += int64(float64(e.currentMana) * 5.0 * 2.0)
		}
		e.currentCurrency += newCurrency
		e.currentMana = 0
		e.grid.ShowMessage(fmt.Sprintf("Received %d currency!", newCurrency))
	}

	// Item Buttons
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		if isInButton(mouseX, mouseY, e.getButtonPosition(lib.NewVec2I(5, 0))) {
			e.ActivateItem(0)
		} else if isInButton(mouseX, mouseY, e.getButtonPosition(lib.NewVec2I(6, 0))) {
			e.ActivateItem(1)
		} else if isInButton(mouseX, mouseY, e.getButtonPosition(lib.NewVec2I(7, 0))) {
			e.ActivateItem(2)
		} else if isInButton(mouseX, mouseY, e.getButtonPosition(lib.NewVec2I(8, 0))) {
			e.ActivateItem(3)
		}
	}

	// Tower Buttons
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		if isInButton(mouseX, mouseY, e.getButtonPosition(e.basicTowerButton)) {
			e.selectTowerType(towers.TowerTypeBasic)
		} else if isInButton(mouseX, mouseY, e.getButtonPosition(e.tackTowerButton)) {
			e.selectTowerType(towers.TowerTypeTacks)
		} else if isInButton(mouseX, mouseY, e.getButtonPosition(e.iceTowerButton)) {
			e.selectTowerType(towers.TowerTypeIce)
		} else if isInButton(mouseX, mouseY, e.getButtonPosition(e.aoeTowerButton)) {
			e.selectTowerType(towers.TowerTypeAoe)
		} else if isInButton(mouseX, mouseY, e.getButtonPosition(e.cashTowerButton)) {
			e.selectTowerType(towers.TowerTypeCash)
		}
	}

	// Tower Hotkeys
	if inpututil.IsKeyJustPressed(ebiten.Key1) {
		e.selectTowerType(towers.TowerTypeBasic)
	} else if inpututil.IsKeyJustPressed(ebiten.Key2) {
		e.selectTowerType(towers.TowerTypeTacks)
	} else if inpututil.IsKeyJustPressed(ebiten.Key3) {
		e.selectTowerType(towers.TowerTypeIce)
	} else if inpututil.IsKeyJustPressed(ebiten.Key4) {
		e.selectTowerType(towers.TowerTypeAoe)
	} else if inpututil.IsKeyJustPressed(ebiten.Key5) {
		e.selectTowerType(towers.TowerTypeCash)
	}

	// Tower Placement
	e.hoveredTile = lib.NewVec2I(mouseX/e.tilePixels, mouseY/e.tilePixels)
	e.hoveredTileIsOnPath = e.isOnPath(e.hoveredTile)
	_, e.hoveredTileHasTower = e.grid.towers[e.hoveredTile]
	if (e.blueprintSelected != towers.TowerTypeNone || e.freeTurretSelected != towers.TowerTypeNone) && isInBounds(e.hoveredTile) && !e.hoveredTileIsOnPath && !e.hoveredTileHasTower {
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			selectedTowerType := towers.TowerTypeBasic
			free := false
			if e.blueprintSelected != towers.TowerTypeNone {
				selectedTowerType = e.blueprintSelected
			} else if e.freeTurretSelected != towers.TowerTypeNone {
				selectedTowerType = e.freeTurretSelected
				free = true
			}
			var tower towers.Tower = nil
			switch selectedTowerType {
			case towers.TowerTypeBasic:
				tower = towers.NewTowerBasic(e.hoveredTile.Mul(e.tilePixels))
			case towers.TowerTypeTacks:
				tower = towers.NewTowerTacks(e.hoveredTile.Mul(e.tilePixels))
			case towers.TowerTypeIce:
				tower = towers.NewTowerIce(e.hoveredTile.Mul(e.tilePixels))
			case towers.TowerTypeAoe:
				tower = towers.NewTowerAoe(e.hoveredTile.Mul(e.tilePixels))
			case towers.TowerTypeCash:
				tower = towers.NewTowerCash(e.hoveredTile.Mul(e.tilePixels))
				// case towers.TowerType...:
			}
			if tower != nil {
				if free || e.currentCurrency >= tower.Price() {
					if free {
						e.RemoveItem(e.selectedItem)
						e.ClearSelectedItem()
					} else {
						e.currentCurrency -= tower.Price()
					}
					e.grid.towers[e.hoveredTile] = tower
					e.grid.selectedTower = e.hoveredTile
				} else {
					e.grid.ShowMessage(fmt.Sprintf("Not enough currency to place tower. Need %d", tower.Price()))
				}
			}
		}
	} else if (e.blueprintSelected != towers.TowerTypeNone || e.freeTurretSelected != towers.TowerTypeNone) && e.hoveredTileIsOnPath && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) && isInBounds(e.hoveredTile) {
		e.grid.ShowMessage("Can't place tower on the path.")
	} else if isInBounds(e.hoveredTile) && e.hoveredTileHasTower {
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			e.blueprintSelected = towers.TowerTypeNone
			if e.freeTurretSelected != towers.TowerTypeNone {
				e.freeTurretSelected = towers.TowerTypeNone
				e.selectedItem = -1
			}
			e.grid.selectedTower = e.hoveredTile
		}
	}

	// Cancel tower selection and blueprint placement
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		e.blueprintSelected = towers.TowerTypeNone
		e.grid.selectedTower = lib.NewVec2I(-1, -1)
	}

	// Unselect Tower
	if e.blueprintSelected == towers.TowerTypeNone && isInBounds(e.hoveredTile) && !e.hoveredTileHasTower {
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			e.grid.selectedTower = lib.NewVec2I(-1, -1)
		}
	}

	// Upgrade Tower Damage Button and Hotkey
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) && isInButton(mouseX, mouseY, e.getButtonPosition(e.damageButton)) {
		e.UpgradeSelectedTowerDamage()
	} else if inpututil.IsKeyJustPressed(ebiten.KeyD) {
		e.UpgradeSelectedTowerDamage()
	}

	// Upgrade Tower Speed Button and Hotkey
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) && isInButton(mouseX, mouseY, e.getButtonPosition(e.firerateButton)) {
		e.UpgradeSelectedTowerSpeed()
	} else if inpututil.IsKeyJustPressed(ebiten.KeyS) {
		e.UpgradeSelectedTowerSpeed()
	}

	// Sell Tower Button and Hotkey
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) && isInButton(mouseX, mouseY, e.getButtonPosition(e.removeButton)) {
		e.SellSelectedTower()
	} else if inpututil.IsKeyJustPressed(ebiten.KeyX) {
		e.SellSelectedTower()
	}

	return nil
}

func (e *EntityInventory) Deinit(EntitySpawner) {

}

func (e *EntityInventory) Draw(screen *ebiten.Image) {
	buttonOutline := color.RGBA{100, 255, 100, 255}

	// Tower Placement
	outlineColor := color.RGBA{100, 255, 100, 255}
	if e.hoveredTileHasTower || e.hoveredTileIsOnPath {
		outlineColor = color.RGBA{255, 100, 100, 255}
	}
	if (e.blueprintSelected != towers.TowerTypeNone || e.freeTurretSelected != towers.TowerTypeNone) && isInBounds(e.hoveredTile) {
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
	geomBord := ebiten.GeoM{}
	geomBord.Translate(0, float64(12*e.tilePixels))
	screen.DrawImage(e.inventoryBarImage, &ebiten.DrawImageOptions{GeoM: geomBord})

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

	// Item Slots
	for index, _ := range e.inventory {
		itemPosition := e.getButtonPosition(lib.NewVec2I(index+5, 0))
		geomItem := ebiten.GeoM{}
		geomItem.Scale(4, 4)
		geomItem.Translate(float64(itemPosition.X), float64(itemPosition.Y))
		screen.DrawImage(e.inventorySlotImage, &ebiten.DrawImageOptions{GeoM: geomItem})
	}

	// Item Icons
	for index, itemType := range e.inventory {
		if itemType != NoItem {
			itemPos := e.getButtonTowerIconPosition(lib.NewVec2I(index+5, 0))
			geomIcon := ebiten.GeoM{}
			geomIcon.Scale(4, 4)
			geomIcon.Translate(float64(itemPos.X), float64(itemPos.Y))
			screen.DrawImage(e.GetItemIcon(itemType), &ebiten.DrawImageOptions{GeoM: geomIcon})
		}
	}

	// Item Slot Highlighting
	if e.selectedItem == 0 {
		e.highlightButton(e.getButtonPosition(lib.NewVec2I(5, 0)), buttonOutline, screen)
	} else if e.selectedItem == 1 {
		e.highlightButton(e.getButtonPosition(lib.NewVec2I(6, 0)), buttonOutline, screen)
	} else if e.selectedItem == 2 {
		e.highlightButton(e.getButtonPosition(lib.NewVec2I(7, 0)), buttonOutline, screen)
	} else if e.selectedItem == 3 {
		e.highlightButton(e.getButtonPosition(lib.NewVec2I(8, 0)), buttonOutline, screen)
	}

	// Tower Buttons
	for i := 2; i < 7; i++ {
		towerButtonPosition := e.getButtonPosition(lib.NewVec2I(i, 1))
		geomT1bg := ebiten.GeoM{}
		geomT1bg.Scale(4, 4)
		geomT1bg.Translate(float64(towerButtonPosition.X), float64(towerButtonPosition.Y))
		screen.DrawImage(e.inventorySlotImage, &ebiten.DrawImageOptions{GeoM: geomT1bg})
	}

	// Tower Button Icons
	// Basic Tower
	basicTowerImgPos := e.getButtonTowerIconPosition(e.basicTowerButton)
	geomT1im := ebiten.GeoM{}
	geomT1im.Scale(4, 4)
	geomT1im.Translate(float64(basicTowerImgPos.X), float64(basicTowerImgPos.Y))
	screen.DrawImage(e.basicTowerImage, &ebiten.DrawImageOptions{GeoM: geomT1im})
	// Tack Tower
	tackTowerImgPos := e.getButtonTowerIconPosition(e.tackTowerButton)
	geomT2im := ebiten.GeoM{}
	geomT2im.Scale(4, 4)
	geomT2im.Translate(float64(tackTowerImgPos.X), float64(tackTowerImgPos.Y))
	screen.DrawImage(e.tackTowerImage, &ebiten.DrawImageOptions{GeoM: geomT2im})
	// Ice Tower
	iceTowerImgPos := e.getButtonTowerIconPosition(e.iceTowerButton)
	geomT3im := ebiten.GeoM{}
	geomT3im.Scale(4, 4)
	geomT3im.Translate(float64(iceTowerImgPos.X), float64(iceTowerImgPos.Y))
	screen.DrawImage(e.iceTowerImage, &ebiten.DrawImageOptions{GeoM: geomT3im})
	// AOE Tower
	aoeTowerImgPos := e.getButtonTowerIconPosition(e.aoeTowerButton)
	geomT4im := ebiten.GeoM{}
	geomT4im.Scale(4, 4)
	geomT4im.Translate(float64(aoeTowerImgPos.X), float64(aoeTowerImgPos.Y))
	screen.DrawImage(e.aoeTowerImage, &ebiten.DrawImageOptions{GeoM: geomT4im})
	// Cash Tower
	cashTowerImgPos := e.getButtonTowerIconPosition(e.cashTowerButton)
	geomT5im := ebiten.GeoM{}
	geomT5im.Scale(4, 4)
	geomT5im.Translate(float64(cashTowerImgPos.X), float64(cashTowerImgPos.Y))
	screen.DrawImage(e.cashTowerImage, &ebiten.DrawImageOptions{GeoM: geomT5im})

	// Select Tower
	if e.blueprintSelected == towers.TowerTypeBasic {
		e.highlightButton(e.getButtonPosition(e.basicTowerButton), buttonOutline, screen)
	} else if e.blueprintSelected == towers.TowerTypeTacks {
		e.highlightButton(e.getButtonPosition(e.tackTowerButton), buttonOutline, screen)
	} else if e.blueprintSelected == towers.TowerTypeIce {
		e.highlightButton(e.getButtonPosition(e.iceTowerButton), buttonOutline, screen)
	} else if e.blueprintSelected == towers.TowerTypeAoe {
		e.highlightButton(e.getButtonPosition(e.aoeTowerButton), buttonOutline, screen)
	} else if e.blueprintSelected == towers.TowerTypeCash {
		e.highlightButton(e.getButtonPosition(e.cashTowerButton), buttonOutline, screen)
	}

	// Menu Buttons
	for i := 0; i < 4; i++ {
		buttonPos := e.getButtonPosition(lib.NewVec2I(i, 0))
		buttonImgOptions := &ebiten.DrawImageOptions{}
		buttonImgOptions.GeoM.Scale(4, 4)
		buttonImgOptions.GeoM.Translate(float64(buttonPos.X), float64(buttonPos.Y))
		screen.DrawImage(e.inventorySlotImage, buttonImgOptions)
	}

	// Menu Button Icons
	// Play Button
	playButtonImagePos := e.getButtonTowerIconPosition(e.playButton)
	geomUI1 := ebiten.GeoM{}
	geomUI1.Scale(4, 4)
	geomUI1.Translate(float64(playButtonImagePos.X), float64(playButtonImagePos.Y))
	screen.DrawImage(e.playButtonImage, &ebiten.DrawImageOptions{GeoM: geomUI1})
	// Remove Button
	removeButtonImagePos := e.getButtonTowerIconPosition(e.removeButton)
	geomUI2 := ebiten.GeoM{}
	geomUI2.Scale(4, 4)
	geomUI2.Translate(float64(removeButtonImagePos.X), float64(removeButtonImagePos.Y))
	screen.DrawImage(e.removeButtonImage, &ebiten.DrawImageOptions{GeoM: geomUI2})
	// Damage Upgrade Button
	damageButtonImagePos := e.getButtonTowerIconPosition(e.damageButton)
	geomUI3 := ebiten.GeoM{}
	geomUI3.Scale(4, 4)
	geomUI3.Translate(float64(damageButtonImagePos.X), float64(damageButtonImagePos.Y))
	screen.DrawImage(e.damageButtonImage, &ebiten.DrawImageOptions{GeoM: geomUI3})
	// Fire Rate Upgrade Button
	firerateButtonImagePos := e.getButtonTowerIconPosition(e.firerateButton)
	geomUI4 := ebiten.GeoM{}
	geomUI4.Scale(4, 4)
	geomUI4.Translate(float64(firerateButtonImagePos.X), float64(firerateButtonImagePos.Y))
	screen.DrawImage(e.firerateButtonImage, &ebiten.DrawImageOptions{GeoM: geomUI4})

	// Upgrade Indicators
	if e.isTowerSelected() {
		tow := e.grid.towers[e.grid.selectedTower]
		// Damage
		dmgButtonPos := e.getButtonPosition(e.damageButton)
		for i := 0; i < int(tow.GetDamageUpgrades()); i++ {
			geomDI := ebiten.GeoM{}
			geomDI.Scale(4, 4)
			geomDI.Translate(float64(dmgButtonPos.X)+(16.0+4.0)*float64(i), float64(dmgButtonPos.Y))
			screen.DrawImage(e.upgradeIndicatorImage, &ebiten.DrawImageOptions{GeoM: geomDI})
		}

		// Speed
		spdButtonPos := e.getButtonPosition(e.firerateButton)
		for i := 0; i < int(tow.GetSpeedUpgrades()); i++ {
			geomSI := ebiten.GeoM{}
			geomSI.Scale(4, 4)
			geomSI.Translate(float64(spdButtonPos.X)+(16.0+4.0)*float64(i), float64(spdButtonPos.Y))
			screen.DrawImage(e.upgradeIndicatorImage, &ebiten.DrawImageOptions{GeoM: geomSI})
		}

	}
}

func isInButton(mouseX int, mouseY int, button lib.Vec2I) bool {
	return mouseX >= button.X && mouseX < button.X+96 && mouseY >= button.Y && mouseY < button.Y+96
}

func (e *EntityInventory) highlightButton(button lib.Vec2I, col color.RGBA, screen *ebiten.Image) {
	vector.StrokeRect(screen,
		float32(button.X),
		float32(button.Y),
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

func (e *EntityInventory) ToggleTowerIndicator() {
	if e.turretRangeIndicator {
		e.grid.ShowMessage("Range indicator disabled.")
	} else {
		e.grid.ShowMessage("Range indicator enabled.")
	}
	e.turretRangeIndicator = !e.turretRangeIndicator
	e.grid.towerRangeIndicator = e.turretRangeIndicator
}

func (e *EntityInventory) getButtonPosition(position lib.Vec2I) lib.Vec2I {
	buttonPos := lib.NewVec2I(0, 0)
	buttonPos.X = int(24 + position.X*(e.buttonPixels+14))
	buttonPos.Y = int(12*e.tilePixels + 18 + position.Y*(e.buttonPixels+18))
	return buttonPos
}

func (e *EntityInventory) getButtonTowerIconPosition(position lib.Vec2I) lib.Vec2I {
	buttonPos := e.getButtonPosition(position)
	iconPos := lib.NewVec2I(buttonPos.X+16, buttonPos.Y+16)
	return iconPos
}

func (e *EntityInventory) selectTowerType(towerType towers.TowerType) {
	if e.blueprintSelected == towerType {
		e.blueprintSelected = towers.TowerTypeNone
	} else {
		e.blueprintSelected = towerType
		if e.freeTurretSelected != towers.TowerTypeNone {
			e.freeTurretSelected = towers.TowerTypeNone
			e.selectedItem = -1
		}
	}
}

func (e *EntityInventory) StartWave() {
	e.currentWave = append(e.currentWave, e.waveController.GenerateNextWave()...)
	e.peace = false
	e.waveCounter++
	e.grid.ShowMessage(fmt.Sprintf("Wave %d started! (Strength: %d)", e.waveCounter, e.waveController.GetResources()))
}

func (e *EntityInventory) SellSelectedTower() {
	if e.isTowerSelected() {
		sellPrice := int64(float64(e.grid.towers[e.grid.selectedTower].Price())*0.5) + int64(e.grid.towers[e.grid.selectedTower].GetTotalUpgrades()*100)
		delete(e.grid.towers, e.grid.selectedTower)
		e.grid.selectedTower = lib.NewVec2I(-1, -1)
		e.currentCurrency += sellPrice
		e.grid.ShowMessage(fmt.Sprintf("Sold selected tower for %d!", sellPrice))
	}
}

func (e *EntityInventory) UpgradeSelectedTowerSpeed() {
	if !e.isTowerSelected() {
		return
	}
	tower := e.grid.towers[e.grid.selectedTower]
	if tower.GetTotalUpgrades() >= 7 {
		e.grid.ShowMessage("This tower has reached the maximum amount of total upgrades!")
	} else if tower.GetSpeedUpgrades() >= 5 {
		e.grid.ShowMessage("This tower has reached the maximum amount of speed upgrades!")
	} else {
		upgradePrice := int64(float64(tower.GetTotalUpgrades()+1) * float64(tower.Price()))
		if e.currentCurrency >= upgradePrice {
			e.currentCurrency -= upgradePrice
			tower.SpeedUpgrade()
			e.grid.ShowMessage("Tower speed upgraded!")
		} else {
			e.grid.ShowMessage(fmt.Sprintf("Not enough currency! (Required: %d)", upgradePrice))
		}
	}
}

func (e *EntityInventory) UpgradeSelectedTowerDamage() {
	if !e.isTowerSelected() {
		return
	}
	tower := e.grid.towers[e.grid.selectedTower]
	if tower.GetTotalUpgrades() >= 7 {
		e.grid.ShowMessage("This tower has reached the maximum amount of total upgrades!")
	} else if tower.GetDamageUpgrades() >= 5 {
		e.grid.ShowMessage("This tower has reached the maximum amount of damage upgrades!")
	} else {
		upgradePrice := int64(float64(tower.GetTotalUpgrades()+1) * float64(tower.Price()))
		if e.currentCurrency >= upgradePrice {
			e.currentCurrency -= upgradePrice
			tower.DamageUpgrade()
			e.grid.ShowMessage("Tower damage upgraded!")
		} else {
			e.grid.ShowMessage(fmt.Sprintf("Not enough currency! (Required: %d)", upgradePrice))
		}
	}

}

func (e *EntityInventory) isTowerSelected() bool {
	return e.grid.selectedTower.X != -1 && e.grid.selectedTower.Y != -1
}

func (e *EntityInventory) ClearSelectedItem() {
	e.freeUpgradeSelected = false
	e.maxUpgradeSelected = false
	e.freeTurretSelected = towers.TowerTypeNone
	e.selectedItem = -1
}

func (e *EntityInventory) SelectFreeTurret(towerType towers.TowerType, itemNumber int) {
	e.blueprintSelected = towers.TowerTypeNone
	e.freeTurretSelected = towerType
	e.selectedItem = itemNumber
}

func (e *EntityInventory) ActivateItem(itemNumber int) {
	prevItem := e.selectedItem
	e.ClearSelectedItem()
	if prevItem == itemNumber {
		return
	}
	switch e.inventory[itemNumber] {
	case NoItem:
		e.grid.ShowMessage("This slot is empty!")
	case BasicTower:
		e.SelectFreeTurret(towers.TowerTypeBasic, itemNumber)
	case TackTower:
		e.SelectFreeTurret(towers.TowerTypeTacks, itemNumber)
	case IceTower:
		e.SelectFreeTurret(towers.TowerTypeIce, itemNumber)
	case AoeTower:
		e.SelectFreeTurret(towers.TowerTypeAoe, itemNumber)
	case ManaTower:
		e.SelectFreeTurret(towers.TowerTypeCash, itemNumber)
	case FreeUpgrade:
		e.selectedItem = itemNumber
		e.freeUpgradeSelected = true
	case MaxUpgrade:
		e.selectedItem = itemNumber
		e.maxUpgradeSelected = true
	case CurrencyGift:
		newCurrency := 50 * e.waveCounter
		e.currentCurrency += newCurrency
		e.RemoveItem(itemNumber)
		e.grid.ShowMessage(fmt.Sprintf("Received %d currency!", newCurrency))
	case SpikeTrap:
		return
	case ClearEnemies:
		return
	case DamageBuff:
		return
	case SpeedBuff:
		return
	}
}

func (e *EntityInventory) RemoveItem(itemSlotNumber int) {
	for i := itemSlotNumber; i < 4; i++ {
		if i == 3 {
			e.inventory[i] = NoItem
			return
		}
		e.inventory[i] = e.inventory[i+1]
	}
}

func (e *EntityInventory) GetItemIcon(itemType Item) *ebiten.Image {
	switch itemType {
	case BasicTower:
		return e.basicTowerImage
	case TackTower:
		return e.tackTowerImage
	case IceTower:
		return e.iceTowerImage
	case AoeTower:
		return e.aoeTowerImage
	case ManaTower:
		return e.cashTowerImage
	case FreeUpgrade:
		// TODO:
		return e.damageButtonImage
	case MaxUpgrade:
		return e.firerateButtonImage
	case CurrencyGift:
		return e.removeButtonImage
	case SpikeTrap:
		return e.removeButtonImage
	case ClearEnemies:
		return e.removeButtonImage
	case DamageBuff:
		return e.removeButtonImage
	case SpeedBuff:
		return e.removeButtonImage
	}
	return e.removeButtonImage
}
