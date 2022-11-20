# fq

Tool, language and decoders for working with binary data.

![fq demo](doc/demo.svg)

Basic usage is `fq . file` or `fq d file`.

For details see [usage.md](doc/usage.md).

### Background

fq is inspired by the well known jq tool and language and allows you to work with binary formats the same way you would using jq. In addition it can present data like a hex viewer, transform, slice and concatenate binary data. It also supports nested formats and has an interactive REPL with auto-completion.

It was originally designed to query, inspect and debug media codecs and containers like mp4, flac, mp3, jpeg. But has since then been extended to support a variety of formats like executables, packet captures (with TCP reassembly) and serialization formats like JSON, YAML, XML, ASN1 BER, Avro, CBOR, protobuf. In addition it also has functions to work with URL:s, convert to/from hex, number bases, search for things etc.

In summary it aims to be jq, hexdump, dd and gdb for files combined into one.

**NOTE:** fq is still early in development so things might change, be broken or do not make sense. That also means that there is a great opportunity to help out!

### Goals

- Make binaries accessible, queryable and sliceable.
- Nested formats and bit-oriented decoding.
- Quick and comfortable CLI tool.
- Bits and bytes transformations.

### Hopes

- Make it useful enough that people want to help improve it.
- Inspire people to create similar tools.

### Supported formats

[fq -rn -L doc 'include "formats"; formats_list']: sh-start

