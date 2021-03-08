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
     |00 01 02 03 04 05 06 07 08 09 0a 0b 0c 0d 0e 0f|                |.headers[0].frames[1]:
0x020|         41 50 49 43                           |   APIC         |  id: "APIC" (Attached picture)
0x020|                     00 00 15 0c               |       ....     |  size: 2700
0x020|                                 00 00         |           ..   | -flags:
0x020|                                       03      |             .  |  text_encoding: UTF-8 (3)
0x020|                                          69 6d|              im|  mime_type: "image/png"
0x030|61 67 65 2f 70 6e 67 00                        |age/png.        |
0x030|                        00                     |        .       |  picture_type: 0
0x030|                           00                  |         .      |  description: ""
0x030|                              89 50 4e 47 0d 0a|          .PNG..| -picture: png
0x040|1a 0a 00 00 00 0d 49 48 44 52 00 00 01 40 00 00|......IHDR...@..|
*    |2665 bytes more until 0xab8+7                  |                |
 
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

|Name                 |Description                                                  |Uses                                                                                                                                                                                                                                 |
|-|-|-|
|`apev2`              |APEv2&nbsp;metadata&nbsp;tag                                 |                                                                                                                                                                                                                                     |
|`av1_ccr`            |AV1&nbsp;codec&nbsp;configuration&nbsp;record                |                                                                                                                                                                                                                                     |
|`av1_frame`          |AV1&nbsp;frame                                               |`av1_obu`                                                                                                                                                                                                                            |
|`av1_obu`            |AV1&nbsp;open&nbsp;bitstream&nbsp;unit                       |                                                                                                                                                                                                                                     |
|`bzip2`              |bzip2&nbsp;compression                                       |`probe`                                                                                                                                                                                                                              |
|`dns`                |DNS&nbsp;packet                                              |                                                                                                                                                                                                                                     |
|`elf`                |Executable&nbsp;and&nbsp;Linkable&nbsp;Format                |                                                                                                                                                                                                                                     |
|`flac`               |Free&nbsp;lossless&nbsp;audio&nbsp;codec&nbsp;file           |`flac_metadatablock`, `flac_frame`                                                                                                                                                                                                   |
|`flac_frame`         |FLAC&nbsp;frame                                              |                                                                                                                                                                                                                                     |
|`flac_metadatablock` |FLAC&nbsp;metadatablock                                      |`flac_picture`, `vorbis_comment`                                                                                                                                                                                                     |
|`flac_picture`       |FLAC&nbsp;metadatablock&nbsp;picture                         |`image`                                                                                                                                                                                                                              |
|`gzip`               |GZIP&nbsp;compression                                        |`probe`                                                                                                                                                                                                                              |
|`icc_profile`        |International&nbsp;Color&nbsp;Consortium&nbsp;profile        |                                                                                                                                                                                                                                     |
|`id3_v1`             |ID3v1&nbsp;metadata                                          |                                                                                                                                                                                                                                     |
|`id3_v11`            |ID3v1.1&nbsp;metadata                                        |                                                                                                                                                                                                                                     |
|`id3_v2`             |ID3v2&nbsp;metadata                                          |`image`                                                                                                                                                                                                                              |
|`jpeg`               |Joint&nbsp;Photographic&nbsp;Experts&nbsp;Group&nbsp;file    |`tiff`                                                                                                                                                                                                                               |
|`jq`                 |                                                             |                                                                                                                                                                                                                                     |
|`mkv`                |Matroska                                                     |`av1_ccr`, `av1_frame`, `flac_frame`, `flac_metadatablock`, `mp3_frame`, `mpeg_aac_frame`, `mpeg_asc`, `mpeg_avc_dcr`, `mpeg_avc`, `mpeg_hevc_dcr`, `mpeg_hevc`, `mpeg_spu`, `opus_packet`, `vorbis_packet`, `vp8_frame`, `vp9_frame`|
|`mp3`                |MP3&nbsp;file                                                |`id3_v2`, `id3_v1`, `id3_v11`, `apev2`, `mp3_frame`                                                                                                                                                                                  |
|`mp3_frame`          |MPEG&nbsp;audio&nbsp;layer&nbsp;3&nbsp;frame                 |`mp3_xing`                                                                                                                                                                                                                           |
|`mp3_xing`           |Xing&nbsp;header                                             |                                                                                                                                                                                                                                     |
|`mp4`                |MPEG-4&nbsp;file                                             |`av1_ccr`, `av1_frame`, `flac_frame`, `flac_metadatablock`, `mp3_frame`, `mpeg_aac_frame`, `mpeg_avc_dcr`, `mpeg_avc`, `mpeg_es`, `mpeg_hevc_dcr`, `mpeg_hevc`, `opus_packet`, `vorbis_packet`, `vp9_frame`                          |
|`mpeg_aac_frame`     |Advanced&nbsp;Audio&nbsp;Coding&nbsp;frame                   |                                                                                                                                                                                                                                     |
|`mpeg_aac_stream`    |Raw&nbsp;audio&nbsp;data&nbsp;transport&nbsp;stream          |`mpeg_adts`                                                                                                                                                                                                                          |
|`mpeg_adts`          |Audio&nbsp;data&nbsp;transport&nbsp;stream&nbsp;packet       |`mpeg_aac_frame`                                                                                                                                                                                                                     |
|`mpeg_asc`           |MPEG-4&nbsp;Audio&nbsp;specific&nbsp;config                  |                                                                                                                                                                                                                                     |
|`mpeg_avc`           |H.264/AVC&nbsp;sample                                        |                                                                                                                                                                                                                                     |
|`mpeg_avc_dcr`       |H.264/AVC&nbsp;Decoder&nbsp;configuration&nbsp;record        |                                                                                                                                                                                                                                     |
|`mpeg_es`            |MPEG&nbsp;elementary&nbsp;stream                             |`mpeg_asc`, `vorbis_packet`                                                                                                                                                                                                          |
|`mpeg_hevc`          |H.265/HEVC&nbsp;sample                                       |                                                                                                                                                                                                                                     |
|`mpeg_hevc_dcr`      |H.265/HEVC&nbsp;Decoder&nbsp;configuration&nbsp;record       |                                                                                                                                                                                                                                     |
|`mpeg_pes`           |MPEG&nbsp;Packetized&nbsp;elementary&nbsp;stream             |`mpeg_pes_packet`, `mpeg_spu`                                                                                                                                                                                                        |
|`mpeg_pes_packet`    |MPEG&nbsp;Packetized&nbsp;elementary&nbsp;stream&nbsp;packet |                                                                                                                                                                                                                                     |
|`mpeg_spu`           |Sub&nbsp;picture&nbsp;unit&nbsp;(dvd&nbsp;subtitle)          |                                                                                                                                                                                                                                     |
|`ogg`                |OGG&nbsp;file                                                |`ogg_page`, `vorbis_packet`, `opus_packet`                                                                                                                                                                                           |
|`ogg_page`           |OGG&nbsp;page                                                |                                                                                                                                                                                                                                     |
|`opus_packet`        |Opus&nbsp;packet                                             |`vorbis_comment`                                                                                                                                                                                                                     |
|`png`                |Portable&nbsp;network&nbsp;graphics&nbsp;file                |`icc_profile`, `tiff`                                                                                                                                                                                                                |
|`raw`                |Raw&nbsp;bits                                                |                                                                                                                                                                                                                                     |
|`tar`                |Tar&nbsp;archive                                             |`probe`                                                                                                                                                                                                                              |
|`tiff`               |Tag&nbsp;Image&nbsp;File&nbsp;Format                         |`icc_profile`                                                                                                                                                                                                                        |
|`vorbis_comment`     |Vorbis&nbsp;comment                                          |`flac_picture`                                                                                                                                                                                                                       |
|`vorbis_packet`      |Vorbis&nbsp;packet                                           |`vorbis_comment`                                                                                                                                                                                                                     |
|`vp8_frame`          |VP8&nbsp;frame                                               |                                                                                                                                                                                                                                     |
|`vp9_frame`          |VP9&nbsp;frame                                               |                                                                                                                                                                                                                                     |
|`wav`                |WAV&nbsp;file                                                |                                                                                                                                                                                                                                     |
|`webp`               |WEBP&nbsp;image                                              |`vp8_frame`                                                                                                                                                                                                                          |
|`image`              |Group                                                        |`jpeg`, `png`, `tiff`, `webp`                                                                                                                                                                                                        |
|`probe`              |Group                                                        |`bzip2`, `elf`, `flac`, `gzip`, `jpeg`, `mkv`, `mp3`, `mp4`, `mpeg_aac_stream`, `mpeg_adts`, `ogg`, `png`, `tar`, `tiff`, `wav`, `webp`                                                                                              |

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
