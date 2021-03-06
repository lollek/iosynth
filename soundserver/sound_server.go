package soundserver

import (
	"io"
	"time"

	"github.com/hajimehoshi/oto/v2"
)

var (
	context *oto.Context

	SampleRate      int = 44100
)

func PlaySound(duration time.Duration, sound io.Reader) {
	player := context.NewPlayer(sound)
	player.Play()
	go func() {
		// It seems like we need to manually close all Players, or we
		// get missing notes from time to time. See
		// https://github.com/lollek/iosynth/issues/2
		time.Sleep(duration)
		player.Close()
	}()
}

func Init() error {
	const numberOfChannels = 2 // Stereo
	const bitDepthInBytes = 2 // 16 bit
	c, ready, err := oto.NewContext(SampleRate, numberOfChannels, bitDepthInBytes)
	if err != nil {
		return err
	}
	<-ready

	context = c
	return nil
}
