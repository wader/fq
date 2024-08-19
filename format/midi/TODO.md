# TODO

- [x] update forked master branch
- [x] move delta into events
- [x] Use FieldUTF8 for MIDI chunk tags
- [ ] discard unknown chunks
- [ ] assert available bytes
- [ ] tests
      - [ ] format 0
      - [ ] format 1
      - [ ] format 2
- (?) example queries
      - tempo changes
      - key changes
      - notes
- (?) add to probe group
- (?) tick field

- [ ] fix gaps
      - [x] SequencerSpecificEvent
      - [x] SMPTEOffset
      - [x] TimeSignature
      - [x] SysExMessage
      - [ ] SysEx - 'continued' flag

- meta events
    - [x] sequence number
    - [x] text
    - [x] copyright
    - [x] track name 
    - [x] instrument name
    - [x] lyric
    - [x] marker
    - [x] cue point
    - [x] program name
    - [x] device name
    - [x] MIDI channel prefix
    - [x] MIDI port
    - [x] end of track
    - [x] tempo
    - [x] SMPTE offset
    - [x] key signature
    - [x] time signature
    - [x] sequencer specific event
    - [x] map manufacturer
    - [x] check key mappings
    - [ ] Combine status + event into metaevent field

- midi events
    - [x] note off
    - [x] note on
    - [x] polyphonic pressure
    - [x] controller
    - [x] program change
    - [x] channel pressure
    - [x] pitch bend
    - [x] running status
    - [x] use context struct for running status
    - [x] map note names

- sysex events
    - [x] message
    - [x] continuation
    - [x] escape
    - [x] use context struct for casio
    - [x] map manufacturer ID


