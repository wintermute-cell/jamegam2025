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
var _ towers.ProjectileManager = &EntityGrid{}

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
	// enemies     []*enemy.Enemy // TODO: maybe use a free list here too
	enemies     *lib.FreeList[*enemy.Enemy]
	spatialHash *spatialhash.SpatialHash
	towers      map[lib.Vec2I]towers.Tower

	// Projectiles
	projectiles *lib.FreeList[towers.Projectile]

	// Resources
	platformImage *ebiten.Image
	floorImage    *ebiten.Image

	// TODO: REMOVE
	REMOVE_enemyspawntimer float64
}

// AddProjectile implements towers.ProjectileManager.
func (e *EntityGrid) AddProjectile(projectile towers.Projectile) int {
	return e.projectiles.Insert(projectile)
}

// RemoveProjectile implements towers.ProjectileManager.
func (e *EntityGrid) RemoveProjectile(idx int) {
	e.projectiles.Remove(idx)
}

// GetEnemies implements towers.EnemyManager.
// NOTE: MUST BE CALLED AFTER SPATIAL HASH IS CONSTRUCTED
func (e *EntityGrid) GetEnemies(point lib.Vec2, radius float32) ([]*enemy.Enemy, []lib.Vec2I) {
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
		// enemy := e.enemies[idx]
		enemy := e.enemies.Get(int(idx))
		lastIdx, nextIdx := enemy.GetPathNodes()
		last := e.enemyPath[lastIdx].ToVec2().Mul(float32(e.tilePixels))
		next := e.enemyPath[nextIdx].ToVec2().Mul(float32(e.tilePixels))
		pos := last.Lerp(next, float32(enemy.GetPathProgress()))
		if pos.Dist(point) < float32(radius) {
			ret = append(ret, enemy)
		}
	}

	return ret, e.enemyPath
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
		projectiles:   lib.NewFreeList[towers.Projectile](2000),
		enemies:       lib.NewFreeList[*enemy.Enemy](2000),
		platformImage: platformImage,
		floorImage:    floorImage,
		spatialHash:   spatialhash.NewSpatialHash(100_000, int32(tilePixels), 50_000),
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
	e.spatialHash.Clear()

	dt := lib.Dt()
	e.REMOVE_enemyspawntimer += dt
	if e.REMOVE_enemyspawntimer > 0.6 {
		e.REMOVE_enemyspawntimer = 0
		enem := enemy.NewEnemy(enemy.EnemyTypeBasic, 0, 1, 0.0)
		// e.enemies = append(e.enemies, enem)
		idx := e.enemies.Insert(enem)
		enem.SetDestroyFunc(func() {
			e.enemies.Remove(idx)
		})
	}

	// Move Enemies
	shElements := []*spatialhash.SHElement{}
	// for idx, enemy := range e.enemies {
	e.enemies.FuncAll(func(idx int, enemy *enemy.Enemy) {
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
	})

	e.spatialHash.Construct(shElements)

	// Update Towers
	for _, tower := range e.towers {
		tower.Update(e, e)
	}

	// Update Projectiles
	e.projectiles.FuncAll(func(_ int, projectile towers.Projectile) {
		projectile.Update(e, e)
	})

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
	// for _, enem := range e.enemies {
	e.enemies.FuncAll(func(_ int, enem *enemy.Enemy) {
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
	})

	// Draw Towers
	for _, tower := range e.towers {
		tower.Draw(screen)
	}

	// Draw Projectiles
	e.projectiles.FuncAll(func(_ int, projectile towers.Projectile) {
		projectile.Draw(screen)
	})
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
