package main

import (
	"io"
	"sync"

	"github.com/hajimehoshi/oto/v2"
)

var (
	players []oto.Player
	m       sync.Mutex
	context *oto.Context
)

func PlaySound(sound io.Reader) {
	p := context.NewPlayer(sound)
	p.Play()
	m.Lock()
	players = append(players, p)
	m.Unlock()
}

func InitSoundServer() error {
	c, ready, err := oto.NewContext(sampleRate, channelNum, bitDepthInBytes)
	if err != nil {
		return err
	}
	<-ready

	context = c
	return nil
}
