### Notes

1. Only supports the MIDI 1.0 specification.
2. Does only basic validation on the MIDI data.

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
- transcriptaze.development@gmail.com

### References

1. [The Complete MIDI 1.0 Detailed Specification](https://www.midi.org/specifications/item/the-midi-1-0-specification)
2. [The MIDI Specification](http://midi.teragonaudio.com/tech/midispec.htm)
