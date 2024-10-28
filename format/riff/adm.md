[Audio Definition Model](https://adm.ebu.io/background/what_is_the_adm.html) including 3D Audio.

RIFF / WAV / Broadcast Wave Format (BWF) chunks:
- `<chna>` Chunk, Track UIDs of Audio Definition Model
- `<axml>` Chunk, BWF XML Metadata, e.g. for Audio Definition Model ambisonics and elements

### Examples
Decode ADM configuration from `<chna>` and `<axml>` chunks:
```bash
$ fq -d wav '.chunks[] | select(.id | IN("chna", "axml")) | tovalue' amd-bwf.wav

# Extract ADM <axml> chunk objects definitions xml content
$ fq -r -d wav '.chunks[] | select(.id | IN("axml")) | .xml | tovalue' amd-bwf.wav | tee axml-content.xml
```

### Authors
- [@johnnymarnell](https://johnnymarnell.github.io), original author

### References
- https://adm.ebu.io/background/what_is_the_adm.html
- https://tech.ebu.ch/publications/tech3285s7
- https://tech.ebu.ch/publications/tech3285s5