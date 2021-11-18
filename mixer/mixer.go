package mixer

import (
	"time"

	"github.com/lollek/iosynth/soundserver"
	"github.com/lollek/iosynth/sound"
)

func RawCommand(data []byte) {
	freq := 523.3
	duration := 3 * time.Second
	sound := sound.NewSineWave(freq, duration)
	soundserver.PlaySound(sound)
}

