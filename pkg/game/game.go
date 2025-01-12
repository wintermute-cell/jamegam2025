package game

import (
	"fmt"
	"image/color"
	"jamegam/pkg/audio"
	"jamegam/pkg/entity"
	"jamegam/pkg/lib"
	"jamegam/pkg/pauser"
	"jamegam/pkg/sprites"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type TileConfig struct {
	width  int
	height int
	scale  int
}

type Game struct {
	tileConfig         TileConfig
	firstClickHappened bool
	entities           []entity.Entity

	// Entities
	grid      *entity.EntityGrid
	inventory *entity.EntityInventory

	isMainMenu bool
}

// NewGame creates a new Game instance
func NewGame() *Game {
	g := &Game{}
	g.Init()
	return g
}

// Init initializes the game.
func (g *Game) Init() {
	audio.Controller.PlayOst()
	g.tileConfig = TileConfig{16, 12, 64}
	mapDef := `
pppppppppppppppp
pppppp........pp
p....p.pppppp.pp
p.pp.p..pp....pp
..pp.pp.pp.ppppp
pppp..p..p.....p
ppppp.pp.ppppp.p
pp....p..p.....p
pp.pppp.pp.ppppp
pp..p...pp.p...p
ppp...pppp...p.p
pppppppppppppp.p
`
	enemyPath := []lib.Vec2I{lib.NewVec2I(14, 12), lib.NewVec2I(14, 11), lib.NewVec2I(14, 10), lib.NewVec2I(14, 9), lib.NewVec2I(13, 9), lib.NewVec2I(12, 9), lib.NewVec2I(12, 10), lib.NewVec2I(11, 10), lib.NewVec2I(10, 10), lib.NewVec2I(10, 9), lib.NewVec2I(10, 8), lib.NewVec2I(10, 7), lib.NewVec2I(11, 7), lib.NewVec2I(12, 7), lib.NewVec2I(13, 7), lib.NewVec2I(14, 7), lib.NewVec2I(14, 6), lib.NewVec2I(14, 5), lib.NewVec2I(13, 5), lib.NewVec2I(12, 5), lib.NewVec2I(11, 5), lib.NewVec2I(10, 5), lib.NewVec2I(10, 4), lib.NewVec2I(10, 3), lib.NewVec2I(11, 3), lib.NewVec2I(12, 3), lib.NewVec2I(13, 3), lib.NewVec2I(13, 2), lib.NewVec2I(13, 1), lib.NewVec2I(12, 1), lib.NewVec2I(11, 1), lib.NewVec2I(10, 1), lib.NewVec2I(9, 1), lib.NewVec2I(8, 1), lib.NewVec2I(7, 1), lib.NewVec2I(6, 1), lib.NewVec2I(6, 2), lib.NewVec2I(6, 3), lib.NewVec2I(7, 3), lib.NewVec2I(7, 4), lib.NewVec2I(7, 5), lib.NewVec2I(8, 5), lib.NewVec2I(8, 6), lib.NewVec2I(8, 7), lib.NewVec2I(7, 7), lib.NewVec2I(7, 8), lib.NewVec2I(7, 9), lib.NewVec2I(6, 9), lib.NewVec2I(5, 9), lib.NewVec2I(5, 10), lib.NewVec2I(4, 10), lib.NewVec2I(3, 10), lib.NewVec2I(3, 9), lib.NewVec2I(2, 9), lib.NewVec2I(2, 8), lib.NewVec2I(2, 7), lib.NewVec2I(3, 7), lib.NewVec2I(4, 7), lib.NewVec2I(5, 7), lib.NewVec2I(5, 6), lib.NewVec2I(5, 5), lib.NewVec2I(4, 5), lib.NewVec2I(4, 4), lib.NewVec2I(4, 3), lib.NewVec2I(4, 2), lib.NewVec2I(3, 2), lib.NewVec2I(2, 2), lib.NewVec2I(1, 2), lib.NewVec2I(1, 3), lib.NewVec2I(1, 4), lib.NewVec2I(0, 4), lib.NewVec2I(-1, 4)}

	mapDef = strings.TrimSpace(mapDef) // remove leading and trailing whitespace
	g.grid = entity.NewEntityGrid(g.tileConfig.width, g.tileConfig.height, g.tileConfig.scale, mapDef, enemyPath)
	g.AddEntity(g.grid)
	g.inventory = entity.NewEntityInventory(g.tileConfig.scale, g.grid)
	g.AddEntity(g.inventory)
}

// Called by main menu
func (g *Game) LateInit() error
}

