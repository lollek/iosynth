package mixer

import (
	"time"
	"strings"

	"github.com/lollek/iosynth/soundserver"
	"github.com/lollek/iosynth/sound"
)

func RawCommand(data []byte) {
	commands := string(data)
	for _, command := range strings.Split(commands, ";") {
		parseCommand(command)
	}

}

func parseCommand(command string) {
	// command = [CHANNEL] [OCTAVE] [NOTE] [Optional] [Optional]
	if len(command) < 3 || len(command) > 5 {
		return
	}

	channel := command[0]
	if channel != '0' {
		return
	}

	octave := command[1] - '0'
	note := command[2]
	for ('H' <= note && note <= 'Z') || ('h' <= note && note <= 'z') {
		note -= 8
		octave += 1
	}
	if note == 'b' {
		note = 'C'
		octave += 1
	}

	if !(('A' <= note && note <= 'G') || ('a' <= note && note <= 'g')) {
		return
	}

	if octave < 0 {
		octave = 0
	} else if octave > 8 {
		octave = 8
	}

	note2index := map[byte]int {
		'C': 0,
		'c': 1,
		'D': 2,
		'd': 3,
		'E': 4,
		'e': 5,
		'F': 5,
		'f': 6,
		'G': 7,
		'g': 8,
		'A': 9,
		'a': 10,
		'B': 11,
	}


	octave2list := map[byte][12]float64{
		0: [12]float64{16.35, 17.32, 18.35, 19.45, 20.60, 21.83, 23.12, 24.50, 25.96, 27.50, 29.14, 30.87 },
		1: [12]float64{32.70, 34.65, 36.71, 38.89, 41.20, 43.65, 46.25, 49.00, 51.91, 55.00, 58.27, 61.74,},
		2: [12]float64{65.41, 69.30, 73.42, 77.78, 82.41, 87.31, 92.50, 98.00, 103.83, 110.00, 116.54, 123.47,},
		3: [12]float64{130.81, 138.59, 146.83, 155.56, 164.81, 174.61, 185.00, 196.00, 207.65, 220.00, 233.08, 246.94,},
		4: [12]float64{261.63, 277.18, 293.66, 311.13, 329.63, 349.23, 369.99, 392.00, 415.30, 440.00, 466.16, 493.88,},
		5: [12]float64{523.25, 554.37, 587.33, 622.25, 659.25, 698.46, 739.99, 783.99, 830.61, 880.00, 932.33, 987.77,},
		6: [12]float64{1046.50, 1108.73, 1174.66, 1244.51, 1318.51, 1396.91, 1479.98, 1567.98, 1661.22, 1760.00, 1864.66, 1975.53,},
		7: [12]float64{2093.00, 2217.46, 2349.32, 2489.02, 2637.02, 2793.83, 2959.96, 3135.96, 3322.44, 3520.00, 3729.31, 3951.07,},
		8: [12]float64{4186.01, 4434.92, 4698.63, 4978.03, 5274.04, 5587.65, 5919.91, 6271.93, 6644.88, 7040.00, 7458.62, 7902.13,},
	}

	freq := octave2list[octave][note2index[note]]
	duration := 250 * time.Millisecond
	sound := sound.NewSineWave(freq, duration)
	soundserver.PlaySound(sound)
}
