## fq

jq for binary files

```
# duration of a mp3 file
$ fq file.mp3 '[.frames[] | .samples_per_frame / .sample_rate] | add'
7504.169795907116

# width/height of embedded id3v2 jpeg picture
$ fq file.mp3 '.header.frames[] | select(.id == "APIC").picture.segments[] | select(.code._symbol == "SOF0")'
   |                                               |                |.header.frames[8].picture.segments[4]:
3b0|ff                                             |.               |  prefix: ff
3b0|   c0                                          | .              |  code: SOF0 (192)
3b0|      00 11                                    |  ..            |  Lf: 17
3b0|            08                                 |    .           |  P: 8
3b0|               01 40                           |     .@         |  Y: 320
3b0|                     01 40                     |       .@       |  X: 320
3b0|                           03                  |         .      |  Nf: 3
3b0|                              01 22 00 02 11 01|          ."....|  frame_components[3]:
3c0|03 11 01                                       |...             |
$  fq file.mp3 '.header.frames[] | select(.id == "APIC").picture.segments[] | select(.code._symbol == "SOF0") | {X,Y}'
{
  "X": 320,
  "Y": 320
}
$ fq file.mp3 '.header.frames[] | select(.id == "APIC").picture._raw' > picture.jpeg
$ file picture.jpeg
JPEG image data, JFIF standard 1.01, resolution (DPI), density 96x96, segment length 16, baseline, precision 8, 320x320, components 3

# bitrate of first two and last frames
$ fq file.mp3 '[(.frames[0:2] + .frames[-3:-1])[].bitrate]'
[
  128000,
  128000,
  128000,
  128000
]

# mp4 sidx
$ fq test.mp4 '.. | select(.type == "sidx") | [.index_table[] | {size, duration}]'
[
  {
    "duration": 94208,
    "size": 104035
  },
  {
    "duration": 33792,
    "size": 38475
  }
]
```

## Install

TODO

## Differences to jq / gojq

fq uses a fork of gojq that has these language additions

- Hex `0xab`, octal `0o77` and binary `0b101` integer literals
- Bitwise operations, `band`, `bor`, `bxor`, `bsl`, `bsr`, `bnot`
- `div` integer division function

fq also has some additions

- TODO: `scope` and `scopedump` functions used to implement REPL completion
- TODO: Custom object interface used to traverse fq's field tree and to allowing a terse
syntax for comparing and working with fields, accessing child fields and special properties like `_range`.
- `open(path) ` opens file
- `probe([name])` try to automatically detect and decode, TODO: rename to `decode`?
- All decoders are available as functions with their name, e.g. `... | mp3_frame`
- `d`/`dump` show field tree
- `v`/`verbose` show field tree verbosely
- `p`/`preview` show preview of field tree
- TODO: more functions

TODO: repl

## How to use

TODO: unknown for gaps
TODO: piping

## Decoders

[_dev/decoders_markdown]: sh-start

|Name|Description|
|-|-|
|aac_adts|Audio data transport stream packet|
|aac_frame|Advanced Audio Coding frame|
|aac_stream|Raw audio data transport stream|
|apev2|APEv2 metadata tag|
|bzip2|bzip2 compression|
|dns|DNS packet|
|elf|Executable and Linkable Format|
|flac|Free lossless audio codec|
|flac_frame|FLAC frame|
|flac_metadatablock|FLAC metadatablock|
|flac_picture|FLAC metadatablock picture|
|gzip|GZIP compression|
|icc_profile|International Color Consortium profile|
|id3_v1|ID3v1 metadata|
|id3_v11|ID3v1.1 metadata|
|id3_v2|ID3v2 metadata|
|jpeg|Joint Photographic Experts Group image|
|jq||
|mkv|Matroska|
|mp3|MPEG audio layer 3 file|
|mp3_frame|MPEG audio layer 3 frame|
|mp3_xing|Xing header|
|mp4|MP4 container|
|mpeg_asc|MPEG-4 Audio specific config|
|mpeg_avc|H.264/AVC sample|
|mpeg_avc_dcr|H.264/AVC Decoder configuration record|
|mpeg_es|MPEG elementary stream|
|mpeg_pes|MPEG Packetized elementary stream|
|mpeg_pes_packet|MPEG Packetized elementary stream packet|
|mpeg_spu|Sub picture unit (dvd subtitle)|
|ogg|OGG container|
|ogg_page|OGG page|
|opus_packet|Opus packet|
|png|Portable network graphics image|
|raw|Raw bits|
|tar|Tar archive|
|tiff|Tag Image File Format|
|vorbis_comment|Vorbis comment|
|vorbis_packet|Vorbis packet|
|vp9_frame|VP9 frame|
|wav|WAV container|

[#]: sh-end

TODO: format graph?

## Own decoders and use as library

TODO

## Known issues, TODOs and ideas

### Known issues

- TODO: concat bitbufs?
- TODO: byte units when outputting

### TODOs

- Nested BitBufs, how to show? what about ranges? for example compressed data, demuxed ogg
- Clean up panics, errors, better partial decode
- bitio.MultiBitReader to save memory
- Cleanup decoder API
- Save encoding for values, LE, BE, varint etc
- Cleanup decoders
- Document decode maturity/completeness
- Embed jq code using go 1.16 embed
- Arbitrary integer base literals

### Ideas

- Some kind of UI, web and cli? maybe would be nice to show hex dump etc with overlapping fields?
- Would it be possible to save memory by just record range/decoder at first decode and
then decode as needed later?
- Move more things to jq code, dumper, CLI, help, REPL?
- Some kind of bit vs bytes position notation/type

## Thanks

Would not be possible without [itchyny](https://github.com/itchyny)'s
jq implementation [gojq](https://github.com/itchyny/gojq).
