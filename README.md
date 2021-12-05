# fq

Tool, language and decoders for inspecting binary data.

![fq demo](doc/demo.svg)

In most cases fq works the same way as jq but instead of reading JSON it reads binary data.
The result is a JSON compatbile structures where each value has a bit range, symbolic
interpretations and know how to be presented in a useful way.

**NOTE:** fq is early in development and many things are missing, broken or do not make sense.
That also means there is a great opportunity to help out.

## Goals

- Make binary formats accessible and queryable.
- Nested formats and bit-oriented decoding.
- Quick and comfortable CLI tool.
- Bit and byte transformations and conversions.
- Programmer's calculator.

## Usage

Basic usage is: `fq . file`.

For details see [usage.md](doc/usage.md)

## Install

Download archive from [releases](https://github.com/wader/fq/releases) page for your
platform, unarchive it and move the executable to `PATH` etc.

### Homebrew

```sh
# install latest release
brew install wader/tap/fq
```

### Build from source

Make sure you have go 1.17 or later installed.

To install directly from git repository do:
```sh
# build and install latest release
go install github.com/wader/fq@latest

# or build and install latest master
go install github.com/wader/fq@master

# copy binary to $PATH if needed
cp "$(go env GOPATH)/bin/fq" /usr/local/bin
```

To build and run tests from source directory:
```sh
make test fq
# copy binary to $PATH if needed
cp fq /usr/local/bin
```

## Supported formats

[./formats_list.jq]: sh-start

aac_frame, adts, adts_frame, apev2, av1_ccr, av1_frame, av1_obu, avc_annexb, avc_au, avc_dcr, avc_nalu, avc_pps, avc_sei, avc_sps, bzip2, dns, dns_tcp, elf, ether8023_frame, exif, flac, flac_frame, flac_metadatablock, flac_metadatablocks, flac_picture, flac_streaminfo, gif, gzip, hevc_annexb, hevc_au, hevc_dcr, hevc_nalu, icc_profile, icmp, id3v1, id3v11, id3v2, ipv4_packet, jpeg, json, matroska, mp3, mp3_frame, mp4, mpeg_asc, mpeg_es, mpeg_pes, mpeg_pes_packet, mpeg_spu, mpeg_ts, ogg, ogg_page, opus_packet, pcap, pcapng, png, protobuf, protobuf_widevine, pssh_playready, raw, sll2_packet, sll_packet, tar, tcp_segment, tiff, udp_datagram, vorbis_comment, vorbis_packet, vp8_frame, vp9_cfm, vp9_frame, vpx_ccr, wav, webp, xing, zip

[#]: sh-end

For details see [formats.md](doc/formats.md)

## TODO and ideas

See [TODO.md](doc/TODO.md)

## Development

See [dev.md](doc/dev.md)

## Thanks and related projects

This project would not have been possible without [itchyny](https://github.com/itchyny)'s
jq implementation [gojq](https://github.com/itchyny/gojq). I also want to thank
[HexFiend](https://github.com/HexFiend/HexFiend) for inspiration and ideas and [stedolan](https://github.com/stedolan)
for inventing the [jq](https://github.com/stedolan/jq) language.

Similar projects:
- https://github.com/HexFiend/HexFiend
- https://github.com/binspector/binspector
- https://kaitai.io
