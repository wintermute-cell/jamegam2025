package entity

import (
	"bufio"
	"fmt"
	"image/color"

	// "jamegam/pkg/audio"
	"jamegam/pkg/enemy"
	"jamegam/pkg/lib"
	"jamegam/pkg/pauser"
	"jamegam/pkg/spatialhash"
	"jamegam/pkg/sprites"
	"jamegam/pkg/towers"
	"log"
	"math"
	"math/rand"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
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

	Health int

	messageTimer   float64
	currentMessage string

	towerRangeIndicator bool

	textFace *text.GoTextFace

	// Map & Path
	mapDef    string
	enemyPath []lib.Vec2I
	mapTiles  [][]mapTileType

	// Enemies and Towers
	// enemies     []*enemy.Enemy // TODO: maybe use a free list here too
	enemies       *lib.FreeList[*enemy.Enemy]
	spatialHash   *spatialhash.SpatialHash
	towers        map[lib.Vec2I]towers.Tower
	selectedTower lib.Vec2I // cant have pointers to towers because of map, so only cell
	droppedMana   int64

	// Projectiles
	projectiles *lib.FreeList[towers.Projectile]

	// Resources
	platformImage *ebiten.Image
	floorImage    *ebiten.Image

	// TODO: REMOVE
	REMOVE_enemyspawntimer float64
}

// AddMana implements towers.EnemyManager.
func (e *EntityGrid) AddMana(mana int64) {
	e.droppedMana += mana
}

// AddProjectile implements towers.ProjectileManager.
func (e *EntityGrid) AddProjectile(projectile towers.Projectile) int {
	return e.projectiles.Insert(projectile)
}

// RemoveProjectile implements towers.ProjectileManager.
func (e *EntityGrid) RemoveProjectile(idx int) {
	e.projectiles.Remove(idx)
}

// NukeEnemies destroys all enemies.
func (e *EntityGrid) NukeEnemies() {
	e.enemies.FuncAll(func(idx int, enemy *enemy.Enemy) {
		// TODO: play sound and animation
		enemy.SetHealth(0)
	})
}

// BuffAllTowersDamage increases the damage of all towers temporarily.
func (e *EntityGrid) BuffAllTowersDamage(mult, duration float32) {
	for _, tower := range e.towers {
		tower.SetDamageBuff(mult, duration)
	}
}

