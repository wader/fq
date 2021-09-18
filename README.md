# fq

Tool, language and decoders for exploring binary data.

<sub>
<pre sh>
<b># Overview of mp3 file</b> 
$ fq . file.mp3 
     |00 01 02 03 04 05 06 07 08 09 0a 0b 0c 0d 0e 0f|0123456789abcdef|.: {} file.mp3 (mp3)
0x000|49 44 33 04 00 00 00 00 15 39 54 53 53 45 00 00|ID3......9TSSE..|  headers: [1]
*    |until 0xac2.7 (2755)                           |                |
0xac0|         ff fb 40 c0 00 00 00 00 00 00 00 00 00|   ..@..........|  frames: [3]
0xad0|00 00 00 00 00 00 00 00 49 6e 66 6f 00 00 00 0f|........Info....|
*    |until 0xd19.7 (end) (599)                      |                |
     |                                               |                |  footers: [0]
 
<b># Show ID3v2 tag in mp3 file</b> 
$ fq '.headers[0]' file.mp3 
     |00 01 02 03 04 05 06 07 08 09 0a 0b 0c 0d 0e 0f|0123456789abcdef|.headers[0]: {} (id3v2)
0x000|49 44 33                                       |ID3             |  magic: "ID3" (Correct)
0x000|         04                                    |   .            |  version: 4
0x000|            00                                 |    .           |  revision: 0
0x000|               00                              |     .          |  flags: {}
0x000|                  00 00 15 39                  |      ...9      |  size: 2745
0x000|                              54 53 53 45 00 00|          TSSE..|  frames: [2]
0x010|00 0f 00 00 03 4c 61 76 66 35 38 2e 37 36 2e 31|.....Lavf58.76.1|
*    |until 0xab8.7 (2735)                           |                |
0xab0|                           00 00 00 00 00 00 00|         .......|  padding: Correct (none) (zero padding)
0xac0|00 00 00                                       |...             |
 
<b># Resolution of ID3v2 cover art</b> 
$ fq '.headers[0].frames[] | select(.id == "APIC").picture.chunks[] | select(.type == "IHDR") | {width, height}' file.mp3 
{
  "height": 240,
  "width": 320
}
 
<b># Extract image file</b> 
$ fq '.headers[].frames[] | select(.id == "APIC")?.picture | tobits' file.mp3 > file.png 
$ file file.png 
file.png: PNG image data, 320 x 240, 8-bit/color RGB, non-interlaced
 
</pre>
</sub>

## Goals

- Make structured binary data accessible using the jq language.
- Nested format decoding and bit-oriented decoding.
- Quick cli tool that tries to mimic jq as much as possible.

**NOTE** fq is early in development and many things are missing, broken or do not make sense. This also means there is much opportunity to help out!

## Install

Download archive from [releases](https://github.com/wader/fq/releases) page for your
platform, unarchive and move the executable to `PATH`.

### Homebrew

```sh
# install latest release
brew install wader/tap/fq
```

### Build from source

Make sure you have go 1.17 or later and then do:
```sh
# build and install latest master
go install github.com/wader/fq@latest
```
and the binary should end up at `$GOPATH/bin/fq`.

## Usage

Basic usage is:

[fq -h | grep Usage: | sed 's/\(.*\)/<pre>\1<\/pre>/']: sh-start

<pre>Usage: fq [OPTIONS] [--] [EXPR] [FILE...]</pre>

[#]: sh-end

For more usage details see [usage.md](doc/usage.md).

## Supported formats

[./formats_list.jq]: sh-start

aac_frame, adts, adts_frame, apev2, av1_ccr, av1_frame, av1_obu, avc_annexb, avc_au, avc_dcr, avc_nalu, avc_pps, avc_sei, avc_sps, bzip2, dns, elf, exif, flac, flac_frame, flac_metadatablocks, flac_picture, gif, gzip, hevc_annexb, hevc_au, hevc_dcr, hevc_nalu, icc_profile, id3v1, id3v11, id3v2, jpeg, json, matroska, mp3, mp3_frame, mp4, mpeg_asc, mpeg_es, mpeg_pes, mpeg_pes_packet, mpeg_spu, mpeg_ts, ogg, ogg_page, opus_packet, png, protobuf, protobuf_widevine, pssh_playready, raw, tar, tiff, vorbis_comment, vorbis_packet, vp8_frame, vp9_cfm, vp9_frame, vpx_ccr, wav, webp, xing

[#]: sh-end


For more format details see [usage.md](doc/usage.md).

## TODO and ideas

See [TODO.md](doc/TODO.md)

## Development

See [dev.md](doc/dev.md)

## Thanks and related projects

This project would not be possible without [itchyny](https://github.com/itchyny)'s
jq implementation [gojq](https://github.com/itchyny/gojq). Also want to thank
[HexFiend](https://github.com/HexFiend/HexFiend) for inspiration and ideas and [stedolan](https://github.com/stedolan)
for inventing the [jq](https://github.com/stedolan/jq) language.

Similar projects:
- https://github.com/HexFiend/HexFiend
- https://github.com/binspector/binspector
- https://kaitai.io
