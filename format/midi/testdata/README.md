# NOTES

## MIDI files

The test and example MIDI files are located in the _testdata/midi_ folder.

1. _format-0.mid_
MIDI format 0 reference file. Contains a single track with all supported MIDI events.

2. _format-1.mid_
MIDI format 1 reference file. Contains two tracks:
- _Track 0_, a tempo track with the _Time Signature_ and _Tempo_ events
- _Track 1_, with all the other supported MIDI events

3. _format-2.mid_
MIDI format 2 reference file. Contains two tracks:
- _Track 0_, a track with all supported MIDI events
- _Track 1_, a reversed version of _Track 0_

4. _empty.mid_
Empty MIDI file to verify MIDI decoder handles empty files.

5. _unknown_chunks.mid_

MIDI file with non-MIDI chunks interleaved between the _MTrk_ track chunks.

6. _key_signatures.mid_

Test file with all supported MIDI key signatures.

7. _notes.mid_

Test file with all supported MIDI notes.

8. _twinkle.mid_

Sample MIDI file for the example queries in the help.


## MIDI event files

MIDI files with a single event for development and debugging are located in the _testdata/events_ folder.

### Metaevents

1. _sequence-number.mid_
```
00 ff 00 02 00 17

delta: 0
sequence number: 23
```

2. _text.mid_
```
00 ff 01 0d 54 68 69 73 20 61 6e 64 20 54 68 61 74

delta: 0
text: This and That
```

3. _copyright.mid_
```
00 ff 02 04 54 68 65 6d

delta: 0
copyright: Them
```

4. _track_name.mid_
```
00 ff 03 0f 41 63 6f 75 73 74 69 63 20 47 75 69 74 61 72

delta: 0
track name: Acoustic Guitar

```

5. _instrument_name.mid_
```
00 ff 04 0a 44 69 64 67 65 72 69 64 6f 6f

delta: 0
instrument: Didgeridoo
```

6. _lyric.mid_
```
00 ff 05 08 4c 61 2d 6c 61 2d 6c 61

delta: 0
lyric: La-la-la
```

7. _marker.mid_
```
00 ff 06 0f 48 65 72 65 20 42 65 20 44 72 61 67 6f 6e 73

delta: 0
marker: Here Be Dragons
```

8. _cuepoint.mid_
```
00 ff 07 0c 4d 6f 72 65 20 63 6f 77 62 65 6c 6c

delta: 0
cue: More cowbell
```

9. _program_name.mid_
```
00 ff 08 06 45 73 63 61 70 65

delta: 0
program: Escape
```

10. _device_name.mid_
```
00 ff 09 08 54 68 65 54 68 69 6e 67

delta: 00
device: TheThing
```

11. _midi-channel-prefix.mid_
```
00 ff 20 01 0d

delta: 00
MIDI channel prefix: 13
```

12. _midi-port.mid_
```
00 ff 21 01 70

delta: 00
MIDI port: 112
```

13. _tempo.mid_
```
00 ff 51 03 07 a1 20

delta: 0
tempo: 500000
```

14. _smpte-offset.mid_
```
00 ff 54 05 4d 2d 3b 07 27

delta: 0
framerate: 25
hour:      13
minute:    45
second:    59
frames:     7
fractions: 39
```

15. _time-signature.mid_
```
00 ff 58 04 04 02 18 08 

delta: 0
numerator:   4 
denominator: 4
ticks_per_click: 24
thirty_seconds_per_quarter: 8 
```

16. _key-signature.mid_
```
00 ff 59 02 00 01 

delta: 0
key: A minor
```

17. _end-of-track.mid_
```
00 ff 2f 00

delta: 0
```

18. _sequencer-specific-event_
```
00 ff 7f 06 00 00 3b 3a 4c 5e

delta: 0
manufacturer: 00 00 3b Mark Of The Unicorn (MOTU)
data:         3a 4c 5e
```

### MIDI events

1. _note-off.mid_
```
00 81 70 60

delta: 0
channel: 1
note: 112 (E3)
velocity: 96
```

2. _note-on.mid_
```
00 90 30 48

delta: 0
channel: 0
note: 48 (C3)
velocity: 72
```

3. _polyphonic-pressure.mid_
```
00 a0 64

delta: 0
channel: 0
pressure: 100
```

4. _controller.mid_
```
00 b0 20 21

delta: 0
channel: 0
controller: 32
value: 33
```

5. _program-change.mid_
```
00 c0 19

delta: 0
channel: 0
program: 25
```

6. _channel-pressure.mid_
```
00 d0 07

delta: 0
channel: 0
pressure: 7
```

7. _pitch-bend.mid_
```
00 e0 40 00    00 e0 60 00    00 e0 20 00
             
delta: 0       delta: 0       delta: 0
channel: 0     channel: 0     channel: 0
bend: 0        bend: 4096     bend: -4096
```
