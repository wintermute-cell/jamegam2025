package sprites

import (
	"jamegam/pkg/lib"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var (
	SpriteMap     *ebiten.Image
	SpriteOverMap *ebiten.Image

	SpritePauseMenu *ebiten.Image
)

func init() {
	var err error

	SpriteMap, _, err = ebitenutil.NewImageFromFile("map.png")
	lib.Must(err)

	SpriteOverMap, _, err = ebitenutil.NewImageFromFile("over_map.png")
	lib.Must(err)

	SpritePauseMenu, _, err = ebitenutil.NewImageFromFile("pausemenu.png")
	lib.Must(err)
}
