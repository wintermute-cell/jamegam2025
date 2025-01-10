package entity

import (
	"bufio"
	"image/color"
	"jamegam/pkg/lib"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// Ensure EntityGrid implements Entity
var _ Entity = &EntityGrid{}

type mapTileType int

const (
	mapTileTypeEmpty mapTileType = iota
	mapTileTypePlatform
)

type EntityGrid struct {
	xTiles     int
	yTiles     int
	tilePixels int

	hoveredTile lib.Vec2I

	// Map & Path
	mapDef    string
	enemyPath []lib.Vec2I
	mapTiles  [][]mapTileType

	// Enemies and Towers
	enemies map[lib.Vec2I][]Entity
	towers  map[lib.Vec2I]Entity

	// Resources
	platformImage *ebiten.Image
	floorImage    *ebiten.Image
}

func NewEntityGrid(
	xTiles int,
	yTiles int,
	tilePixels int,
	mapDef string,
	enemyPath []lib.Vec2I,
) *EntityGrid {
	platformImage, _, err := ebitenutil.NewImageFromFile("test_platform.png")
	lib.Must(err)
	floorImage, _, err := ebitenutil.NewImageFromFile("test_floor.png")
	lib.Must(err)
	newEnt := &EntityGrid{
		xTiles:        xTiles,
		yTiles:        yTiles,
		tilePixels:    tilePixels,
		mapDef:        mapDef,
		enemyPath:     enemyPath,
		platformImage: platformImage,
		floorImage:    floorImage,
		enemies:       make(map[lib.Vec2I][]Entity),
		towers:        make(map[lib.Vec2I]Entity),
	}
	return newEnt
}

func (e *EntityGrid) Init(EntitySpawner) {
	scanner := bufio.NewScanner(strings.NewReader(e.mapDef))
	rowCount := 0
	for scanner.Scan() {
		line := scanner.Text()
		rowCount++
		if len(line) != e.xTiles {
			panic("Map definition is not the right size")
		}
		var row []mapTileType
		for _, char := range line {
			switch char {
			case '.':
				row = append(row, mapTileTypeEmpty)
			case 'p':
				row = append(row, mapTileTypePlatform)
			}
		}
		e.mapTiles = append(e.mapTiles, row)
	}
	if rowCount != e.yTiles {
		panic("Map definition is not the right size")
	}
}

func (e *EntityGrid) Update(EntitySpawner) error {
	mouseX, mouseY := ebiten.CursorPosition()
	e.hoveredTile = lib.NewVec2I(mouseX/e.tilePixels, mouseY/e.tilePixels)

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		towerImage, _, err := ebitenutil.NewImageFromFile("test_tower.png")
		lib.Must(err)
		tower := NewEntityTower(
			TowerTypeBasic,
			lib.NewVec2(
				float32(e.hoveredTile.X*e.tilePixels),
				float32(e.hoveredTile.Y*e.tilePixels),
			),
			towerImage)
		e.towers[e.hoveredTile] = tower
	}

	return nil
}

func (e *EntityGrid) Deinit(EntitySpawner) {

}

func (e *EntityGrid) Draw(screen *ebiten.Image) {
	// Draw Grid
	for x := 0; x <= e.xTiles; x++ {
		for y := 0; y <= e.yTiles; y++ {
			drawGridLine(screen, x, y, e.tilePixels)

			// grid lines have to be drawn +1 tile to the right and down
			// so we have to check if we are in bounds
			if x < e.xTiles && y < e.yTiles {
				if e.mapTiles[y][x] == mapTileTypeEmpty {
					geom := ebiten.GeoM{}
					geom.Scale(4, 4)
					geom.Translate(float64(x*e.tilePixels), float64(y*e.tilePixels))
					screen.DrawImage(e.floorImage, &ebiten.DrawImageOptions{
						GeoM: geom,
					})
				}
				if e.mapTiles[y][x] == mapTileTypePlatform {
					geom := ebiten.GeoM{}
					geom.Scale(4, 4)
					geom.Translate(float64(x*e.tilePixels), float64(y*e.tilePixels))
					screen.DrawImage(e.platformImage, &ebiten.DrawImageOptions{
						GeoM: geom,
					})
				}
			}
		}
	}

	// Draw Hovered Tile
	vector.StrokeRect(screen,
		float32(e.hoveredTile.X*e.tilePixels),
		float32(e.hoveredTile.Y*e.tilePixels),
		float32(e.tilePixels),
		float32(e.tilePixels),
		3.0,
		color.RGBA{255, 100, 100, 255},
		false,
	)

	// Draw Enemy Path
	for i := 0; i < len(e.enemyPath)-1; i++ {
		vector.StrokeLine(screen,
			float32(e.enemyPath[i].X*e.tilePixels+e.tilePixels/2),
			float32(e.enemyPath[i].Y*e.tilePixels+e.tilePixels/2),
			float32(e.enemyPath[i+1].X*e.tilePixels+e.tilePixels/2),
			float32(e.enemyPath[i+1].Y*e.tilePixels+e.tilePixels/2),
			3.0,
			color.RGBA{255, 0, 0, 255},
			false)
	}

	// Draw Towers
	for _, tower := range e.towers {
		tower.Draw(screen)
	}
}

func drawGridLine(screen *ebiten.Image, x, y, tilePixels int) {
	var thickness float32 = 1.0
	vector.StrokeLine(screen,
		float32(x*tilePixels),
		float32(y*tilePixels),
		float32(x+1*tilePixels),
		float32(y*tilePixels),
		thickness,
		color.RGBA{255, 255, 255, 255},
		false)
	vector.StrokeLine(screen,
		float32(x*tilePixels),
		float32(y*tilePixels),
		float32(x*tilePixels),
		float32(y+1*tilePixels),
		thickness,
		color.RGBA{255, 255, 255, 255},
		false)
}
