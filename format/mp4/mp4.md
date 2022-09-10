#### `mp4_path($path)` - Lookup mp4 box using a mp4 box path.

```sh
# <decode value box> | mp4_path($path) -> <decode value box>
$ fq 'mp4_path(".moov.trak[1]")' file.mp4
```

#### `mp4_path` - Return mp4 box path for a decode value box.

```sh
# <decode value box> | mp4_path -> string
$ fq 'grep_by(.type == "trak") | mp4_path' file.mp4
```

#### Force decode a single box

```sh
$ fq -n '"AAAAHGVsc3QAAAAAAAAAAQAAADIAAAQAAAEAAA==" | frombase64 | mp4({force:true}) | d'
```

#### Don't decode samples and manually decode first sample for first track as a `aac_frame`

```sh
$ fq -o decode_samples=false '.tracks[0].samples[0] | aac_frame | d' file.mp4
```

#### Entries for first edit list as values

```sh
$ fq 'first(grep_by(.type=="elst").entries) | tovalue' file.mp4
```

#### References

- [ISO/IEC base media file format (MPEG-4 Part 12)](https://en.wikipedia.org/wiki/ISO/IEC_base_media_file_format)
- [Quicktime file format](https://developer.apple.com/standards/qtff-2001.pdf)
