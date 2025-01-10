package entity

import (
	"bufio"
	"image/color"
	"jamegam/pkg/enemy"
	"jamegam/pkg/lib"
	"jamegam/pkg/spatialhash"
	"jamegam/pkg/towers"
	"log"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
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
	// enemies    map[lib.Vec2I][]*enemy.Enemy
	// newEnemies map[lib.Vec2I][]*enemy.Enemy // Used for updating of enemy positions
	enemies     []*enemy.Enemy
	spatialHash *spatialhash.SpatialHash
	towers      map[lib.Vec2I]towers.Tower

	// Resources
	platformImage *ebiten.Image
	floorImage    *ebiten.Image

	// TODO: REMOVE
	REMOVE_enemyspawntimer float64
}

// GetEnemies implements towers.EnemyManager.
// NOTE: MUST BE CALLED AFTER SPATIAL HASH IS CONSTRUCTED
func (e *EntityGrid) GetEnemies(point lib.Vec2, radius float32) []*enemy.Enemy {
	ret := []*enemy.Enemy{}
	shBounds := spatialhash.SHBounds{
		Mx:      int32(point.X),
		My:      int32(point.Y),
		HWidth:  int32(radius),
		HHeight: int32(radius),
	}
	hitIdxs := e.spatialHash.InBounds(shBounds)
	for _, hit := range hitIdxs {
		idx := hit.ID
		enemy := e.enemies[idx]
		lastIdx, nextIdx := enemy.GetPathNodes()
		last := e.enemyPath[lastIdx].ToVec2().Mul(float32(e.tilePixels))
		next := e.enemyPath[nextIdx].ToVec2().Mul(float32(e.tilePixels))
		pos := last.Lerp(next, float32(enemy.GetPathProgress()))
		if pos.Dist(point) < float32(radius) {
			ret = append(ret, enemy)
		}
	}

	return ret
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
		spatialHash:   spatialhash.NewSpatialHash(100_000, int32(tilePixels), 50_000),
		// enemies:       make(map[lib.Vec2I][]*enemy.Enemy),
		// newEnemies:    make(map[lib.Vec2I][]*enemy.Enemy),
		towers: make(map[lib.Vec2I]towers.Tower),
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
	e.spatialHash.Clear()

	dt := lib.Dt()
	e.REMOVE_enemyspawntimer += dt
	if e.REMOVE_enemyspawntimer > 5.0 {
		e.REMOVE_enemyspawntimer = 0
		enem := enemy.NewEnemy(enemy.EnemyTypeBasic, 0, 1, 0.0)
		e.enemies = append(e.enemies, enem)
	}

	// Move Enemies
	shElements := []*spatialhash.SHElement{}
	for idx, enemy := range e.enemies {
		lastIdx, nextIdx := enemy.GetPathNodes()
		progress := enemy.GetPathProgress()
		progress += 1.0 * dt
		enemy.SetPathProgress(progress)

		if progress >= 1.0 {
			progress = 0
			if nextIdx == len(e.enemyPath)-1 {
				log.Println("Enemy reached the end")
				panic("unimplemented")
			}
			enemy.SetPathNodes(lastIdx+1, nextIdx+1)
			enemy.SetPathProgress(progress)
			enemy.SetNumPassedNodes(enemy.GetNumPassedNodes() + 1.0)
		}

		shElements = append(shElements, &spatialhash.SHElement{
			ID: int32(idx),
			Bounds: spatialhash.SHBounds{
				Mx:      int32(e.enemyPath[nextIdx].X*e.tilePixels + (e.tilePixels / 2)),
				My:      int32(e.enemyPath[nextIdx].Y*e.tilePixels + (e.tilePixels / 2)),
				HWidth:  int32(e.tilePixels / 2),
				HHeight: int32(e.tilePixels / 2),
			},
		})

	}

	e.spatialHash.Construct(shElements)

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
	for _, enem := range e.enemies {
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

		// for debugging, draw line from enemy to the node that contains the enemy
		// vector.StrokeLine(screen,
		// 	float32(pos.X*float32(e.tilePixels)+float32(e.tilePixels)/2),
		// 	float32(pos.Y*float32(e.tilePixels)+float32(e.tilePixels)/2),
		// 	float32(cellKey.X*e.tilePixels+e.tilePixels/2),
		// 	float32(cellKey.Y*e.tilePixels+e.tilePixels/2),
		// 	1.0,
		// 	color.RGBA{0, 0, 255, 255},
		// 	false)
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
