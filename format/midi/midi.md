### Notes

1. Only supports the MIDI 1.0 MIDI file specification.
2. Only supports _MThd_ and _MTrk_ chunks.
3. Does only basic validation on the MIDI data.

### Sample queries

1. Extract the track names from a MIDI file
```
fq -d midi -d midi '.. | select(.event=="track_name")? | "\(.track_name)"' midi/twinkle.mid 
```

2. Extract the tempo changes from a MIDI file
```
fq -d midi '.. | select(.event=="tempo")?.tempo' midi/twinkle.mid
```

3. Extract the key changes from a MIDI file
```
fq -d midi '.. | select(.event=="key_signature")?.key_signature' midi/twinkle.mid
```

4. Extract NoteOn events:
```
fq -d midi 'grep_by(.event=="note_on") | [.time.tick, .note_on.note] | join(" ")' midi/twinkle.mid
```

### Authors
- [transcriptaze](https://github.com/transcriptaze)

### References

1. [The Complete MIDI 1.0 Detailed Specification](https://www.midi.org/specifications/item/the-midi-1-0-specification)
2. [Standard MIDI Files](https://midi.org/standard-midi-files)
3. [Standard MIDI File (SMF) Format](http://midi.teragonaudio.com/tech/midifile.htm)
4. [MIDI Files Specification](http://www.somascape.org/midi/tech/mfile.html)
5. [MIDI SMPTE Offset meta message](https://www.recordingblogs.com/wiki/midi-smpte-offset-meta-message)
6. [Somascape MIDI Files Specification](http://www.somascape.org/midi/tech/mfile.html#meta)
