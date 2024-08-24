### Notes

1. Only supports the MIDI 1.0 specification.
2. Does only basic validation on the MIDI data.

### Sample queries

1. Extract the track names from a MIDI file
```
fq -d midi -d midi '.. | select(.event=="TrackName")? | "\(.name)"' twinkle.mid 
```

2. Extract the tempo changes from a MIDI file
```
fq -d midi '.. | select(.event=="Tempo")?.tempo' twinkle.mid
```

3. Extract the key changes from a MIDI file
```
fq -d midi '.. | select(.event=="KeySignature")?.key' key-signatures.mid
```

4. Extract NoteOn and NoteOff events:
```
fq -d midi 'grep_by(.event=="NoteOn" or .event=="NoteOff") | "\(.event)  \(.time.tick)  \(.note)"' twinkle.mid
```

### Authors
- transcriptaze.development@gmail.com

### References

1. [The Complete MIDI 1.0 Detailed Specification](https://www.midi.org/specifications/item/the-midi-1-0-specification)
