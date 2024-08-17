# TODO

- [x] update forked master branch
- [ ] discard unknown chunks
- [ ] tests
      - [ ] format 0
      - [ ] format 1
      - [ ] format 2
- [x] move delta into events
- [ ] fix gaps
- (?) warn only for e.g. invalid format tracks
- (?) add to probe group

- formats
    - [ ] format 0
    - [ ] format 1
    - [ ] format 2

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
          - [ ] check key mappings
    - [x] time signature
    - [ ] sequencer specific event
          - [ ] map manufacturer

- midi events
    - [x] note off
          - [ ] map note name
    - [x] note on
          - [ ] map note name
    - [x] polyphonic pressure
    - [x] controller
    - [x] program change
    - [x] channel pressure
    - [x] pitch bend
    - [ ] running status

- sysex events
    - [ ] message
    - [ ] continuation
    - [ ] escape

