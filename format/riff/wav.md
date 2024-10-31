WAVE audio file format.

Also includes support for [Audio Definition Model](https://adm.ebu.io/background/what_is_the_adm.html) and 3D Audio.

RIFF / WAV / Broadcast Wave Format (BWF) chunks:

- `RIFF`: primary container chunk specifying the file type and containing sub-chunks (e.g., fmt, data)
- `fmt`: describes format / stream encoding in data chunk
- `data`: indicates size and contains encoded raw sound data
- `bext`: broadcast extension chunk, containing broadcast-specific metadata such as description, originator, creation date, time reference, and more
- `LIST`: organizes additional metadata in sub-chunks, often used to include information like artist, genre, or title in INFO or other standardized formats
- `smpl`: sample metadata chunk, containing looping and sampling information, such as start and end points for loops, sample rate, and MIDI pitch
- `fact`: contains metadata on the original uncompressed data, such as the number of samples, typically used in non-PCM (compressed) formats to aid in playback and synchronization
- `chna`: track UIDs of Audio Definition Model
- `axml`: XML metadata, e.g. for Audio Definition Model ambisonics and elements as in [EBUCore spec](https://tech.ebu.ch/docs/tech/tech3293.pdf)
- `dbmd`: Dolby specific metadata like loudness and binaural settings, see also [`dolby_metadata` format](#dolby_metadata)


### Examples
Decode ADM configuration from `<chna>` and `<axml>` chunks:
```bash
$ fq -d wav '.chunks[] | select(.id | IN("chna", "axml")) | tovalue' amd-bwf.wav

# Extract ADM <axml> chunk objects definitions xml content
$ fq -r -d wav '.chunks[] | select(.id | IN("axml")) | .xml | tovalue' amd-bwf.wav | tee axml-content.xml
```

### Authors
- [@wader](https://github.com/wader), original author
- [@johnnymarnell](https://johnnymarnell.github.io), ADM and Dolby support

### References
- http://soundfile.sapp.org/doc/WaveFormat/
- https://github.com/FFmpeg/FFmpeg/blob/master/libavformat/wavdec.c
- https://tech.ebu.ch/docs/tech/tech3285.pdf
- http://www-mmsp.ece.mcgill.ca/Documents/AudioFormats/WAVE/WAVE.html
- https://adm.ebu.io/background/what_is_the_adm.html
- https://tech.ebu.ch/docs/tech/tech3285s7.pdf
- https://tech.ebu.ch/docs/tech/tech3285s5.pdf
