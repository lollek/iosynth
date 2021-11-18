package soundServer

import (
	"io"
	"sync"

	"github.com/hajimehoshi/oto/v2"
)

var (
	players []oto.Player
	m       sync.Mutex
	context *oto.Context

	SampleRate      int = 44100
	ChannelNum      int = 2
	BitDepthInBytes int = 2
)

func PlaySound(sound io.Reader) {
	p := context.NewPlayer(sound)
	p.Play()
	m.Lock()
	players = append(players, p)
	m.Unlock()
}

func Init() error {
	c, ready, err := oto.NewContext(SampleRate, ChannelNum, BitDepthInBytes)
	if err != nil {
		return err
	}
	<-ready

	context = c
	return nil
}
