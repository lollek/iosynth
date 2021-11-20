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

type SoundWave struct {
	freq   float64
	length int64
	mutex  sync.Mutex

	pos    int64
}

func NewSoundWave(freq float64, duration time.Duration) *SoundWave {
	bitRate := int64(numChannels) * int64(soundserver.BitDepthInBytes) * int64(soundserver.SampleRate)
	length := (bitRate * int64(duration)) / int64(time.Second)
	if length%4 != 0 {
		length -= length%4
	}
	return &SoundWave{
		freq:   freq,
		length: length,
	}
}

func (s *SoundWave) Read(buf []byte) (int, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.pos >= s.length {
		return 0, io.EOF
	}

	eof := false
	if s.pos+int64(len(buf)) >= s.length {
		buf = buf[:s.length-s.pos]
		eof = true
	}

	bufLen := len(buf)
	if bufLen%4 > 0 {
		buf = buf[:bufLen - bufLen%4]
		bufLen = len(buf)
	}

	length := float64(soundserver.SampleRate) / float64(s.freq)

	bufferSegmentSize := soundserver.BitDepthInBytes * numChannels
	posOnSoundWave := s.pos / int64(bufferSegmentSize)
	numSegmentsInBuffer := bufLen/bufferSegmentSize
	switch soundserver.BitDepthInBytes {
	case 1:
		for i := 0; i < numSegmentsInBuffer; i++ {
			const max = 127
			b := int(math.Sin(2*math.Pi*float64(posOnSoundWave)/length) * 0.3 * max)
			for ch := 0; ch < numChannels; ch++ {
				buf[i*bufferSegmentSize+ch] = byte(b + 128)
			}
			posOnSoundWave++
		}
	case 2:
		for i := 0; i < numSegmentsInBuffer; i++ {
			const max = 32767
			b := int16(math.Sin(2*math.Pi*float64(posOnSoundWave)/length) * 0.3 * max)
			for ch := 0; ch < numChannels; ch++ {
				buf[i*bufferSegmentSize+2*ch] = byte(b)
				buf[i*bufferSegmentSize+1+2*ch] = byte(b >> 8)
			}
			posOnSoundWave++
		}
	}

	s.pos += int64(bufLen)

	if eof {
		return bufLen, io.EOF
	}
	return bufLen, nil
}
