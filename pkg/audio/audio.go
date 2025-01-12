package audio

import (
	"bytes"
	"io"
	"jamegam/pkg/lib"
	"log"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/solarlune/resound/effects"
)

type AudioController struct {
	audioCtx     *audio.Context
	soundPlayers map[string][]PitchedPlayer
	sounds       map[string][]byte
}

type PitchedPlayer struct {
	player       *audio.Player
	pitchShifter *effects.PitchShift
}

func init() {
	Controller = AudioController{
		audioCtx:     audio.NewContext(44100),
		soundPlayers: make(map[string][]PitchedPlayer),
		sounds:       make(map[string][]byte),
	}

	sounds := []string{
		"audio",
		"test_pew",
		"test_deathsound",
		"tower_cash_shot",
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

func (a *AudioController) Play(sound string, variance float64) {
	variance = (rand.Float64() - 0.5) * variance * 2
	if _, ok := a.sounds[sound]; !ok {
		log.Printf("Sound %s not loaded", sound)
	}
	for _, player := range a.soundPlayers[sound] {
		if !player.player.IsPlaying() {
			player.pitchShifter.SetPitch(1.0 + variance)
			player.player.Rewind()
			player.player.Play()
			return
		}
	}

	// reader := bytes.NewReader(a.sounds[sound])
	// player, err := a.audioCtx.NewPlayer(reader)
	// lib.Must(err)
	// a.soundPlayers[sound] = append(a.soundPlayers[sound], player)
	// player.Play()

	reader := bytes.NewReader(a.sounds[sound])
	pshift := effects.NewPitchShift(2048).SetSource(reader).SetPitch(1.0 + variance)
	player, err := a.audioCtx.NewPlayer(pshift)
	lib.Must(err)
	a.soundPlayers[sound] = append(a.soundPlayers[sound], PitchedPlayer{
		player:       player,
		pitchShifter: pshift,
	})
	player.Play()
}

var Controller AudioController
