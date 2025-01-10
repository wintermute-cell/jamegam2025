package entity

import (
	"bufio"
	"image/color"
	"jamegam/pkg/enemy"
	"jamegam/pkg/lib"
	"jamegam/pkg/towers"
	"log"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// Ensure EntityGrid implements Entity
var _ Entity = &EntityGrid{}
var _ towers.EnemyManager = &EntityGrid{}

type mapTileType int

const (
	mapTileTypeEmpty mapTileType = iota
	mapTileTypePlatform
)

type EntityGrid struct {
	xTiles     int
	yTiles     int
	tilePixels int

	hoveredTile         lib.Vec2I
	hoveredTileHasTower bool

	// Map & Path
	mapDef    string
	enemyPath []lib.Vec2I
	mapTiles  [][]mapTileType

	// Enemies and Towers
	enemies    map[lib.Vec2I][]*enemy.Enemy
	newEnemies map[lib.Vec2I][]*enemy.Enemy // Used for updating of enemy positions
	towers     map[lib.Vec2I]towers.Tower

	// Resources
	platformImage *ebiten.Image
	floorImage    *ebiten.Image

	// TODO: REMOVE
	REMOVE_enemyspawntimer float64
}

// GetEnemies implements towers.EnemyManager.
func (e *EntityGrid) GetEnemies(point lib.Vec2, radius int) []enemy.Enemy {
	panic("unimplemented")
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
		enemies:       make(map[lib.Vec2I][]*enemy.Enemy),
		newEnemies:    make(map[lib.Vec2I][]*enemy.Enemy),
		towers:        make(map[lib.Vec2I]towers.Tower),
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
	dt := lib.Dt()
	e.REMOVE_enemyspawntimer += dt
	if e.REMOVE_enemyspawntimer > 5.0 {
		e.REMOVE_enemyspawntimer = 0
		enem := enemy.NewEnemy(enemy.EnemyTypeBasic, 0, 1, 0.0)
		e.enemies[lib.NewVec2I(0, 1)] = append(e.enemies[lib.NewVec2I(0, 1)], enem)
	}

	// Tower Placement
	mouseX, mouseY := ebiten.CursorPosition()
	e.hoveredTile = lib.NewVec2I(mouseX/e.tilePixels, mouseY/e.tilePixels)
	_, e.hoveredTileHasTower = e.towers[e.hoveredTile]
	if !e.hoveredTileHasTower {
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			tower := towers.NewTowerBasic(e.hoveredTile.Mul(e.tilePixels))
			e.towers[e.hoveredTile] = tower
		}
	}

	// Move Enemies
	clear(e.newEnemies)
	for _, cell := range e.enemies {
		for _, enemy := range cell {
			lastIdx, nextIdx := enemy.GetPathNodes()
			progress := enemy.GetPathProgress()
			progress += 1.0 * dt
			log.Println("Enemy progress", progress)
			enemy.SetPathProgress(progress)

			if progress >= 0.5 {
				// move to next cell
				e.newEnemies[e.enemyPath[nextIdx]] = append(e.newEnemies[e.enemyPath[nextIdx]], enemy)
			} else {
				e.newEnemies[e.enemyPath[lastIdx]] = append(e.newEnemies[e.enemyPath[lastIdx]], enemy)
			}

			if progress >= 1.0 {
				progress = 0
				if nextIdx == len(e.enemyPath)-1 {
					log.Println("Enemy reached the end")
					panic("unimplemented")
				}
				enemy.SetPathNodes(lastIdx+1, nextIdx+1)
				enemy.SetPathProgress(progress)
			}
		}
	}
	// e.enemies = e.newEnemies
	clear(e.enemies)
	for k, v := range e.newEnemies {
		e.enemies[k] = v
	}

	// Update Towers
	for _, tower := range e.towers {
		tower.Update(e)
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

	// Draw Enemies
	for _, cell := range e.enemies {
		for _, enem := range cell {
			progress := enem.GetPathProgress()
			lastIdx, nextIdx := enem.GetPathNodes()
			last := e.enemyPath[lastIdx].ToVec2()
			next := e.enemyPath[nextIdx].ToVec2()
			pos := last.Lerp(next, float32(progress))

			geom := ebiten.GeoM{}
			geom.Scale(4, 4)
			geom.Translate(float64(pos.X*float32(e.tilePixels)), float64(pos.Y*float32(e.tilePixels)))
			screen.DrawImage(enemy.SpriteEnemyBasic, &ebiten.DrawImageOptions{
				GeoM: geom,
			})

			geom.Reset()
			geom.Scale(4, 4)
		}
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
