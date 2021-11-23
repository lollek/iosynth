package sound

import (
	"io"
	"math"
	"math/rand"
	"sync"
	"time"

	"github.com/lollek/iosynth/soundserver"
)

const (
	numChannels = 2
)

type WaveFn func(tick, ticksPerCycle float64) float64

type SoundWave struct {
	freq   float64
	length int64
	waveFn WaveFn
	mutex  sync.Mutex

	pos int64
}

func SineWaveFn(tick, ticksPerCycle float64) float64 {
	return math.Sin(2 * math.Pi * (tick / ticksPerCycle))
}

func SquareWaveFn(tick, ticksPerCycle float64) float64 {
	halfACycle := int(ticksPerCycle) / 2
	tickOfCycle := int(tick) % int(ticksPerCycle)
	if tickOfCycle < halfACycle {
		return 0.5
	} else {
		return -0.5
	}
}

func TriangleWaveFn(tick, ticksPerCycle float64) float64 {
	halfACycle := int(ticksPerCycle) / 2
	tickOfCycle := int(tick) % int(ticksPerCycle)
	percentageOfCycle := float64(tickOfCycle) / ticksPerCycle
	if tickOfCycle < halfACycle {
		return (-1 + percentageOfCycle*4)
	} else {
		return (1 - (percentageOfCycle-0.5)*4)
	}
}

func SawtoothWaveFn(tick, ticksPerCycle float64) float64 {
	tickOfCycle := int(tick) % int(ticksPerCycle)
	percentageOfCycle := float64(tickOfCycle) / ticksPerCycle
	return 1 - percentageOfCycle*2
}

func NoiseWaveFn(tick, ticksPerCycle float64) float64 {
	return 1 - rand.Float64() + rand.Float64()
}

func NewSoundWave(freq float64, duration time.Duration, waveFn WaveFn) *SoundWave {
	bitRate := int64(numChannels) * int64(soundserver.BitDepthInBytes) * int64(soundserver.SampleRate)
	length := (bitRate * int64(duration)) / int64(time.Second)
	if length%4 != 0 {
		length += 4 - length%4
	}
	return &SoundWave{
		freq:   freq,
		length: length,
		waveFn: waveFn,
	}
}

func (s *SoundWave) Read(buf []byte) (bytesRead int, err error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.pos >= s.length {
		return 0, io.EOF
	}

	// Shorten the buffer if we don't need to calculate data for the whole
	// buffer (e.g. if we are nearing the end)
	if s.pos+int64(len(buf)) >= s.length {
		buf = buf[:s.length-s.pos]
	}

	// Align the buffer to 4 bytes.
	bufLen := len(buf)
	if bufLen%4 > 0 {
		buf = buf[:bufLen-bufLen%4]
		bufLen = len(buf)
	}

	// Fill buffer with data
	ticksPerCycle := float64(soundserver.SampleRate) / float64(s.freq)
	tickDataSize := soundserver.BitDepthInBytes * numChannels
	tick := s.pos / int64(tickDataSize)
	ticksInBuffer := bufLen / tickDataSize

	for i := 0; i < ticksInBuffer; i++ {
		const amp = 0.3
		waveData := s.waveFn(float64(tick), ticksPerCycle) * amp
		switch soundserver.BitDepthInBytes {
		case 1:
			const max = 0x7F
			b := int(waveData * max)
			for ch := 0; ch < numChannels; ch++ {
				buf[i*tickDataSize+ch] = byte(b + 0x80)
			}
		case 2:
			const max = 0x7FFF
			b := int16(waveData * max)
			for ch := 0; ch < numChannels; ch++ {
				buf[i*tickDataSize+2*ch] = byte(b)
				buf[i*tickDataSize+1+2*ch] = byte(b >> 8)
			}
		}
		tick++
	}

	s.pos += int64(bufLen)

	return bufLen, nil
}
