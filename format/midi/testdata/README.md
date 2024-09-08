# NOTES

## MIDI files

The test and example MIDI files are located in the _testdata/midi_ folder.

1. _reference.mid_
Reference MIDI file for testing/development only (it is not a valid MIDI file), with two tracks:
- _Track 0_: _empty_ track with only an _end-of-track_ event
- _Track 1_: _instrument_ track with sequential examples of all decoded MIDI events

2. _format-0.mid_
Basic MIDI format 0 test file. Contains a single track with only a _track name_ and _end-of-track_ events.

3. _format-1.mid_
Basic MIDI format 1 test file. Contains two tracks, each with only a _track name_ and _end-of-track_ events.

4. _format-2.mid_
Basic MIDI format 2 test file.  Contains two tracks, each with only a _track name_ and _end-of-track_ events.

5. _smpte-timecode.mid_
MIDI format 0 test file with an SMPTE timecode for the divisions field.

6. _empty.mid_
Empty MIDI file to verify MIDI decoder handles empty files.

7. _key_signatures.mid_

Test file with all supported MIDI key signatures.

8. _notes.mid_

Test file with all supported MIDI notes.

9. _unknown-chunks.mid_

Test file with 'alien' chunks interleaved with the _MTrk_ track chunks.

10. _invalid-MThd-length.mid_

Test file with invalid _MThd_ chunk length.

11. _invalid-MTrk-length.mid_

Test file with invalid _MTrk_ chunk length.

12. _twinkle.mid_

Sample valid MIDI file for the example queries in the help.


## MIDI event files

MIDI files with a single event for development and debugging are located in the _testdata/events_ folder.

### Meta events

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
manufacturer: 00 00 3b (Mark Of The Unicorn (MOTU))
data:         3a 4c 5e
```

### MIDI events

1. _note-off.mid_
```
00 81 70 60      40 71 48

delta: 0         delta: 64
channel: 1
note: 112 (E8)   note: 113 (F8)
velocity: 96     velocity: 72
```

2. _note-on.mid_
```
00 92 30 48      40 32 48

delta: 0         delta: 64
channel: 0
note: 48 (C3)    note: 50 (D3)
velocity: 72     velocity: 72
```

3. _polyphonic-pressure.mid_
```
00 a0 64         40 48

delta: 0         delta: 64
channel: 0
pressure: 100    pressure: 72
```

4. _controller.mid_
```
00 b0 20 21      40 20 22

delta: 0         delta: 64
channel: 0
controller: 32   controller: 32
value: 33        value: 34
```

5. _program-change.mid_
```
00 c0 19         40 20

delta: 0         delta: 64
channel: 0
program: 25      program: 32
```

6. _channel-pressure.mid_
```
00 d0 07         40 48

delta: 0         delta: 64
channel: 0
pressure: 7      pressure: 72
```

7. _pitch-bend.mid_
```
00 e5 40 00    20 60 00       40 20 00
             
delta: 0       delta: 32      delta: 64
channel: 5
bend: 0        bend: 4096     bend: -4096
```

### System Exclusive events

1. _sysex-message.mid_
```
00 f0 05 7e 00 09 01 f7

delta: 0
manufacturer: 7e (Non-RealTime Extensions)
data: 00 09 01
```

2. _sysex-continuation.mid_
```
00 f0 03 43 01 23           00 f7 06 45 67 89 ab cd ef    00 f7 04 01 23 45 f7

delta: 0                    delta: 0                      delta: 0
manufacturer: 43 (Yamaha)   data: 45 67 89 ab cd ef       data: 01 23 45
data: 01 23
```

3. _sysex-escape.mid_
```
00 f7 02 f3 01

delta: 0
data: f3 01
```