// BuffAllTowersSpeed increases the speed of all towers temporarily.
func (e *EntityGrid) BuffAllTowersSpeed(mult, duration float32) {
	for _, tower := range e.towers {
		tower.SetSpeedBuff(mult, duration)
	}
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
		if enemy.IsDead {
			continue
		}
		lastIdx, nextIdx := enemy.GetPathNodes()
		last := e.enemyPath[lastIdx].ToVec2().Mul(float32(e.tilePixels))
		next := e.enemyPath[nextIdx].ToVec2().Mul(float32(e.tilePixels))
		pos := last.Lerp(next, float32(enemy.GetPathProgress())).Add(lib.NewVec2(32, 32))
		var enemyRadius float32 = 24
		if pos.Dist(point) < float32(radius)+enemyRadius {
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
	arialFile, err := ebitenutil.OpenFile("font.ttf")
	lib.Must(err)
	textFaceSource, err := text.NewGoTextFaceSource(arialFile)
	lib.Must(err)

	newEnt := &EntityGrid{
		xTiles:              xTiles,
		yTiles:              yTiles,
		tilePixels:          tilePixels,
		mapDef:              mapDef,
		enemyPath:           enemyPath,
		projectiles:         lib.NewFreeList[towers.Projectile](2000),
		enemies:             lib.NewFreeList[*enemy.Enemy](2000),
		platformImage:       platformImage,
		floorImage:          floorImage,
		spatialHash:         spatialhash.NewSpatialHash(100_000, int32(tilePixels), 50_000),
		textFace:            &text.GoTextFace{Source: textFaceSource, Size: 20},
		towers:              make(map[lib.Vec2I]towers.Tower),
		droppedMana:         0,
		towerRangeIndicator: true,
		selectedTower:       lib.NewVec2I(-1, -1),
		Health:              100,
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

func (e *EntityGrid) SpawnEnemy(enType enemy.EnemyType) {
	enem := enemy.NewEnemy(enType, 0, 1, 0.0)
	enValue := enem.GetValue()
	idx := e.enemies.Insert(enem)
	enem.SetDestroyFunc(func() {
		e.enemies.Remove(idx)
		e.droppedMana += enValue
	})
}

func (e *EntityGrid) ShowMessage(message string) {
	// audio.Controller.Play("error", 0.0) // TODO: notification sound?
	e.currentMessage = "> " + message
	e.messageTimer = 3.0
}

func (e *EntityGrid) Update(EntitySpawner) error {

	if e.Health <= 0 {
		e.Restart()
		return nil
	}

	e.messageTimer -= lib.Dt()
	if e.messageTimer <= 0 {
		e.messageTimer = -2
		e.currentMessage = ""
	}

	e.spatialHash.Clear()

	dt := lib.Dt()

	// Move Enemies
	shElements := []*spatialhash.SHElement{}
	// for idx, enemy := range e.enemies {
	killedEnemies := []int{}
	e.enemies.FuncAll(func(idx int, enemy *enemy.Enemy) {
		lastIdx, nextIdx := enemy.GetPathNodes()
		progress := enemy.GetPathProgress()
		progress += float64(enemy.GetSpeed()) * dt
		enemy.SetPathProgress(progress)

		if progress >= 1.0 {
			progress = 0
			if nextIdx == len(e.enemyPath)-2 {
				killedEnemies = append(killedEnemies, idx)
				log.Println("Enemy reached the end")
				e.Health--

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

	for _, idx := range killedEnemies {
		e.enemies.Get(idx).SetHealth(0)
	}

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
	geom := ebiten.GeoM{}
	geom.Scale(4, 4)
	screen.DrawImage(sprites.SpriteMap, &ebiten.DrawImageOptions{
		GeoM: geom,
	})

	// for x := 0; x <= e.xTiles; x++ {
	// 	for y := 0; y <= e.yTiles; y++ {
	// 		drawGridLine(screen, x, y, e.tilePixels)
	//
	// 		// // grid lines have to be drawn +1 tile to the right and down
	// 		// // so we have to check if we are in bounds
	// 		// if x < e.xTiles && y < e.yTiles {
	// 		// 	if e.mapTiles[y][x] == mapTileTypeEmpty {
	// 		// 		geom := ebiten.GeoM{}
	// 		// 		geom.Scale(4, 4)
	// 		// 		geom.Translate(float64(x*e.tilePixels), float64(y*e.tilePixels))
	// 		// 		screen.DrawImage(e.floorImage, &ebiten.DrawImageOptions{
	// 		// 			GeoM: geom,
	// 		// 		})
	// 		// 	}
	// 		// 	if e.mapTiles[y][x] == mapTileTypePlatform {
	// 		// 		geom := ebiten.GeoM{}
	// 		// 		geom.Scale(4, 4)
	// 		// 		geom.Translate(float64(x*e.tilePixels), float64(y*e.tilePixels))
	// 		// 		screen.DrawImage(e.platformImage, &ebiten.DrawImageOptions{
	// 		// 			GeoM: geom,
	// 		// 		})
	// 		// 	}
	// 		// }
	// 	}
	// }

	// Draw Enemy Path
	// for i := 0; i < len(e.enemyPath)-1; i++ {
	// 	vector.StrokeLine(screen,
	// 		float32(e.enemyPath[i].X*e.tilePixels+e.tilePixels/2),
	// 		float32(e.enemyPath[i].Y*e.tilePixels+e.tilePixels/2),
	// 		float32(e.enemyPath[i+1].X*e.tilePixels+e.tilePixels/2),
	// 		float32(e.enemyPath[i+1].Y*e.tilePixels+e.tilePixels/2),
	// 		3.0,
	// 		color.RGBA{255, 0, 0, 255},
	// 		false)
	// }

	// Draw Enemies
	e.enemies.FuncAll(func(_ int, enem *enemy.Enemy) {
		progress := enem.GetPathProgress()
		lastIdx, nextIdx := enem.GetPathNodes()
		last := e.enemyPath[lastIdx].ToVec2().Mul(float32(e.tilePixels))
		next := e.enemyPath[nextIdx].ToVec2().Mul(float32(e.tilePixels))
		pos := last.Lerp(next, float32(progress))

		newWander := enem.GetWander() + float32(lib.Dt())*enem.WanderVelocity
		enem.WanderVelocity = enem.WanderVelocity*0.95 + (rand.Float32()-0.5)*200*float32(lib.Dt())
		newWander = max(-10, min(10, newWander))
		if !pauser.IsPaused {
			enem.SetWander(newWander)
		}
		wanderDirection := next.Sub(last).Normalize().Rotate(90).Mul(enem.GetWander())

		newBounce := enem.GetBounce() + float32(lib.Dt())*float32(math.Sqrt(float64(enem.GetSpeed())))*10
		if !pauser.IsPaused {
			enem.SetBounce(newBounce)
		}

		geom := ebiten.GeoM{}
		geom.Scale(4, 4)
		geom.Translate(float64(pos.X), float64(pos.Y))
		geom.Translate(float64(wanderDirection.X), float64(wanderDirection.Y))
		geom.Translate(0, math.Sin(float64(newBounce))*5)
		screen.DrawImage(enem.GetSprite(), &ebiten.DrawImageOptions{
			GeoM: geom,
		})
		if enem.GetSpeedMod() < 1.0 {
			// Draw a slow effect
			screen.DrawImage(enemy.SpriteSlowEffect, &ebiten.DrawImageOptions{
				GeoM: geom,
			})
		}

		// Draw Hitbox
		// vector.StrokeCircle(screen,
		// 	float32(pos.X)+32,
		// 	float32(pos.Y)+32,
		// 	24,
		// 	1,
		// 	color.RGBA{255, 0, 0, 255},
		// 	false)

		geom.Reset()
		geom.Scale(4, 4)
	})

	screen.DrawImage(sprites.SpriteOverMap, &ebiten.DrawImageOptions{
		GeoM: geom,
	})

	// Draw Towers
	for _, tower := range e.towers {
		tower.Draw(screen)
	}
	if e.selectedTower.X >= 0 && e.selectedTower.Y >= 0 {
		selectedTower := e.towers[e.selectedTower]
		if selectedTower != nil {
			radius := selectedTower.Radius()
			vector.DrawFilledCircle(screen,
				float32(e.selectedTower.X*64+32),
				float32(e.selectedTower.Y*64+32),
				radius,
				color.RGBA{0, 0, 0, 80},
				false)
			selectedTower.Draw(screen) // draw that one again on top
		}
	}

	// Draw Projectiles
	e.projectiles.FuncAll(func(_ int, projectile towers.Projectile) {
		projectile.Draw(screen)
	})

	// Draw Message
	if e.messageTimer > -0.5 {
		geom := ebiten.GeoM{}
		geom.Translate(10, (12*64)-34)
		vector.DrawFilledRect(screen, 5, 12*64-39, 16*64-10, 34, color.RGBA{0, 0, 0, 160}, false)
		text.Draw(screen, fmt.Sprintf("%v", e.currentMessage), e.textFace, &text.DrawOptions{
			DrawImageOptions: ebiten.DrawImageOptions{GeoM: geom},
		})
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

func (e *EntityGrid) Restart() {
	e.NukeEnemies()
	e.projectiles.Clear()
	e.selectedTower = lib.NewVec2I(-1, -1)
	e.towers = make(map[lib.Vec2I]towers.Tower)
	e.Health = 100
}
