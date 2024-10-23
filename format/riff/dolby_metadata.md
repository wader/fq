Dolby Metadata from `<dbmd>` chunk of RIFF / WAV / Broadcast Wave Format (BWF),
including Dolby Atmos, AC3, Dolby Digital \[Plus\], and Dolby Audio Info (e.g. LUFS, True Peak).

### Examples
Decode Dolby metadata from `<dbmd>` chunk:
```
$ fq -d wav '.chunks[] | select(.id | IN("dbmd")) | tovalue' bwf.wav
```

RIFF / WAV / Broadcast Wave Format (BWF) chunks:
- `<chna>` Track UIDs of Audio Definition Model
- `<axml>` BWF XML Metadata, e.g. for Audio Definition Model ambisonics and elements

### Authors
- [@johnnymarnell](https://johnnymarnell.github.io), original author

### References
- https://adm.ebu.io/background/what_is_the_adm.html
- https://tech.ebu.ch/publications/tech3285s7
- https://tech.ebu.ch/publications/tech3285s5
- https://tech.ebu.ch/files/live/sites/tech/files/shared/tech/tech3285s6.pdf
- https://github.com/DolbyLaboratories/dbmd-atmos-parser
