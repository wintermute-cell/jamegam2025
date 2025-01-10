package game

import (
	"image/color"
	"jamegam/pkg/entity"
	"jamegam/pkg/lib"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
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
	grid *entity.EntityGrid
}

// NewGame creates a new Game instance
func NewGame() *Game {
	g := &Game{}
	g.Init()
	return g
}

// Init initializes the game.
func (g *Game) Init() {
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
}

// Update is part of the ebiten.Game interface.
func (g *Game) Update() error {
	specialUpdate(g)
	for _, entity := range g.entities {
		if err := entity.Update(g); err != nil {
			return err
		}
	}
	return nil
}

// Draw is part of the ebiten.Game interface.
func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.Black)
	ebitenutil.DebugPrint(screen, "Bye, World!")
	for _, entity := range g.entities {
		entity.Draw(screen)
	}
}

// Layout is part of the ebiten.Game interface.
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return g.tileConfig.width * g.tileConfig.scale, g.tileConfig.height * g.tileConfig.scale
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
