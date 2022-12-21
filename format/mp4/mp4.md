### Speed up decoding by not decoding samples

```sh
# manually decode first sample as a aac_frame
$ fq -o decode_samples=false '.tracks[0].samples[0] | aac_frame | d' file.mp4
```

### Entries for first edit list as values

```sh
$ fq 'first(grep_by(.type=="elst").entries) | tovalue' file.mp4
```

### Whole box tree as JSON (exclude mdat data and tracks)

```sh
$ fq 'del(.tracks) | grep_by(.type=="mdat").data = "<excluded>" | tovalue' file.mp4
```

### Force decode a single box

```sh
$ fq -n '"AAAAHGVsc3QAAAAAAAAAAQAAADIAAAQAAAEAAA==" | from_base64 | mp4({force:true}) | d'
```

### Lookup mp4 box using a mp4 box path.

```sh
# <decode value box> | mp4_path($path) -> <decode value box>
$ fq 'mp4_path(".moov.trak[1]")' file.mp4
```

### Get mp4 box path for a decode value box.

```sh
# <decode value box> | mp4_path -> string
$ fq 'grep_by(.type == "trak") | mp4_path' file.mp4
```

### References

- [ISO/IEC base media file format (MPEG-4 Part 12)](https://en.wikipedia.org/wiki/ISO/IEC_base_media_file_format)
- [Quicktime file format](https://developer.apple.com/standards/qtff-2001.pdf)
