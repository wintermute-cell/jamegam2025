package game

import (
	"image/color"
	"jamegam/pkg/entity"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Game struct {
	firstClickHappened bool
	entities           []entity.Entity
}

// NewGame creates a new Game instance
func NewGame() *Game {
	g := &Game{}
	g.Init()
	return g
}

// Init initializes the game.
func (g *Game) Init() {

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
	return 320, 240
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
			ebiten.SetCursorMode(ebiten.CursorModeCaptured)
		}
	}
	return nil
}
