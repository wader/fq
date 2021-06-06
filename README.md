# fq

Tool, language and decoders for querying and exploring binary data.

## Usage

<sub>
<pre sh>
<b># Overview of mp3 file</b> 
$ fq file.mp3 
     |00 01 02 03 04 05 06 07 08 09 0a 0b 0c 0d 0e 0f|                |.: {} (mp3)
0x000|49 44 33 04 00 00 00 00 15 39 54 53 53 45 00 00|ID3......9TSSE..|  headers: [1]
*    |until 0xac2.7 (2755)                           |                |
0xac0|         ff fb 40 c0 00 00 00 00 00 00 00 00 00|   ..@..........|  frames: [3]
0xad0|00 00 00 00 00 00 00 00 49 6e 66 6f 00 00 00 0f|........Info....|
*    |until 0xd19.7 (end) (599)                      |                |
     |                                               |                |  footers: [0]
 
<b># Show ID3v2 APIC frame</b> 
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
0x030|                              89 50 4e 47 0d 0a|          .PNG..|  picture: {} (png)
0x040|1a 0a 00 00 00 0d 49 48 44 52 00 00 01 40 00 00|......IHDR...@..|
*    |until 0xab8.7 (2687)                           |                |
 
<b># Resolution of embedded PNG file</b> 
$ fq file.mp3 '.headers[].frames[] | select(.id == "APIC").picture.chunks[] | select(.type == "IHDR") | {width, height}' 
{
  "height": 240,
  "width": 320
}
 
<b># Extract PNG file</b> 
$ fq file.mp3 '.headers[].frames[] | select(.id == "APIC")?.picture._bits' > file.png 
$ file file.png 
file.png: PNG image data, 320 x 240, 8-bit/color RGB, non-interlaced
 
<b># Codecs in a mp4 file</b> 
$ fq file.mp4 '[.. | select(.type == "stsd")?.sample_descriptions[].data_format]' 
[
  "avc1",
  "mp4a"
]
</pre>
</sub>

## Install

Currently there are no binary releases, but it's quite easy to build fq yourself. Make sure you have go 1.16
or later and then do:
```sh
go install github.com/wader/fq@latest
```
and the binary should end up at `$GOPATH/bin/fq`.

## Langauge