// Update is part of the ebiten.Game interface.
func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		pauser.IsPaused = !pauser.IsPaused
	}
	specialUpdate(g)
	if !pauser.IsPaused {
		for _, entity := range g.entities {
			if err := entity.Update(g); err != nil {
				return err
			}
		}
	} else {
		// TODO: pause menu buttons
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			x, y := ebiten.CursorPosition()
			// vector.DrawFilledRect(fakeScreen, 32+312, 20+300, 336, 88, color.RGBA{0, 0, 0, 100}, false)
			// vector.DrawFilledRect(fakeScreen, 32+312, 136+300, 336, 88, color.RGBA{0, 0, 0, 100}, false)
			if x > 312+32 && x < 312+32+336 && y > 300+20 && y < 300+20+88 {
				g.inventory.RestartGame()
			}
			if x > 312+32 && x < 312+32+336 && y > 300+136 && y < 300+136+88 {
				audio.Controller.ToggleMute()
			}

		}
	}
	return nil
}

// Draw is part of the ebiten.Game interface.
func (g *Game) Draw(screen *ebiten.Image) {
	fakeScreen := ebiten.NewImage(1024, 1024)
	screen.Fill(color.Black)
	fakeScreen.Fill(color.Black)

	for _, entity := range g.entities {
		entity.Draw(fakeScreen)
	}

	if pauser.IsPaused {
		vector.DrawFilledRect(fakeScreen, 0, 0, 2000, 2000, color.RGBA{0, 0, 0, 100}, false)
		geom := ebiten.GeoM{}
		geom.Translate(312, 300)
		fakeScreen.DrawImage(sprites.SpritePauseMenu, &ebiten.DrawImageOptions{
			GeoM: geom,
		})

		// Button hitboxes
		// vector.DrawFilledRect(fakeScreen, 32+312, 20+300, 336, 88, color.RGBA{0, 0, 0, 100}, false)
		// vector.DrawFilledRect(fakeScreen, 32+312, 136+300, 336, 88, color.RGBA{0, 0, 0, 100}, false)
	}

	ebitenutil.DebugPrint(fakeScreen, fmt.Sprintf("dt: %f", lib.Dt()))
	geom := ebiten.GeoM{}
	// geom.Translate(20, 20)
	screen.DrawImage(fakeScreen, &ebiten.DrawImageOptions{
		GeoM: geom,
	})
}

// Layout is part of the ebiten.Game interface.
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	// return (g.tileConfig.width)*g.tileConfig.scale + 40, (g.tileConfig.height+2)*g.tileConfig.scale + 40
	return (outsideWidth), (outsideHeight)
}

// AddEntity adds an entity to the game
func (g *Game) AddEntity(e entity.Entity) {
	e.Init(g)
	g.entities = append(g.entities, e)
}

// RemoveEntity removes an entity from the game
func (g *Game) RemoveEntity(e entity.Entity) {
	for i, entity := range g.entities {
		if entity == e {
			entity.Deinit(g)
			g.entities = append(g.entities[:i], g.entities[i+1:]...)
			return
		}
	}
}

// specialUpdate is a part of update that does not contain game specific logic,
// but general behaviour handling
func specialUpdate(g *Game) error {
	if !g.firstClickHappened {
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			g.firstClickHappened = true
			// ebiten.SetCursorMode(ebiten.CursorModeCaptured)
		}
	}
	return nil
}
