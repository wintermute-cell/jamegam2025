package audio

import (
	"bytes"
	"io"
	"jamegam/pkg/lib"
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/solarlune/resound/effects"
)

type AudioController struct {
	audioCtx     *audio.Context
	soundPlayers map[string][]PitchedPlayer
	sounds       map[string][]byte

	ostPlayer *audio.Player
	isMuted   bool
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
		"enemy_death_poof",
		"tower_cash_shot",
		"tower_cash_shot2",
		"basic_tower_shoot",
		"ice_tower_shoot",
		"aoe_tower_shoot",
		"aoe_tower_explosion",
		"click",
		"error",
		"notification",
		"build_tower",
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

func (a *AudioController) PlayOst() {
	ostReader, err := ebitenutil.OpenFile("ost.ogg")
	lib.Must(err)
	streamReader, err := vorbis.Decode(a.audioCtx, ostReader)
	lib.Must(err)
	loop := audio.NewInfiniteLoop(streamReader, streamReader.Length())
	player, err := a.audioCtx.NewPlayer(loop)
	lib.Must(err)
	player.SetVolume(0.5)
	player.Play()
	a.ostPlayer = player
	// player, err := a.audioCtx.NewPlayer(streamReader)
	// lib.Must(err)
	// audio.NewInfiniteLoop(player, streamReader.Length()).Play()
	// player.SetVolume(0.5)
	// player.Play()
}

func (a *AudioController) ToggleMute() {
	a.isMuted = !a.isMuted
	if a.isMuted {
		a.ostPlayer.SetVolume(0)
	} else {
		a.ostPlayer.SetVolume(0.5)
	}
}

func (a *AudioController) Play(sound string, variance float64) {
	if a.isMuted {
		return
	}
	variance = (rand.Float64() - 0.5) * variance * 2
	log.Printf("Playing sound %s with variance %f", sound, variance)
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
	pshift := effects.NewPitchShift(256).SetSource(reader).SetStrength(0.7).SetPitch(1.0 + variance)
	player, err := a.audioCtx.NewPlayer(pshift)
	player.SetBufferSize(time.Millisecond * 500)
	lib.Must(err)
	a.soundPlayers[sound] = append(a.soundPlayers[sound], PitchedPlayer{
		player:       player,
		pitchShifter: pshift,
	})
	player.Play()
}

var Controller AudioController
