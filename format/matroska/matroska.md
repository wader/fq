### Lookup element using path

```sh
$ fq 'matroska_path(".Segment.Tracks[0)")' file.mkv
```

### Get path to element

```sh
$ fq 'grep_by(.id == "Tracks") | matroska_path' file.mkv
```

### References
- https://tools.ietf.org/html/draft-ietf-cellar-ebml-00
- https://matroska.org/technical/specs/index.html
- https://www.matroska.org/technical/basics.html
- https://www.matroska.org/technical/codec_specs.html
- https://wiki.xiph.org/MatroskaOpus
