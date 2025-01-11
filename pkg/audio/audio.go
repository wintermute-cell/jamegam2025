package audio

import (
	"bytes"
	"io"
	"jamegam/pkg/lib"
	"log"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type AudioController struct {
	audioCtx     *audio.Context
	soundPlayers map[string][]*audio.Player
	sounds       map[string][]byte
}

func init() {
	Controller = AudioController{
		audioCtx:     audio.NewContext(44100),
		soundPlayers: make(map[string][]*audio.Player),
		sounds:       make(map[string][]byte),
	}

	sounds := []string{
		"audio",
		"test_pew",
		"test_deathsound",
	}

	for _, sound := range sounds {
		Controller.loadSound(sound)
	}

}

func (a *AudioController) loadSound(name string) {
	soundReader, err := ebitenutil.OpenFile(name + ".ogg")
	lib.Must(err)
	streamReader, err := vorbis.Decode(a.audioCtx, soundReader)
	lib.Must(err)
	a.sounds[name], err = io.ReadAll(streamReader)
	lib.Must(err)
}

func (a *AudioController) Play(sound string) {
	if _, ok := a.sounds[sound]; !ok {
		log.Printf("Sound %s not loaded", sound)
	}
	for _, player := range a.soundPlayers[sound] {
		if !player.IsPlaying() {
			player.Rewind()
			player.Play()
			return
		}
	}

	reader := bytes.NewReader(a.sounds[sound])
	player, err := a.audioCtx.NewPlayer(reader)
	lib.Must(err)
	a.soundPlayers[sound] = append(a.soundPlayers[sound], player)
	player.Play()
}

var Controller AudioController