fq is based on the [jq language](https://stedolan.github.io/jq/) and for basic usage its syntax
is similar to how object and array access looks in JavaScript or JSON path, `.food[10]` etc.

To get the most out of fq it's recommended to learn more about jq, here are some good starting points:

- [jq manual](https://stedolan.github.io/jq/manual/)
- [jq Cookbook](https://github.com/stedolan/jq/wiki/Cookbook),
[FAQ](https://github.com/stedolan/jq/wiki/FAQ),
[Pitfalls](https://github.com/stedolan/jq/wiki/How-to:-Avoid-Pitfalls)

The most common beginner gotcha is probably jq's use of `;` and `,`. jq uses `;` as argument separator.
To call `f` with two arguments use `f(a; b)`. If you do `f(a, b)` you will pass a single generator
expression `a, b` to `f`.

### Differences to jq

- Supports hexdecimal `0xab`, octal `0o77` and binary `0b101` integer literals
- Has bitwise operations, `band`, `bor`, `bxor`, `bsl`, `bsr`, `bnot`
- Has `div` integer division operator
- Try include `include "file?";` that don't fail if file is missing
- Can have keys that are hidden from `keys` and `[]` used for some `_name` proprties like `_bytes`
- Value can be hybrid array and object at the same time

### Functions

An addition to the standard library functions from jq fq has these functions:

- `open(path)` open file
- `probe` or `decode` try to automatically detect format and decode
- `mp3`, `matroska`, ..., `decode([name])` try decode as format
- `d`/`display` display value
- `v`/`verbose` display value verbosely
- `p`/`preview` show preview of field tree
- `repl` nested REPL

## Configuration

To add own functions you can use `init.fq` that will be read from
- `$HOME/Library/Application Support/fq/init.jq` on macOS
- `$HOME/.config/fq/init.jq` on Linux, BSD etc
- `%AppData%\fq\init.jq` on Windows (TODO: not tested)

## How to use

TODO: unknown for gaps
TODO: piping

## Decoders

[./formats_markdown.jq]: sh-start

|Name                 |Description                                                   |Uses|
|-                    |-                                                             |-|
|`aac_frame`          |Advanced&nbsp;Audio&nbsp;Coding&nbsp;frame                    |<sub></sub>|
|`aac_stream`         |Raw&nbsp;audio&nbsp;data&nbsp;transport&nbsp;stream           |<sub>`adts`</sub>|
|`adts`               |Audio&nbsp;Data&nbsp;Transport&nbsp;Stream&nbsp;packet        |<sub>`aac_frame`</sub>|
|`apev2`              |APEv2&nbsp;metadata&nbsp;tag                                  |<sub></sub>|
|`av1_ccr`            |AV1&nbsp;Codec&nbsp;Configuration&nbsp;Record                 |<sub></sub>|
|`av1_frame`          |AV1&nbsp;frame                                                |<sub>`av1_obu`</sub>|
|`av1_obu`            |AV1&nbsp;Open&nbsp;Bitstream&nbsp;Unit                        |<sub></sub>|
|`avc_au`             |H.264/AVC&nbsp;Access&nbsp;Unit                               |<sub>`avc_nalu`</sub>|
|`avc_dcr`            |H.264/AVC&nbsp;Decoder&nbsp;Configuration&nbsp;Record         |<sub>`avc_nalu`</sub>|
|`avc_nalu`           |H.264/AVC&nbsp;Network&nbsp;Access&nbsp;Layer&nbsp;Unit       |<sub>`avc_sps`, `avc_pps`, `avc_sei`</sub>|
|`avc_pps`            |H.264/AVC&nbsp;Picture&nbsp;Parameter&nbsp;Set                |<sub></sub>|
|`avc_sei`            |H.264/AVC&nbsp;Supplemental&nbsp;Enhancement&nbsp;Information |<sub></sub>|
|`avc_sps`            |H.264/AVC&nbsp;Sequence&nbsp;Parameter&nbsp;Set               |<sub></sub>|
|`bzip2`              |bzip2&nbsp;compression                                        |<sub>`probe`</sub>|
|`dns`                |DNS&nbsp;packet                                               |<sub></sub>|
|`elf`                |Executable&nbsp;and&nbsp;Linkable&nbsp;Format                 |<sub></sub>|
|`exif`               |Exchangeable&nbsp;Image&nbsp;File&nbsp;Format                 |<sub>`icc_profile`</sub>|
|`flac`               |Free&nbsp;Lossless&nbsp;Audio&nbsp;Codec&nbsp;file            |<sub>`flac_metadatablock`, `flac_frame`</sub>|
|`flac_frame`         |FLAC&nbsp;frame                                               |<sub></sub>|
|`flac_metadatablock` |FLAC&nbsp;metadatablock                                       |<sub>`flac_picture`, `vorbis_comment`</sub>|
|`flac_picture`       |FLAC&nbsp;metadatablock&nbsp;picture                          |<sub>`image`</sub>|
|`gif`                |Graphics&nbsp;Interchange&nbsp;Format                         |<sub></sub>|
|`gzip`               |gzip&nbsp;compression                                         |<sub>`probe`</sub>|
|`hevc_au`            |H.265/HEVC&nbsp;Access&nbsp;Unit                              |<sub>`hevc_nalu`</sub>|
|`hevc_dcr`           |H.265/HEVC&nbsp;Decoder&nbsp;Configuration&nbsp;Record        |<sub>`hevc_nalu`</sub>|
|`hevc_nalu`          |H.265/HEVC&nbsp;Network&nbsp;Access&nbsp;Layer&nbsp;Unit      |<sub></sub>|
|`icc_profile`        |International&nbsp;Color&nbsp;Consortium&nbsp;profile         |<sub></sub>|
|`id3_v1`             |ID3v1&nbsp;metadata                                           |<sub></sub>|
|`id3_v11`            |ID3v1.1&nbsp;metadata                                         |<sub></sub>|
|`id3_v2`             |ID3v2&nbsp;metadata                                           |<sub>`image`</sub>|
|`jpeg`               |Joint&nbsp;Photographic&nbsp;Experts&nbsp;Group&nbsp;file     |<sub>`exif`, `icc_profile`</sub>|
|`matroska`           |Matroska&nbsp;file                                            |<sub>`av1_ccr`, `av1_frame`, `flac_frame`, `flac_metadatablock`, `mp3_frame`, `aac_frame`, `mpeg_asc`, `avc_dcr`, `avc_au`, `hevc_dcr`, `hevc_au`, `mpeg_spu`, `mpeg_pes_packet`, `opus_packet`, `vorbis_packet`, `vp8_frame`, `vp9_frame`, `vp9_cfm`</sub>|
|`mp3`                |MP3&nbsp;file                                                 |<sub>`id3_v2`, `id3_v1`, `id3_v11`, `apev2`, `mp3_frame`</sub>|
|`mp3_frame`          |MPEG&nbsp;audio&nbsp;layer&nbsp;3&nbsp;frame                  |<sub>`mp3_xing`</sub>|
|`mp3_xing`           |Xing&nbsp;header                                              |<sub></sub>|
|`mp4`                |MPEG-4&nbsp;file                                              |<sub>`av1_ccr`, `av1_frame`, `flac_frame`, `flac_metadatablock`, `mp3_frame`, `aac_frame`, `avc_dcr`, `avc_au`, `mpeg_es`, `hevc_dcr`, `hevc_au`, `mpeg_pes_packet`, `opus_packet`, `vorbis_packet`, `vp9_frame`, `vpx_ccr`, `jpeg`, `id3_v2`, `protobuf_widevine`</sub>|
|`mpeg_asc`           |MPEG-4&nbsp;Audio&nbsp;Specific&nbsp;Config                   |<sub></sub>|
|`mpeg_es`            |MPEG&nbsp;Elementary&nbsp;Stream                              |<sub>`mpeg_asc`, `vorbis_packet`</sub>|
|`mpeg_pes`           |MPEG&nbsp;Packetized&nbsp;elementary&nbsp;stream              |<sub>`mpeg_pes_packet`, `mpeg_spu`</sub>|
|`mpeg_pes_packet`    |MPEG&nbsp;Packetized&nbsp;elementary&nbsp;stream&nbsp;packet  |<sub></sub>|
|`mpeg_spu`           |Sub&nbsp;Picture&nbsp;Unit&nbsp;(DVD&nbsp;subtitle)           |<sub></sub>|
|`mpeg_ts`            |MPEG&nbsp;Transport&nbsp;Stream                               |<sub></sub>|
|`ogg`                |OGG&nbsp;file                                                 |<sub>`ogg_page`, `vorbis_packet`, `opus_packet`</sub>|
|`ogg_page`           |OGG&nbsp;page                                                 |<sub></sub>|
|`opus_packet`        |Opus&nbsp;packet                                              |<sub>`vorbis_comment`</sub>|
|`png`                |Portable&nbsp;Network&nbsp;Graphics&nbsp;file                 |<sub>`icc_profile`, `exif`</sub>|
|`protobuf`           |Protobuf                                                      |<sub></sub>|
|`protobuf_widevine`  |Widevine&nbsp;protobuf                                        |<sub>`protobuf`</sub>|
|`raw`                |Raw&nbsp;bits                                                 |<sub></sub>|
|`tar`                |Tar&nbsp;archive                                              |<sub>`probe`</sub>|
|`tiff`               |Tag&nbsp;Image&nbsp;File&nbsp;Format                          |<sub>`icc_profile`</sub>|
|`vorbis_comment`     |Vorbis&nbsp;comment                                           |<sub>`flac_picture`</sub>|
|`vorbis_packet`      |Vorbis&nbsp;packet                                            |<sub>`vorbis_comment`</sub>|
|`vp8_frame`          |VP8&nbsp;frame                                                |<sub></sub>|
|`vp9_cfm`            |VP9&nbsp;Codec&nbsp;Feature&nbsp;Metadata                     |<sub></sub>|
|`vp9_frame`          |VP9&nbsp;frame                                                |<sub></sub>|
|`vpx_ccr`            |VPX&nbsp;Codec&nbsp;Configuration&nbsp;Record                 |<sub></sub>|
|`wav`                |WAV&nbsp;file                                                 |<sub>`id3_v2`, `id3_v1`, `id3_v11`</sub>|
|`webp`               |WebP&nbsp;image                                               |<sub>`vp8_frame`</sub>|
|`image`              |Group                                                         |<sub>`gif`, `jpeg`, `png`, `tiff`, `webp`</sub>|
|`probe`              |Group                                                         |<sub>`aac_stream`, `adts`, `bzip2`, `elf`, `flac`, `gif`, `gzip`, `jpeg`, `matroska`, `mp3`, `mp4`, `mpeg_ts`, `ogg`, `png`, `tar`, `tiff`, `wav`, `webp`</sub>|

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

### Useful tricks

#### Run pipelines using CLI arguments
<pre sh>
$ fq file.mp3 .frames[0].header.bitrate radix2 
"1101101011000000"
</pre>
instead of:
<pre sh>
$ fq file.mp3 '.frames[0].header.bitrate | radix2' 
"1101101011000000"
</pre>
this can also be used with interactive mode
```sh
$ fq -i file.flac .metadatablocks[0] 
.metadatablocks[0] flac_metadatablock> 
```

#### appending to array is slow

Try to use `map` or `foreach` instead.

#### Use `print` and `println` to produce more friendly compact output

```
> [[0,"a"],[1,"b"]]
[
  [
    0,
    "a"
  ],
  [
    1,
    "b"
  ]
]
> [[0,"a"],[1,"b"]] | .[] | "\(.[0]): \(.[1])" | println
0: a
1: b
```

#### Run interactive mode with no input
```sh
fq -i
null>
```

### Ideas

- Suppose writing decoder in scripting language, jq, js, tango etc?
- Some kind of UI, web and cli? would be nice to visualize overlapping fields
- Is it possible to save memory by just record range/decoder at first decode and
then decode as needed later?
- Move more things to jq code, dumper?
- Some kind of bit vs bytes position notation/type
- Range/field user annotations

## Development

- TODO: `scope` and `scopedump` functions used to implement REPL completion
- TODO: Custom object interface used to traverse fq's field tree and to allowing a terse
syntax for comparing and working with fields, accessing child fields and special properties like `_range`.

## Thanks and related projects

This project would not be possible without [itchyny](https://github.com/itchyny)'s
jq implementation [gojq](https://github.com/itchyny/gojq).

Also want to thank
- [HexFiend](https://github.com/HexFiend/HexFiend) for inspiration and ideas
- [stedolan](https://github.com/stedolan) for inventing the [jq](https://github.com/stedolan/jq)
language.
