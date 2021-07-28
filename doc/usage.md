### Arguments

<pre sh>
$ fq -h 
Usage: fq [OPTIONS] [EXPR] [FILE...]
--compact,-c     Compact output
--decode,-d      Decode format (probe)
--file,-f        Read script from file
--formats        Show formats
--help,-h        Show help
-n               Null input
--option,-o=ARG  Set option, eg: color=true
                   addrbase=16
                   bytecolors=0-0xff=brightwhite,0=brightblack,32-126:9-13=white
                   color=false
                   colors=array=white,dumpaddr=yellow,dumpheader=yellow+underline,error=brightred,false=yellow,index=white,null=brightblack,number=cyan,object=white,objectkey=brightblue,string=green,true=yellow,value=white
                   depth=0
                   displaybytes=16
                   linebytes=16
                   raw=true
                   sizebase=10
                   unicode=false
                   verbose=false
-r               Raw strings
--repl,-i        Interactive REPL
--version,-v     Show version (dev)
</pre>

- TODO: null input
- TODO: expressions

### Running

- TODO: stdin/stdout

### Interactive REPL

- TODO: tab completion, ctrl-d, ctrl-d, help
- TODO: nested, nested with generator

### Script

- TODO: #!

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

- [gojq's differences to jq](https://github.com/itchyny/gojq#difference-to-jq),
notable is support for arbitrary-precision integers.
- Supports hexdecimal `0xab`, octal `0o77` and binary `0b101` integer literals
- Has bitwise operations, `band`, `bor`, `bxor`, `bsl`, `bsr`, `bnot`
- Has `div` integer division operator
- Try include `include "file?";` that don't fail if file is missing
- Possible for a value to act as a object with keys even when it's an array, number etc.
- There can be keys hidden from `keys` and `[]`. Used for, `_format`, `_bytes` etc.
- Some values do not support to be updated

### Functions

- All standard library functions from jq
- `open(path)` open file for reading
- `probe` or `decode` try to automatically detect format and decode
- `mp3`, `matroska`, ..., `<name>`, `decode([name])` try decode as format
- `d`/`display` display value
- `v`/`verbose` display value verbosely
- `p`/`preview` show preview of field tree
- `hd`/`hexdump` hexdump value
- `repl` nested REPL

### Decoded values (TODO: better name?)

When you decode something successfully in fq you will get a value. A value work a bit like
jq object with special abilities and is used to represent a tree structure of the decoded
binary data. Each value always has a name, type and a bit range.

A value has these special keys:

- `_name` name of value
- `_value` jq value of value
- `_start` bit range start
- `_stop` bit range stop
- `_len` bit range length (TODO: rename)
- `_bits` bits in range as a binary
- `_bytes` bits in range as binary using byte units
- `_path` jq path to value
- `_unknown` value is un-decoded gap
- `_symbol` symbolic string representation of value (optional)
- `_description` longer description of value (optional)
- `_format` name of decoded format (optional)
- `_error` error message (optional)

- TODO: unknown gaps

### Binary and IO lists

- TODO: similar to erlang io lists, [], binary, string (utf8) and numbers

## Configuration

To add own functions you can use `init.fq` that will be read from
- `$HOME/Library/Application Support/fq/init.jq` on macOS
- `$HOME/.config/fq/init.jq` on Linux, BSD etc
- `%AppData%\fq\init.jq` on Windows (TODO: not tested)

## Decoders

[./formats_table.jq]: sh-start

|Name                 |Description                                                   |Uses|
|-                    |-                                                             |-|
|`aac_frame`          |Advanced&nbsp;Audio&nbsp;Coding&nbsp;frame                    |<sub></sub>|
|`adts`               |Audio&nbsp;Data&nbsp;Transport&nbsp;Stream                    |<sub>`adts_frame`</sub>|
|`adts_frame`         |Audio&nbsp;Data&nbsp;Transport&nbsp;Stream&nbsp;frame         |<sub>`aac_frame`</sub>|
|`apev2`              |APEv2&nbsp;metadata&nbsp;tag                                  |<sub>`image`</sub>|
|`av1_ccr`            |AV1&nbsp;Codec&nbsp;Configuration&nbsp;Record                 |<sub></sub>|
|`av1_frame`          |AV1&nbsp;frame                                                |<sub>`av1_obu`</sub>|
|`av1_obu`            |AV1&nbsp;Open&nbsp;Bitstream&nbsp;Unit                        |<sub></sub>|
|`avc_au`             |H.264/AVC&nbsp;Access&nbsp;Unit                               |<sub>`avc_nalu`</sub>|
|`avc_dcr`            |H.264/AVC&nbsp;Decoder&nbsp;Configuration&nbsp;Record         |<sub>`avc_nalu`</sub>|
|`avc_nalu`           |H.264/AVC&nbsp;Network&nbsp;Access&nbsp;Layer&nbsp;Unit       |<sub>`avc_sps` `avc_pps` `avc_sei`</sub>|
|`avc_pps`            |H.264/AVC&nbsp;Picture&nbsp;Parameter&nbsp;Set                |<sub></sub>|
|`avc_sei`            |H.264/AVC&nbsp;Supplemental&nbsp;Enhancement&nbsp;Information |<sub></sub>|
|`avc_sps`            |H.264/AVC&nbsp;Sequence&nbsp;Parameter&nbsp;Set               |<sub></sub>|
|`bzip2`              |bzip2&nbsp;compression                                        |<sub>`probe`</sub>|
|`dns`                |DNS&nbsp;packet                                               |<sub></sub>|
|`elf`                |Executable&nbsp;and&nbsp;Linkable&nbsp;Format                 |<sub></sub>|
|`exif`               |Exchangeable&nbsp;Image&nbsp;File&nbsp;Format                 |<sub></sub>|
|`flac`               |Free&nbsp;Lossless&nbsp;Audio&nbsp;Codec&nbsp;file            |<sub>`flac_metadatablock` `flac_frame`</sub>|
|`flac_frame`         |FLAC&nbsp;frame                                               |<sub></sub>|
|`flac_metadatablock` |FLAC&nbsp;metadatablock                                       |<sub>`flac_picture` `vorbis_comment`</sub>|
|`flac_picture`       |FLAC&nbsp;metadatablock&nbsp;picture                          |<sub>`image`</sub>|
|`gif`                |Graphics&nbsp;Interchange&nbsp;Format                         |<sub></sub>|
|`gzip`               |gzip&nbsp;compression                                         |<sub>`probe`</sub>|
|`hevc_au`            |H.265/HEVC&nbsp;Access&nbsp;Unit                              |<sub>`hevc_nalu`</sub>|
|`hevc_dcr`           |H.265/HEVC&nbsp;Decoder&nbsp;Configuration&nbsp;Record        |<sub>`hevc_nalu`</sub>|
|`hevc_nalu`          |H.265/HEVC&nbsp;Network&nbsp;Access&nbsp;Layer&nbsp;Unit      |<sub></sub>|
|`icc_profile`        |International&nbsp;Color&nbsp;Consortium&nbsp;profile         |<sub></sub>|
|`id3v1`              |ID3v1&nbsp;metadata                                           |<sub></sub>|
|`id3v11`             |ID3v1.1&nbsp;metadata                                         |<sub></sub>|
|`id3v2`              |ID3v2&nbsp;metadata                                           |<sub>`image`</sub>|
|`jpeg`               |Joint&nbsp;Photographic&nbsp;Experts&nbsp;Group&nbsp;file     |<sub>`exif` `icc_profile`</sub>|
|`matroska`           |Matroska&nbsp;file                                            |<sub>`aac_frame` `av1_ccr` `av1_frame` `flac_frame` `flac_metadatablock` `mp3_frame` `mpeg_asc` `avc_au` `avc_dcr` `hevc_au` `hevc_dcr` `mpeg_pes_packet` `mpeg_spu` `opus_packet` `vorbis_packet` `vp8_frame` `vp9_cfm` `vp9_frame`</sub>|
|`mp3`                |MP3&nbsp;file                                                 |<sub>`id3v2` `id3v1` `id3v11` `apev2` `mp3_frame`</sub>|
|`mp3_frame`          |MPEG&nbsp;audio&nbsp;layer&nbsp;3&nbsp;frame                  |<sub>`xing`</sub>|
|`mp4`                |MPEG-4&nbsp;file&nbsp;and&nbsp;similar                        |<sub>`aac_frame` `av1_ccr` `av1_frame` `flac_frame` `flac_metadatablock` `id3v2` `image` `jpeg` `mp3_frame` `avc_au` `avc_dcr` `mpeg_es` `hevc_au` `hevc_dcr` `mpeg_pes_packet` `opus_packet` `protobuf_widevine` `vorbis_packet` `vp9_frame` `vpx_ccr`</sub>|
|`mpeg_asc`           |MPEG-4&nbsp;Audio&nbsp;Specific&nbsp;Config                   |<sub></sub>|
|`mpeg_es`            |MPEG&nbsp;Elementary&nbsp;Stream                              |<sub>`mpeg_asc` `vorbis_packet`</sub>|
|`mpeg_pes`           |MPEG&nbsp;Packetized&nbsp;elementary&nbsp;stream              |<sub>`mpeg_pes_packet` `mpeg_spu`</sub>|
|`mpeg_pes_packet`    |MPEG&nbsp;Packetized&nbsp;elementary&nbsp;stream&nbsp;packet  |<sub></sub>|
|`mpeg_spu`           |Sub&nbsp;Picture&nbsp;Unit&nbsp;(DVD&nbsp;subtitle)           |<sub></sub>|
|`mpeg_ts`            |MPEG&nbsp;Transport&nbsp;Stream                               |<sub></sub>|
|`ogg`                |OGG&nbsp;file                                                 |<sub>`ogg_page` `vorbis_packet` `opus_packet`</sub>|
|`ogg_page`           |OGG&nbsp;page                                                 |<sub></sub>|
|`opus_packet`        |Opus&nbsp;packet                                              |<sub>`vorbis_comment`</sub>|
|`png`                |Portable&nbsp;Network&nbsp;Graphics&nbsp;file                 |<sub>`icc_profile` `exif`</sub>|
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
|`wav`                |WAV&nbsp;file                                                 |<sub>`id3v2` `id3v1` `id3v11`</sub>|
|`webp`               |WebP&nbsp;image                                               |<sub>`vp8_frame`</sub>|
|`xing`               |Xing&nbsp;header                                              |<sub></sub>|
|`image`              |Group                                                         |<sub>`gif` `jpeg` `png` `tiff` `webp`</sub>|
|`probe`              |Group                                                         |<sub>`adts` `bzip2` `elf` `flac` `gif` `gzip` `jpeg` `matroska` `mp3` `mp4` `mpeg_ts` `ogg` `png` `tar` `tiff` `wav` `webp`</sub>|

[#]: sh-end

TODO: format graph?

## Own decoders and use as library

TODO


### Useful tricks

#### `.. | select(...)` fails with `expected an ... but got: ...`

Try add `select(...)?` the select expression assumes it will get and object etc.

#### Manual decode

Sometimes fq fails to decode or you know there is valid data buried inside some binary or maybe
you know the format of some unknown value. Then you can decode manually.

<pre>
# try decode a `mp3_frame` that failed to decode
$ fq file.mp3 .unknown0 mp3_frame
# skip first 10 bytes then decode as `mp3_frame`
$ fq file.mp3 .unknown0._bytes[10:] mp3_frame
</pre>

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
