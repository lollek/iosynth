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
	remaining []byte
}

func NewSoundWave(freq float64, duration time.Duration) *SoundWave {
	bitRate := int64(numChannels) * int64(soundserver.BitDepthInBytes) * int64(soundserver.SampleRate)

	return &SoundWave{
		freq:   freq,
		length: (((bitRate * int64(duration)) / int64(time.Second)) / 4) * 4,
	}
}

func (s *SoundWave) Read(buf []byte) (int, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if len(s.remaining) > 0 {
		n := copy(buf, s.remaining)
		s.remaining = s.remaining[n:]
		return n, nil
	}

	if s.pos == s.length {
		return 0, io.EOF
	}

	eof := false
	if s.pos+int64(len(buf)) > s.length {
		buf = buf[:s.length-s.pos]
		eof = true
	}

	var origBuf []byte
	if len(buf)%4 > 0 {
		origBuf = buf
		buf = make([]byte, len(origBuf)+4-len(origBuf)%4)
	}

	length := float64(soundserver.SampleRate) / float64(s.freq)

	num := (soundserver.BitDepthInBytes) * (numChannels)
	p := s.pos / int64(num)
	switch soundserver.BitDepthInBytes {
	case 1:
		for i := 0; i < len(buf)/num; i++ {
			const max = 127
			b := int(math.Sin(2*math.Pi*float64(p)/length) * 0.3 * max)
			for ch := 0; ch < numChannels; ch++ {
				buf[num*i+ch] = byte(b + 128)
			}
			p++
		}
	case 2:
		for i := 0; i < len(buf)/num; i++ {
			const max = 32767
			b := int16(math.Sin(2*math.Pi*float64(p)/length) * 0.3 * max)
			for ch := 0; ch < numChannels; ch++ {
				buf[num*i+2*ch] = byte(b)
				buf[num*i+1+2*ch] = byte(b >> 8)
			}
			p++
		}
	}

	s.pos += int64(len(buf))

	n := len(buf)
	if origBuf != nil {
		n = copy(origBuf, buf)
		s.remaining = buf[n:]
	}

	if eof {
		return n, io.EOF
	}
	return n, nil
}
