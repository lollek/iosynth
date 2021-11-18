## iosynth

iosynth is an simple synthesizer, heavily inspired by
[hundredrabbits/Pilot](https://github.com/hundredrabbits/Pilot) and is created
to be controlled externally by [hundredrabbits/orca](https://git.sr.ht/~rabbits/orca).

By default, it starts a listening port to UDP 49161, and only has a single
channel (Channel 0) with a sine wave sound. You can send it a string of three characters in
order to make a sound, in the format CHANNEL-OCTAVE-NOTE. E.g. 04C will play
"middle C" and 04c will play "middle c#"
