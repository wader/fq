/*
Package midi implements an fq plugin to decode [standard MIDI files].

The MIDI decoder is a member of the 'probe' group and fq should automatically invoke the 
decoder when opening a MIDI file. The decoder can be explicitly specified with the '-d midi'
command line option.

The decoder currently only supports MIDI 1.0 files and does only basic validation on the 
MIDI file structure.

[standard MIDI files]: https://midi.org/standard-midi-files.
*/
package midi
