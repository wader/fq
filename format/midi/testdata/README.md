# NOTES

## Files

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

