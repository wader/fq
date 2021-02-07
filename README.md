# fq

jq for binary files

Tool and framework for querying and exploring binary formats.

##

``` (exec)
# duration of a mp3 file
$ fq file.mp3 '[.frames[] | .samples_per_frame / .sample_rate] | add'
0.0783673469387755
 
# embedded id3v2 png picture
$ fq file.mp3 '.headers[].frames[] | select(.id == "APIC").picture'
   |                                               |                |.headers[0].frames[1].picture: png
030|                              89 50 4e 47 0d 0a|          .PNG..|  signature: Correct ("\x89PNG\r\n\x1a\n")
040|1a 0a                                          |..              |
040|      00 00 00 0d 49 48 44 52 00 00 01 40 00 00|  ....IHDR...@..|  chunks[4]:
050|00 f0 08 02 00 00 00 fe 4f 2a 3c 00 00 00 09 70|........O*<....p|
*  |2649 bytes more until ab8+7                    |                |
 
# resolution of embedded png picture
$ fq file.mp3 '.headers[].frames[] | select(.id == "APIC").picture.chunks[] | select(.type == "IHDR") | {width, height}'
{
  "height": 240,
  "width": 320
}
 
$ fq file.mp3 '.headers[].frames[] | select(.id == "APIC").picture._raw' > file.png
$ file file.png
file.png: PNG image data, 320 x 240, 8-bit/color RGB, non-interlaced
 
# codecs in mp4 file
$ fq file.mp4 '[.. | select(.type == "stsd").sample_descriptions[].data_format]'
[
  "avc1",
  "mp4a"
]
```

## Install

TODO

## Differences to jq / gojq

fq uses a fork of [gojq](https://github.com/itchyny/gojq) that has these additions

Language:

- Hex `0xab`, octal `0o77` and binary `0b101` integer literals
- Bitwise operations, `band`, `bor`, `bxor`, `bsl`, `bsr`, `bnot`
- `div` integer division function

Functions:

- `open(path) ` opens file
- `decode([name])` decode as named format or try to automatically detect
- All decoders an groups are available as decode functions with their name, e.g. `... | mp3_frame`, `image`, `probe`.
- `d`/`dump` show field tree
- `v`/`verbose` show field tree verbosely
- `p`/`preview` show preview of field tree
- TODO: more functions
- TODO: `scope` and `scopedump` functions used to implement REPL completion
- TODO: Custom object interface used to traverse fq's field tree and to allowing a terse
syntax for comparing and working with fields, accessing child fields and special properties like `_range`.

TODO: repl

## How to use

TODO: unknown for gaps

TODO: piping

## Decoders

[./decoders_markdown]: sh-start

|Name|Description|
|-|-|
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
|mpeg_aac_frame|Advanced Audio Coding frame|
|mpeg_aac_stream|Raw audio data transport stream|
|mpeg_adts|Audio data transport stream packet|
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
- CRC fields, how to update with actual? fix flac
- Clean up panics, errors, better partial decode
- bitio.MultiBitReader to save memory
- Cleanup decoder API, nested bufs, decoders, try decode loop? decodebuf?
- Save encoding for values, LE, BE, varint etc
- Cleanup decoders
- Document decode maturity/completeness
- Embed jq code using go 1.16 embed
- Arbitrary integer base literals
- Make jq functions that change state fail if called more than once? decode etc?
- REPL push/pop, variables etc?
- REPL tests

### Ideas

- Suppose writing decoder in scripting language, jq, js, tango etc?
- Some kind of UI, web and cli? would be nice to visualize overlapping fields
- Is it possible to save memory by just record range/decoder at first decode and
then decode as needed later?
- Move more things to jq code, dumper, CLI, help, REPL?
- Some kind of bit vs bytes position notation/type

## Thanks

fq would not be possible without [itchyny](https://github.com/itchyny)'s
jq implementation [gojq](https://github.com/itchyny/gojq).
