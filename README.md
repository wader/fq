# fq

jq for binaries

Tool and framework for querying and exploring binary data.

##

<sub>
<pre sh>
<b># duration of a mp3 file</b> 
$ fq file.mp3 '[.frames[] | .samples_per_frame / .sample_rate] | add' 
0.0783673469387755
 
<b># embedded id3v2 png picture</b> 
$ fq file.mp3 '.headers[].frames[] | select(.id == "APIC")' 
     |00 01 02 03 04 05 06 07 08 09 0a 0b 0c 0d 0e 0f|                |.headers[0].frames[1]: {}
0x020|         41 50 49 43                           |   APIC         |  id: "APIC" (Attached picture)
0x020|                     00 00 15 0c               |       ....     |  size: 2700
0x020|                                 00 00         |           ..   |  flags: {}
0x020|                                       03      |             .  |  text_encoding: UTF-8 (3)
0x020|                                          69 6d|              im|  mime_type: "image/png"
0x030|61 67 65 2f 70 6e 67 00                        |age/png.        |
0x030|                        00                     |        .       |  picture_type: 0
0x030|                           00                  |         .      |  description: ""
0x030|                              89 50 4e 47 0d 0a|          .PNG..|  picture: png
0x040|1a 0a 00 00 00 0d 49 48 44 52 00 00 01 40 00 00|......IHDR...@..|
*    |2665 bytes more until 0xab8.7                  |                |
 
<b># resolution of embedded png picture</b> 
$ fq file.mp3 '.headers[].frames[] | select(.id == "APIC").picture.chunks[] | select(.type == "IHDR") | {width, height}' 
{
  "height": 240,
  "width": 320
}
 
<b># extract png</b> 
$ fq file.mp3 '.headers[].frames[] | select(.id == "APIC").picture._bits' > file.png 
$ file file.png 
file.png: PNG image data, 320 x 240, 8-bit/color RGB, non-interlaced
 
<b># codecs in a mp4 file</b> 
$ fq file.mp4 '[.. | select(.type == "stsd").sample_descriptions[].data_format]' 
[
  "avc1",
  "mp4a"
]
</pre>
</sub>

## Install

TODO

## Differences to jq

