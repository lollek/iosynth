package sound

import (
	"io"
	"math"
	"time"
	"sync"

	"github.com/lollek/iosynth/soundserver"
)

const (
	numChannels = 2
)

type WaveFn func(phasePosition float64) float64

type SoundWave struct {
	freq   float64
	length int64
	waveFn func(float64) float64
	mutex  sync.Mutex

	pos    int64
}

func SineWaveFn(phasePosition float64) float64 {
	return math.Sin(2 * math.Pi * phasePosition)
}

func SawtoothWaveFn(phasePosition float64) float64 {
	return 1.0 - phasePosition * 2
}

func NewSoundWave(freq float64, duration time.Duration, waveFn WaveFn) *SoundWave {
	bitRate := int64(numChannels) * int64(soundserver.BitDepthInBytes) * int64(soundserver.SampleRate)
	length := (bitRate * int64(duration)) / int64(time.Second)
	if length%4 != 0 {
		length -= length%4
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

	// Shorten the buffer if we are near the end
	if s.pos+int64(len(buf)) >= s.length {
		buf = buf[:s.length-s.pos]
	}

	// Align the buffer to 4 bytes.
	bufLen := len(buf)
	if bufLen%4 > 0 {
		buf = buf[:bufLen - bufLen%4]
		bufLen = len(buf)
	}

	// Fill buffer with data
	length := float64(soundserver.SampleRate) / float64(s.freq)
	bufferSegmentSize := soundserver.BitDepthInBytes * numChannels
	posOnSoundWave := s.pos / int64(bufferSegmentSize)
	numSegmentsInBuffer := bufLen/bufferSegmentSize

	for i := 0; i < numSegmentsInBuffer; i++ {
		const vol = 0.3
		phasePosition := float64(posOnSoundWave)/length
		waveData := s.waveFn(phasePosition) * vol
		switch soundserver.BitDepthInBytes {
		case 1:
			const max = 0x7F
			b := int(waveData * max)
			for ch := 0; ch < numChannels; ch++ {
				buf[i*bufferSegmentSize+ch] = byte(b + 0x80)
			}
		case 2:
			const max = 0x7FFF
			b := int16(waveData * max)
			for ch := 0; ch < numChannels; ch++ {
				buf[i*bufferSegmentSize+2*ch] = byte(b)
				buf[i*bufferSegmentSize+1+2*ch] = byte(b >> 8)
			}
		}
		posOnSoundWave++
	}


	s.pos += int64(bufLen)

	return bufLen, nil
}
