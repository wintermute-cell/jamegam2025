package sprites

import (
	"jamegam/pkg/lib"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var (
	SpriteMap     *ebiten.Image
	SpriteOverMap *ebiten.Image

	SpritePauseMenu      *ebiten.Image
	SpriteMainMenu       *ebiten.Image
	SpriteMainMenuButton *ebiten.Image
	SpriteTutorial       *ebiten.Image
)

func init() {
	var err error

	SpriteMap, _, err = ebitenutil.NewImageFromFile("map.png")
	lib.Must(err)

	SpriteOverMap, _, err = ebitenutil.NewImageFromFile("over_map.png")
	lib.Must(err)

	SpritePauseMenu, _, err = ebitenutil.NewImageFromFile("pausemenu.png")
	lib.Must(err)

	SpriteMainMenu, _, err = ebitenutil.NewImageFromFile("mainmenu.png")
	lib.Must(err)

	SpriteMainMenuButton, _, err = ebitenutil.NewImageFromFile("mainmenu_button.png")
	lib.Must(err)

	SpriteTutorial, _, err = ebitenutil.NewImageFromFile("tutorial.png")
	lib.Must(err)
}
