## iosynth

iosynth is an simple synthesizer, heavily inspired by
[hundredrabbits/Pilot](https://github.com/hundredrabbits/Pilot) and is created
to be controlled externally by [hundredrabbits/orca](https://git.sr.ht/~rabbits/orca).

By default, it starts a listening port to UDP 49161. You can send it a string of three characters in
order to make a sound, in the format TRACK-OCTAVE-NOTE. See examples below

Example:
* 04C = Middle C as a sine wave
* 04c = Middle C# as a sine wave
* 01C = Middle C as a square wave

Tracks:
* 0: Sine wave
* 1: Square wave
* 2: Triangle wave
* 3: Sawtooth wave

Note that it's still very much work in progress.