fq uses a fork of [gojq](https://github.com/itchyny/gojq) that has these additions:

Language:

- Hex `0xab`, octal `0o77` and binary `0b101` integer literals
- Bitwise operations, `band`, `bor`, `bxor`, `bsl`, `bsr`, `bnot`
- `div` integer division operator

Functions:

- `open(path)` open file
- `decode([name])` decode as named format or try to automatically detect
- All decoders and groups are available as functions with their name, e.g. `... | mp3_frame`, `image`, `probe`.
- `repl` nested REPL
- `d`/`display` show field tree
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

[./formats_markdown.jq]: sh-start

|Name                 |Description                                                  |Uses|
|-                    |-                                                            |-|
|`aac_frame`          |Advanced&nbsp;Audio&nbsp;Coding&nbsp;frame                   |<sub></sub>|
|`aac_stream`         |Raw&nbsp;audio&nbsp;data&nbsp;transport&nbsp;stream          |<sub>`adts`</sub>|
|`adts`               |Audio&nbsp;data&nbsp;transport&nbsp;stream&nbsp;packet       |<sub>`aac_frame`</sub>|
|`apev2`              |APEv2&nbsp;metadata&nbsp;tag                                 |<sub></sub>|
|`av1_ccr`            |AV1&nbsp;codec&nbsp;configuration&nbsp;record                |<sub></sub>|
|`av1_frame`          |AV1&nbsp;frame                                               |<sub>`av1_obu`</sub>|
|`av1_obu`            |AV1&nbsp;open&nbsp;bitstream&nbsp;unit                       |<sub></sub>|
|`avc_dcr`            |H.264/AVC&nbsp;Decoder&nbsp;configuration&nbsp;record        |<sub></sub>|
|`avc_nal`            |H.264/AVC&nbsp;sample                                        |<sub></sub>|
|`bzip2`              |bzip2&nbsp;compression                                       |<sub>`probe`</sub>|
|`dns`                |DNS&nbsp;packet                                              |<sub></sub>|
|`elf`                |Executable&nbsp;and&nbsp;Linkable&nbsp;Format                |<sub></sub>|
|`exif`               |Exchangeable&nbsp;Image&nbsp;File&nbsp;Format                |<sub>`icc_profile`</sub>|
|`flac`               |Free&nbsp;lossless&nbsp;audio&nbsp;codec&nbsp;file           |<sub>`flac_metadatablock`, `flac_frame`</sub>|
|`flac_frame`         |FLAC&nbsp;frame                                              |<sub></sub>|
|`flac_metadatablock` |FLAC&nbsp;metadatablock                                      |<sub>`flac_picture`, `vorbis_comment`</sub>|
|`flac_picture`       |FLAC&nbsp;metadatablock&nbsp;picture                         |<sub>`image`</sub>|
|`gif`                |Graphics&nbsp;Interchange&nbsp;Format                        |<sub></sub>|
|`gzip`               |gzip&nbsp;compression                                        |<sub>`probe`</sub>|
|`hevc_dcr`           |H.265/HEVC&nbsp;Decoder&nbsp;configuration&nbsp;record       |<sub></sub>|
|`hevc_nal`           |H.265/HEVC&nbsp;sample                                       |<sub></sub>|
|`icc_profile`        |International&nbsp;Color&nbsp;Consortium&nbsp;profile        |<sub></sub>|
|`id3_v1`             |ID3v1&nbsp;metadata                                          |<sub></sub>|
|`id3_v11`            |ID3v1.1&nbsp;metadata                                        |<sub></sub>|
|`id3_v2`             |ID3v2&nbsp;metadata                                          |<sub>`image`</sub>|
|`jpeg`               |Joint&nbsp;Photographic&nbsp;Experts&nbsp;Group&nbsp;file    |<sub>`exif`</sub>|
|`mkv`                |Matroska&nbsp;file                                           |<sub>`av1_ccr`, `av1_frame`, `flac_frame`, `flac_metadatablock`, `mp3_frame`, `aac_frame`, `mpeg_asc`, `avc_dcr`, `avc_nal`, `hevc_dcr`, `hevc_nal`, `mpeg_spu`, `opus_packet`, `vorbis_packet`, `vp8_frame`, `vp9_frame`</sub>|
|`mp3`                |MP3&nbsp;file                                                |<sub>`id3_v2`, `id3_v1`, `id3_v11`, `apev2`, `mp3_frame`</sub>|
|`mp3_frame`          |MPEG&nbsp;audio&nbsp;layer&nbsp;3&nbsp;frame                 |<sub>`mp3_xing`</sub>|
|`mp3_xing`           |Xing&nbsp;header                                             |<sub></sub>|
|`mp4`                |MPEG-4&nbsp;file                                             |<sub>`av1_ccr`, `av1_frame`, `flac_frame`, `flac_metadatablock`, `mp3_frame`, `aac_frame`, `avc_dcr`, `avc_nal`, `mpeg_es`, `hevc_dcr`, `hevc_nal`, `opus_packet`, `vorbis_packet`, `vp9_frame`</sub>|
|`mpeg_asc`           |MPEG-4&nbsp;Audio&nbsp;specific&nbsp;config                  |<sub></sub>|
|`mpeg_es`            |MPEG&nbsp;elementary&nbsp;stream                             |<sub>`mpeg_asc`, `vorbis_packet`</sub>|
|`mpeg_pes`           |MPEG&nbsp;Packetized&nbsp;elementary&nbsp;stream             |<sub>`mpeg_pes_packet`, `mpeg_spu`</sub>|
|`mpeg_pes_packet`    |MPEG&nbsp;Packetized&nbsp;elementary&nbsp;stream&nbsp;packet |<sub></sub>|
|`mpeg_spu`           |Sub&nbsp;picture&nbsp;unit&nbsp;(dvd&nbsp;subtitle)          |<sub></sub>|
|`ogg`                |OGG&nbsp;file                                                |<sub>`ogg_page`, `vorbis_packet`, `opus_packet`</sub>|
|`ogg_page`           |OGG&nbsp;page                                                |<sub></sub>|
|`opus_packet`        |Opus&nbsp;packet                                             |<sub>`vorbis_comment`</sub>|
|`png`                |Portable&nbsp;network&nbsp;graphics&nbsp;file                |<sub>`icc_profile`, `exif`</sub>|
|`raw`                |Raw&nbsp;bits                                                |<sub></sub>|
|`tar`                |Tar&nbsp;archive                                             |<sub>`probe`</sub>|
|`tiff`               |Tag&nbsp;Image&nbsp;File&nbsp;Format                         |<sub>`icc_profile`</sub>|
|`vorbis_comment`     |Vorbis&nbsp;comment                                          |<sub>`flac_picture`</sub>|
|`vorbis_packet`      |Vorbis&nbsp;packet                                           |<sub>`vorbis_comment`</sub>|
|`vp8_frame`          |VP8&nbsp;frame                                               |<sub></sub>|
|`vp9_frame`          |VP9&nbsp;frame                                               |<sub></sub>|
|`wav`                |WAV&nbsp;file                                                |<sub></sub>|
|`webp`               |WebP&nbsp;image                                              |<sub>`vp8_frame`</sub>|
|`image`              |Group                                                        |<sub>`gif`, `jpeg`, `png`, `tiff`, `webp`</sub>|
|`probe`              |Group                                                        |<sub>`aac_stream`, `adts`, `bzip2`, `elf`, `flac`, `gif`, `gzip`, `jpeg`, `mkv`, `mp3`, `mp4`, `ogg`, `png`, `tar`, `tiff`, `wav`, `webp`</sub>|

[#]: sh-end

TODO: format graph?

## Own decoders and use as library

TODO

## Known issues, TODOs and ideas

### Known issues

- TODO: concat bitbufs?
- TODO: byte units when outputting

### TODOs

- Function documentation in code, generate md etc
- Copy/pasteable output, add base prefixes
- Nested BitBufs, how to show? what about ranges? for example compressed data, demuxed ogg
- CRC fields, how to update with actual? fix flac
- Clean up panics, errors, better partial decode
- bitio.MultiBitReader to save memory
- Cleanup decoder API, nested bufs, decoders, try decode loop? decodebuf?
- Save encoding for values, LE, BE, varint etc
- Cleanup decoders
- Document decode maturity/completeness
- Arbitrary base integer literals
- Make jq functions that change state fail if called more than once? decode etc?
- REPL push/pop, variables etc?
- REPL tests
- Refactor *[]decode.Format into something more abstract, group?

### Ideas

- Suppose writing decoder in scripting language, jq, js, tango etc?
- Some kind of UI, web and cli? would be nice to visualize overlapping fields
- Is it possible to save memory by just record range/decoder at first decode and
then decode as needed later?
- Move more things to jq code, dumper?
- Some kind of bit vs bytes position notation/type
- Range/field user annotations

## Thanks

This project would not be possible without [itchyny](https://github.com/itchyny)'s
jq implementation [gojq](https://github.com/itchyny/gojq). Also want to thank
[HexFiend](https://github.com/HexFiend/HexFiend) for inspiration and ideas and
[stedolan](https://github.com/stedolan) for inventing the [jq](https://github.com/stedolan/jq)
language.
