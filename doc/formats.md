## Supported formats

[./formats_table.sh]: sh-start

|Name                    |Description                                                                     |Dependencies|
|-                       |-                                                                               |-|
|`aac_frame`             |Advanced&nbsp;Audio&nbsp;Coding&nbsp;frame                                      |<sub></sub>|
|`adts`                  |Audio&nbsp;Data&nbsp;Transport&nbsp;Stream                                      |<sub>`adts_frame`</sub>|
|`adts_frame`            |Audio&nbsp;Data&nbsp;Transport&nbsp;Stream&nbsp;frame                           |<sub>`aac_frame`</sub>|
|`amf0`                  |Action&nbsp;Message&nbsp;Format&nbsp;0                                          |<sub></sub>|
|`apev2`                 |APEv2&nbsp;metadata&nbsp;tag                                                    |<sub>`image`</sub>|
|`ar`                    |Unix&nbsp;archive                                                               |<sub>`probe`</sub>|
|[`asn1_ber`](#asn1_ber) |ASN1&nbsp;Basic&nbsp;Encoding&nbsp;Rules&nbsp;(also&nbsp;CER&nbsp;and&nbsp;DER) |<sub></sub>|
|`av1_ccr`               |AV1&nbsp;Codec&nbsp;Configuration&nbsp;Record                                   |<sub></sub>|
|`av1_frame`             |AV1&nbsp;frame                                                                  |<sub>`av1_obu`</sub>|
|`av1_obu`               |AV1&nbsp;Open&nbsp;Bitstream&nbsp;Unit                                          |<sub></sub>|
|`avc_annexb`            |H.264/AVC&nbsp;Annex&nbsp;B                                                     |<sub>`avc_nalu`</sub>|
|`avc_au`                |H.264/AVC&nbsp;Access&nbsp;Unit                                                 |<sub>`avc_nalu`</sub>|
|`avc_dcr`               |H.264/AVC&nbsp;Decoder&nbsp;Configuration&nbsp;Record                           |<sub>`avc_nalu`</sub>|
|`avc_nalu`              |H.264/AVC&nbsp;Network&nbsp;Access&nbsp;Layer&nbsp;Unit                         |<sub>`avc_sps` `avc_pps` `avc_sei`</sub>|
|`avc_pps`               |H.264/AVC&nbsp;Picture&nbsp;Parameter&nbsp;Set                                  |<sub></sub>|
|`avc_sei`               |H.264/AVC&nbsp;Supplemental&nbsp;Enhancement&nbsp;Information                   |<sub></sub>|
|`avc_sps`               |H.264/AVC&nbsp;Sequence&nbsp;Parameter&nbsp;Set                                 |<sub></sub>|
|[`avro_ocf`](#avro_ocf) |Avro&nbsp;object&nbsp;container&nbsp;file                                       |<sub></sub>|
|`bencode`               |BitTorrent&nbsp;bencoding                                                       |<sub></sub>|
|`bsd_loopback_frame`    |BSD&nbsp;loopback&nbsp;frame                                                    |<sub>`inet_packet`</sub>|
|[`bson`](#bson)         |Binary&nbsp;JSON                                                                |<sub></sub>|
|`bzip2`                 |bzip2&nbsp;compression                                                          |<sub>`probe`</sub>|
|[`cbor`](#cbor)         |Concise&nbsp;Binary&nbsp;Object&nbsp;Representation                             |<sub></sub>|
|`dns`                   |DNS&nbsp;packet                                                                 |<sub></sub>|
|`dns_tcp`               |DNS&nbsp;packet&nbsp;(TCP)                                                      |<sub></sub>|
|`elf`                   |Executable&nbsp;and&nbsp;Linkable&nbsp;Format                                   |<sub></sub>|
|`ether8023_frame`       |Ethernet&nbsp;802.3&nbsp;frame                                                  |<sub>`inet_packet`</sub>|
|`exif`                  |Exchangeable&nbsp;Image&nbsp;File&nbsp;Format                                   |<sub></sub>|
|`flac`                  |Free&nbsp;Lossless&nbsp;Audio&nbsp;Codec&nbsp;file                              |<sub>`flac_metadatablocks` `flac_frame`</sub>|
|`flac_frame`            |FLAC&nbsp;frame                                                                 |<sub></sub>|
|`flac_metadatablock`    |FLAC&nbsp;metadatablock                                                         |<sub>`flac_streaminfo` `flac_picture` `vorbis_comment`</sub>|
|`flac_metadatablocks`   |FLAC&nbsp;metadatablocks                                                        |<sub>`flac_metadatablock`</sub>|
|`flac_picture`          |FLAC&nbsp;metadatablock&nbsp;picture                                            |<sub>`image`</sub>|
|`flac_streaminfo`       |FLAC&nbsp;streaminfo                                                            |<sub></sub>|
|`gif`                   |Graphics&nbsp;Interchange&nbsp;Format                                           |<sub></sub>|
|`gzip`                  |gzip&nbsp;compression                                                           |<sub>`probe`</sub>|
|`hevc_annexb`           |H.265/HEVC&nbsp;Annex&nbsp;B                                                    |<sub>`hevc_nalu`</sub>|
|`hevc_au`               |H.265/HEVC&nbsp;Access&nbsp;Unit                                                |<sub>`hevc_nalu`</sub>|
|`hevc_dcr`              |H.265/HEVC&nbsp;Decoder&nbsp;Configuration&nbsp;Record                          |<sub>`hevc_nalu`</sub>|
|`hevc_nalu`             |H.265/HEVC&nbsp;Network&nbsp;Access&nbsp;Layer&nbsp;Unit                        |<sub>`hevc_vps` `hevc_pps` `hevc_sps`</sub>|
|`hevc_pps`              |H.265/HEVC&nbsp;Picture&nbsp;Parameter&nbsp;Set                                 |<sub></sub>|
|`hevc_sps`              |H.265/HEVC&nbsp;Sequence&nbsp;Parameter&nbsp;Set                                |<sub></sub>|
|`hevc_vps`              |H.265/HEVC&nbsp;Video&nbsp;Parameter&nbsp;Set                                   |<sub></sub>|
|`icc_profile`           |International&nbsp;Color&nbsp;Consortium&nbsp;profile                           |<sub></sub>|
|`icmp`                  |Internet&nbsp;Control&nbsp;Message&nbsp;Protocol                                |<sub></sub>|
|`icmpv6`                |Internet&nbsp;Control&nbsp;Message&nbsp;Protocol&nbsp;v6                        |<sub></sub>|
|`id3v1`                 |ID3v1&nbsp;metadata                                                             |<sub></sub>|
|`id3v11`                |ID3v1.1&nbsp;metadata                                                           |<sub></sub>|
|`id3v2`                 |ID3v2&nbsp;metadata                                                             |<sub>`image`</sub>|
|`ipv4_packet`           |Internet&nbsp;protocol&nbsp;v4&nbsp;packet                                      |<sub>`ip_packet`</sub>|
|`ipv6_packet`           |Internet&nbsp;protocol&nbsp;v6&nbsp;packet                                      |<sub>`ip_packet`</sub>|
|`jpeg`                  |Joint&nbsp;Photographic&nbsp;Experts&nbsp;Group&nbsp;file                       |<sub>`exif` `icc_profile`</sub>|
|`json`                  |JSON                                                                            |<sub></sub>|
|[`macho`](#macho)       |Mach-O&nbsp;macOS&nbsp;executable                                               |<sub></sub>|
|[`matroska`](#matroska) |Matroska&nbsp;file                                                              |<sub>`aac_frame` `av1_ccr` `av1_frame` `avc_au` `avc_dcr` `flac_frame` `flac_metadatablocks` `hevc_au` `hevc_dcr` `image` `mp3_frame` `mpeg_asc` `mpeg_pes_packet` `mpeg_spu` `opus_packet` `vorbis_packet` `vp8_frame` `vp9_cfm` `vp9_frame`</sub>|
|`mp3`                   |MP3&nbsp;file                                                                   |<sub>`id3v2` `id3v1` `id3v11` `apev2` `mp3_frame`</sub>|
|`mp3_frame`             |MPEG&nbsp;audio&nbsp;layer&nbsp;3&nbsp;frame                                    |<sub>`xing`</sub>|
|[`mp4`](#mp4)           |MPEG-4&nbsp;file&nbsp;and&nbsp;similar                                          |<sub>`aac_frame` `av1_ccr` `av1_frame` `flac_frame` `flac_metadatablocks` `id3v2` `image` `jpeg` `mp3_frame` `avc_au` `avc_dcr` `mpeg_es` `hevc_au` `hevc_dcr` `mpeg_pes_packet` `opus_packet` `protobuf_widevine` `pssh_playready` `vorbis_packet` `vp9_frame` `vpx_ccr` `icc_profile`</sub>|
|`mpeg_asc`              |MPEG-4&nbsp;Audio&nbsp;Specific&nbsp;Config                                     |<sub></sub>|
|`mpeg_es`               |MPEG&nbsp;Elementary&nbsp;Stream                                                |<sub>`mpeg_asc` `vorbis_packet`</sub>|
|`mpeg_pes`              |MPEG&nbsp;Packetized&nbsp;elementary&nbsp;stream                                |<sub>`mpeg_pes_packet` `mpeg_spu`</sub>|
|`mpeg_pes_packet`       |MPEG&nbsp;Packetized&nbsp;elementary&nbsp;stream&nbsp;packet                    |<sub></sub>|
|`mpeg_spu`              |Sub&nbsp;Picture&nbsp;Unit&nbsp;(DVD&nbsp;subtitle)                             |<sub></sub>|
|`mpeg_ts`               |MPEG&nbsp;Transport&nbsp;Stream                                                 |<sub></sub>|
|[`msgpack`](#msgpack)   |MessagePack                                                                     |<sub></sub>|
|`ogg`                   |OGG&nbsp;file                                                                   |<sub>`ogg_page` `vorbis_packet` `opus_packet` `flac_metadatablock` `flac_frame`</sub>|
|`ogg_page`              |OGG&nbsp;page                                                                   |<sub></sub>|
|`opus_packet`           |Opus&nbsp;packet                                                                |<sub>`vorbis_comment`</sub>|
|`pcap`                  |PCAP&nbsp;packet&nbsp;capture                                                   |<sub>`link_frame` `tcp_stream` `ipv4_packet`</sub>|
|`pcapng`                |PCAPNG&nbsp;packet&nbsp;capture                                                 |<sub>`link_frame` `tcp_stream` `ipv4_packet`</sub>|
|`png`                   |Portable&nbsp;Network&nbsp;Graphics&nbsp;file                                   |<sub>`icc_profile` `exif`</sub>|
|[`protobuf`](#protobuf) |Protobuf                                                                        |<sub></sub>|
|`protobuf_widevine`     |Widevine&nbsp;protobuf                                                          |<sub>`protobuf`</sub>|
|`pssh_playready`        |PlayReady&nbsp;PSSH                                                             |<sub></sub>|
|`raw`                   |Raw&nbsp;bits                                                                   |<sub></sub>|
|[`rtmp`](#rtmp)         |Real-Time&nbsp;Messaging&nbsp;Protocol                                          |<sub>`amf0`</sub>|
|`sll2_packet`           |Linux&nbsp;cooked&nbsp;capture&nbsp;encapsulation&nbsp;v2                       |<sub>`inet_packet`</sub>|
|`sll_packet`            |Linux&nbsp;cooked&nbsp;capture&nbsp;encapsulation                               |<sub>`inet_packet`</sub>|
|`tar`                   |Tar&nbsp;archive                                                                |<sub>`probe`</sub>|
|`tcp_segment`           |Transmission&nbsp;control&nbsp;protocol&nbsp;segment                            |<sub></sub>|
|`tiff`                  |Tag&nbsp;Image&nbsp;File&nbsp;Format                                            |<sub>`icc_profile`</sub>|
|`udp_datagram`          |User&nbsp;datagram&nbsp;protocol                                                |<sub>`udp_payload`</sub>|
|`vorbis_comment`        |Vorbis&nbsp;comment                                                             |<sub>`flac_picture`</sub>|
|`vorbis_packet`         |Vorbis&nbsp;packet                                                              |<sub>`vorbis_comment`</sub>|
|`vp8_frame`             |VP8&nbsp;frame                                                                  |<sub></sub>|
|`vp9_cfm`               |VP9&nbsp;Codec&nbsp;Feature&nbsp;Metadata                                       |<sub></sub>|
|`vp9_frame`             |VP9&nbsp;frame                                                                  |<sub></sub>|
|`vpx_ccr`               |VPX&nbsp;Codec&nbsp;Configuration&nbsp;Record                                   |<sub></sub>|
|`wav`                   |WAV&nbsp;file                                                                   |<sub>`id3v2` `id3v1` `id3v11`</sub>|
|`webp`                  |WebP&nbsp;image                                                                 |<sub>`vp8_frame`</sub>|
|`xing`                  |Xing&nbsp;header                                                                |<sub></sub>|
|`zip`                   |ZIP&nbsp;archive                                                                |<sub>`probe`</sub>|
|`image`                 |Group                                                                           |<sub>`gif` `jpeg` `mp4` `png` `tiff` `webp`</sub>|
|`inet_packet`           |Group                                                                           |<sub>`ipv4_packet` `ipv6_packet`</sub>|
|`ip_packet`             |Group                                                                           |<sub>`icmp` `icmpv6` `tcp_segment` `udp_datagram`</sub>|
|`link_frame`            |Group                                                                           |<sub>`bsd_loopback_frame` `ether8023_frame` `sll2_packet` `sll_packet`</sub>|
|`probe`                 |Group                                                                           |<sub>`adts` `ar` `avro_ocf` `bzip2` `elf` `flac` `gif` `gzip` `jpeg` `json` `macho` `matroska` `mp3` `mp4` `mpeg_ts` `ogg` `pcap` `pcapng` `png` `tar` `tiff` `wav` `webp` `zip`</sub>|
|`tcp_stream`            |Group                                                                           |<sub>`dns` `rtmp`</sub>|
|`udp_payload`           |Group                                                                           |<sub>`dns`</sub>|

[#]: sh-end

## Format options

Currently the only option is `force` and is used to ignore some format assertion errors. It can be used as a decode option or as a CLI `-o` option:

```
fq -d mp4 -o force=true file.mp4
fq -d raw 'mp4({force: true})' file.mp4
```

## Format details

[./formats_collect.sh]: sh-start

### asn1_ber

Supports decoding BER, CER and DER ([X.690]([X.690_1297.pdf)).

- Currently no extra validation is done for CER and DER.
- Does not support specifying a schema.
- Supports `torepr` but without schema all sequences and sets will be arrays.

```
fq -d asn1_ber torepr file.ber
```

Functions `frompem` and `topem` can help working with PEM format:

```
fq -d raw 'frompem | asn1_ber | d' cert.pem
```

If the schema is known and not that complicated it can be reproduced:

```
fq -d asn1_ber 'torepr as $r | ["version", "modulus", "private_exponent", "private_exponen", "prime1", "prime2", "exponent1", "exponent2", "coefficient"] | with_entries({key: .value, value: $r[.key]})' pkcs1.der
```

Can be used to decode nested parts:

```
fq -d asn1_ber '.constructed[1].value | asn1_ber' file.ber
```

References and tools:
- https://www.itu.int/ITU-T/studygroups/com10/languages/X.690_1297.pdf
- https://en.wikipedia.org/wiki/X.690
- https://letsencrypt.org/docs/a-warm-welcome-to-asn1-and-der/
- https://lapo.it/asn1js/

### avro_ocf

Supports reading Avro Object Container Format (OCF) files based on the [1.11.0 specification](https://avro.apache.org/docs/current/spec.html#Object+Container+Files).

Capable of handling null, deflate, and snappy codecs for data compression.

Limitations:
 - Schema does not support self-referential types, only built-in types.
 - Decimal logical types are not supported for decoding, will just be treated as their primitive type
### becode

Supports `torepr`:

```
fq -d bencode torepr file.torrent
```

### bson

Supports `torepr`:

```
fq -d bson torepr file.bson
```

### cbor

Supports `torepr`:

```
fq -d cbor torepr file.cbor
fq -d cbor 'torepr.field' file.cbor
fq -d cbor 'torepr | .field' file.cbor
fq -d cbor 'torepr | grep("abc")' file.cbor
```

### macho

Supports decoding vanilla and FAT Mach-O binaries.

#### Examples

To decode the macOS build of `fq`:

```
fq . /path/to/fq
```

```
fq '.load_commands[] | select(.cmd=="segment_64")' /path/to/fq
```

Note you can use `-d macho` to decode a broken Mach-O binary.

#### References:
- https://github.com/aidansteele/osx-abi-macho-file-format-reference

### matroska

Supports `matroska_path`:

```
$ fq 'matroska_path(".Segment.Tracks[0]")' file.mkv
     │00 01 02 03 04 05 06 07 08 09│0123456789│.elements[1].elements[3]{}:
0x122│         16 54 ae 6b         │   .T.k   │  id: "Tracks" (0x1654ae6b) (A Top-Level Element of information with many tracks described.)
     │                             │          │  type: "master" (7)
0x122│                     4d bf   │       M. │  size: 3519
0x122│                           bf│         .│  elements[0:3]:
0x12c│84 cf 8b db a0 ae 01 00 00 00│..........│
0x136│00 00 00 78 d7 81 01 73 c5 88│...x...s..│
*    │until 0xee9.7 (3519)         │          │
```

```
$ fq 'first(grep_by(.id == "Tracks")) | matroska_path' test.mkv
".Segment.Tracks"
```

### mp4

Supports `mp4_path`:

```
$ fq 'mp4_path(".moov.trak[1]")' file.mp4
     │00 01 02 03 04 05 06 07 08 09│0123456789│.boxes[3].boxes[1]{}:
0x4f6│                     00 00 02│       ...│  size: 573
0x500│3d                           │=         │
0x500│   74 72 61 6b               │ trak     │  type: "trak" (Container for an individual track or stream)
0x500│               00 00 00 5c 74│     ...\t│  boxes[0:3]:
0x50a│6b 68 64 00 00 00 03 00 00 00│khd.......│
0x514│00 00 00 00 00 00 00 00 01 00│..........│
*    │until 0x739.7 (565)          │          │
```

```
$ fq 'first(grep_by(.type == "trak")) | mp4_path' file.mp4
".moov.trak"
```

### msgpack

Supports `torepr`:

```
fq -d msgpack torepr file.msgpack
```

### protobuf

`protobuf` decoder can be used to decode sub messages:

```
fq -d protobuf '.fields[6].wire_value | protobuf | d'
```

### rtmp

Current only supports plain RTMP (not RTMPT or encrypted variants etc) with AMF0 (not AMF3).

[#]: sh-end


## Dependency graph

![alt text](formats.svg "Format diagram")
