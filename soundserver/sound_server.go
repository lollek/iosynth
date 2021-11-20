package soundserver

import (
	"io"

	"github.com/hajimehoshi/oto/v2"
)

var (
	context *oto.Context

	SampleRate      int = 44100
	BitDepthInBytes int = 2
)

func PlaySound(sound io.Reader) {
	context.NewPlayer(sound).Play()
}

func Init() error {
	c, ready, err := oto.NewContext(SampleRate, /* number of channels */ 2, BitDepthInBytes)
	if err != nil {
		return err
	}
	<-ready

	context = c
	return nil
}