[aac_frame](doc/formats.md#aac_frame),
adts,
adts_frame,
amf0,
apev2,
ar,
[asn1_ber](doc/formats.md#asn1_ber),
av1_ccr,
av1_frame,
av1_obu,
avc_annexb,
[avc_au](doc/formats.md#avc_au),
avc_dcr,
avc_nalu,
avc_pps,
avc_sei,
avc_sps,
[avi](doc/formats.md#avi),
[avro_ocf](doc/formats.md#avro_ocf),
[bencode](doc/formats.md#bencode),
bitcoin_blkdat,
[bitcoin_block](doc/formats.md#bitcoin_block),
bitcoin_script,
bitcoin_transaction,
[bits](doc/formats.md#bits),
[bplist](doc/formats.md#bplist),
bsd_loopback_frame,
[bson](doc/formats.md#bson),
[bytes](doc/formats.md#bytes),
bzip2,
[cbor](doc/formats.md#cbor),
[csv](doc/formats.md#csv),
dns,
dns_tcp,
elf,
ether8023_frame,
exif,
fairplay_spc,
flac,
[flac_frame](doc/formats.md#flac_frame),
flac_metadatablock,
flac_metadatablocks,
flac_picture,
flac_streaminfo,
gif,
gzip,
hevc_annexb,
[hevc_au](doc/formats.md#hevc_au),
hevc_dcr,
hevc_nalu,
hevc_pps,
hevc_sps,
hevc_vps,
[html](doc/formats.md#html),
icc_profile,
icmp,
icmpv6,
id3v1,
id3v11,
id3v2,
ipv4_packet,
ipv6_packet,
jpeg,
json,
jsonl,
[macho](doc/formats.md#macho),
macho_fat,
[markdown](doc/formats.md#markdown),
[matroska](doc/formats.md#matroska),
[mp3](doc/formats.md#mp3),
mp3_frame,
mp3_frame_tags,
[mp4](doc/formats.md#mp4),
mpeg_asc,
mpeg_es,
mpeg_pes,
mpeg_pes_packet,
mpeg_spu,
mpeg_ts,
[msgpack](doc/formats.md#msgpack),
ogg,
ogg_page,
opus_packet,
[pcap](doc/formats.md#pcap),
pcapng,
png,
prores_frame,
[protobuf](doc/formats.md#protobuf),
protobuf_widevine,
pssh_playready,
[rtmp](doc/formats.md#rtmp),
sll2_packet,
sll_packet,
tar,
tcp_segment,
tiff,
toml,
udp_datagram,
vorbis_comment,
vorbis_packet,
vp8_frame,
vp9_cfm,
vp9_frame,
vpx_ccr,
[wasm](doc/formats.md#wasm),
wav,
webp,
[xml](doc/formats.md#xml),
yaml,
[zip](doc/formats.md#zip)

[#]: sh-end

It can also work with some common text formats like URL:s, hex, base64, PEM etc and for some serialization formats like XML, YAML etc it can transform both from and to jq values.

For details see [formats.md](doc/formats.md) and [usage.md](doc/usage.md).

## Presentations

- "fq - jq for binary formats" at [No time to wait 6](https://mediaarea.net/NoTimeToWait6) - [video](https://www.youtube.com/watch?v=-Pwt5KL-xRs&t=1450s) - [slides](doc/presentations/nttw6/fq-nttw6-slides.pdf)
- "fq - jq for binary formats" at [Binary Tools Summit 2022](https://binary-tools.net/summit.html) - [video](https://www.youtube.com/watch?v=GJOq_b0eb-s&list=PLTj8twuHdQz-JcX7k6eOwyVPDB8CyfZc8&index=1) - [slides](doc/presentations/bts2022/fq-bts2022-v1.pdf)

## Install

Use one of the methods listed below or download a pre-built [release](https://github.com/wader/fq/releases) for macOS, Linux or Windows. Unarchive it and move the executable to `PATH` etc.

On macOS if you don't install using a method below you might have to manually allow the binary to run. This can be done by trying to run the binary, ignore the warning and then go into security preference and allow it. Same can be done with this command:

```sh
xattr -d com.apple.quarantine fq && spctl --add fq
```

### Homebrew (macOS)

```sh
brew install wader/tap/fq
```

### MacPorts

On macOS, `fq` can also be installed via [MacPorts](https://www.macports.org).  More details [here](https://ports.macports.org/port/fq/).

```sh
sudo port install fq
```

### Windows

`fq` can be installed via [scoop](https://scoop.sh/).

```powershell
scoop install fq
```

### Arch Linux

`fq` can be installed from the [community repository](https://archlinux.org/packages/community/x86_64/fq/) using [pacman](https://wiki.archlinux.org/title/Pacman):

```sh
pacman -S fq
```

You can also build and install the development (VCS) package using an [AUR helper](https://wiki.archlinux.org/index.php/AUR_helpers):

```sh
paru -S fq-git
```

### Nix

```sh
nix-shell -p fq
```

### FreeBSD

Use the [fq](https://cgit.freebsd.org/ports/tree/misc/fq) port.

### Alpine

Currently in edge testing but should work fine in stable also.

```
apk add -X http://dl-cdn.alpinelinux.org/alpine/edge/testing fq
```

### Build from source

Make sure you have [go](https://go.dev) 1.18 or later installed.

To install directly from git repository (no clone needed) do:
```sh
# build and install latest release
go install github.com/wader/fq@latest

# build and install latest master
go install github.com/wader/fq@master

# copy binary to $PATH if needed
cp "$(go env GOPATH)/bin/fq" /usr/local/bin
```

To build, run and test from source:
```sh
# build an run
go run fq.go
# or
go run .
# just build
go build -o fq .
# run all tests and build binary
make test fq
```

## TODO and ideas

See [TODO.md](doc/TODO.md)

## Development and adding a new decoder

See [dev.md](doc/dev.md)

## Thanks and related projects

This project would not have been possible without [itchyny](https://github.com/itchyny)'s
jq implementation [gojq](https://github.com/itchyny/gojq). I also want to thank
[HexFiend](https://github.com/HexFiend/HexFiend) for inspiration and ideas and [stedolan](https://github.com/stedolan)
for inventing the [jq](https://github.com/stedolan/jq) language.

### Similar or related projects

- [HexFiend](https://github.com/HexFiend/HexFiend) Hex editor for macOS with format template support.
- [binspector](https://github.com/binspector/binspector) Binary format analysis tool with query langauge and REPL.
- [kaitai](https://kaitai.io) Declarative binary format parsing.
- [Wireshark](https://www.wireshark.org) Decodes network traffic (tip: `tshark -T json`).
- [MediaInfo](https://mediaarea.net/en/MediaInfo) Analyze media files  (tip `mediainfo --Output=JSON` and `mediainfo --Details=1`).
- [GNU poke](https://www.jemarch.net/poke) The extensible editor for structured binary data.
- [ffmpeg/ffprobe](https://ffmpeg.org) Powerful media libraries and tools.
- [hexdump](https://git.kernel.org/pub/scm/utils/util-linux/util-linux.git/tree/text-utils/hexdump.c) Hex viewer tool.
- [hex](https://git.janouch.name/p/hex) Interactive hex viewer with format support via lua.
- [hachoir](https://github.com/vstinner/hachoir) General python library for working binary data.
- [scapy](https://scapy.net) Decode/Encode formats, focus on network protocols.

## License

`fq` is distributed under the terms of the MIT License.

See the [LICENSE](LICENSE) file for license details.

Licenses of direct dependencies:

- Forked version of gojq https://github.com/itchyny/gojq/blob/main/LICENSE (MIT)
- Forked version of readline https://github.com/chzyer/readline/blob/master/LICENSE (MIT)
- github.com/BurntSushi/toml https://github.com/BurntSushi/toml/blob/master/COPYING (MIT)
- github.com/creasty/defaults https://github.com/creasty/defaults/blob/master/LICENSE (MIT)
- github.com/gomarkdown/markdown https://github.com/gomarkdown/markdown/blob/master/LICENSE.txt (BSD)
- github.com/gopacket/gopacket https://github.com/gopacket/gopacket/blob/master/LICENSE (BSD)
- github.com/mitchellh/copystructure https://github.com/mitchellh/copystructure/blob/master/LICENSE (MIT)
- github.com/mitchellh/mapstructure https://github.com/mitchellh/mapstructure/blob/master/LICENSE (MIT)
- github.com/pmezard/go-difflib https://github.com/pmezard/go-difflib/blob/master/LICENSE (BSD)
- golang/snappy https://github.com/golang/snappy/blob/master/LICENSE (BSD)
- golang/x/* https://github.com/golang/text/blob/master/LICENSE (BSD)
- gopkg.in/yaml.v3 https://github.com/go-yaml/yaml/blob/v3/LICENSE (MIT)
