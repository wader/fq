# fq

Tool, language and decoders for querying and exploring binary data.

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
$ fq file.mp4 '[.. | select(.type == "stsd")?.boxes[].type]' 
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

## Usage

Basic usage is `fq [OPTIONS] [FILE] [EXPR]...`.

For more details and support formats see [usage.md](doc/usage.md).

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
- https://github.com/binspector/binspector
