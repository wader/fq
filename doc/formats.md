## Format specific functions

[./formats_collect.sh]: sh-start

### bencode

Supports `torepr`, ex:

```
fq -d bencode torepr file.torrent
```
### bson

Supports `torepr`, ex:

```
fq -d bson torepr file.bson
```
### cbor

Supports `torepr`, ex:

```
fq -d cbor torepr file.cbor
fq -d cbor 'torepr.field' file.cbor
fq -d cbor 'torepr | .field' file.cbor
fq -d cbor 'torepr | grep("abc")' file.cbor
```
### matroska

Supports `matroska_path`, ex:

```
$ fq 'matroska_path(".Segment.Tracks[0]") file.mkv
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
$ fq 'first(grep_by(.id == "Tracks")) | matroska_path' testfiles/test.mkv
".Segment.Tracks"
```
### mp4

Supports `mp4_path`, ex:

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

Supports `torepr`, ex:

```
fq -d msgpack torepr file.msgpack
```

[#]: sh-end

## Supported formats

[./formats_table.jq]: sh-start

|Name                  |Description                                                   |Dependencies|
|-                     |-                                                             |-|
|`aac_frame`           |Advanced&nbsp;Audio&nbsp;Coding&nbsp;frame                    |<sub></sub>|
|`adts`                |Audio&nbsp;Data&nbsp;Transport&nbsp;Stream                    |<sub>`adts_frame`</sub>|
|`adts_frame`          |Audio&nbsp;Data&nbsp;Transport&nbsp;Stream&nbsp;frame         |<sub>`aac_frame`</sub>|
|`apev2`               |APEv2&nbsp;metadata&nbsp;tag                                  |<sub>`image`</sub>|
|`ar`                  |Unix&nbsp;archive                                             |<sub>`probe`</sub>|
|`av1_ccr`             |AV1&nbsp;Codec&nbsp;Configuration&nbsp;Record                 |<sub></sub>|
|`av1_frame`           |AV1&nbsp;frame                                                |<sub>`av1_obu`</sub>|
|`av1_obu`             |AV1&nbsp;Open&nbsp;Bitstream&nbsp;Unit                        |<sub></sub>|
|`avc_annexb`          |H.264/AVC&nbsp;Annex&nbsp;B                                   |<sub>`avc_nalu`</sub>|
|`avc_au`              |H.264/AVC&nbsp;Access&nbsp;Unit                               |<sub>`avc_nalu`</sub>|
|`avc_dcr`             |H.264/AVC&nbsp;Decoder&nbsp;Configuration&nbsp;Record         |<sub>`avc_nalu`</sub>|
|`avc_nalu`            |H.264/AVC&nbsp;Network&nbsp;Access&nbsp;Layer&nbsp;Unit       |<sub>`avc_sps` `avc_pps` `avc_sei`</sub>|
|`avc_pps`             |H.264/AVC&nbsp;Picture&nbsp;Parameter&nbsp;Set                |<sub></sub>|
|`avc_sei`             |H.264/AVC&nbsp;Supplemental&nbsp;Enhancement&nbsp;Information |<sub></sub>|
|`avc_sps`             |H.264/AVC&nbsp;Sequence&nbsp;Parameter&nbsp;Set               |<sub></sub>|
|`bencode`             |BitTorrent&nbsp;bencoding                                     |<sub></sub>|
|`bsd_loopback_frame`  |BSD&nbsp;loopback&nbsp;frame                                  |<sub>`ipv4_packet`</sub>|
|`bson`                |Binary&nbsp;JSON                                              |<sub></sub>|
|`bzip2`               |bzip2&nbsp;compression                                        |<sub>`probe`</sub>|
|`cbor`                |Concise&nbsp;Binary&nbsp;Object&nbsp;Representation           |<sub></sub>|
|`dns`                 |DNS&nbsp;packet                                               |<sub></sub>|
|`dns_tcp`             |DNS&nbsp;packet&nbsp;(TCP)                                    |<sub></sub>|
|`elf`                 |Executable&nbsp;and&nbsp;Linkable&nbsp;Format                 |<sub></sub>|
|`ether8023_frame`     |Ethernet&nbsp;802.3&nbsp;frame                                |<sub>`ipv4_packet`</sub>|
|`exif`                |Exchangeable&nbsp;Image&nbsp;File&nbsp;Format                 |<sub></sub>|
|`flac`                |Free&nbsp;Lossless&nbsp;Audio&nbsp;Codec&nbsp;file            |<sub>`flac_metadatablocks` `flac_frame`</sub>|
|`flac_frame`          |FLAC&nbsp;frame                                               |<sub></sub>|
|`flac_metadatablock`  |FLAC&nbsp;metadatablock                                       |<sub>`flac_streaminfo` `flac_picture` `vorbis_comment`</sub>|
|`flac_metadatablocks` |FLAC&nbsp;metadatablocks                                      |<sub>`flac_metadatablock`</sub>|
|`flac_picture`        |FLAC&nbsp;metadatablock&nbsp;picture                          |<sub>`image`</sub>|
|`flac_streaminfo`     |FLAC&nbsp;streaminfo                                          |<sub></sub>|
|`gif`                 |Graphics&nbsp;Interchange&nbsp;Format                         |<sub></sub>|
|`gzip`                |gzip&nbsp;compression                                         |<sub>`probe`</sub>|
|`hevc_annexb`         |H.265/HEVC&nbsp;Annex&nbsp;B                                  |<sub>`hevc_nalu`</sub>|
|`hevc_au`             |H.265/HEVC&nbsp;Access&nbsp;Unit                              |<sub>`hevc_nalu`</sub>|
|`hevc_dcr`            |H.265/HEVC&nbsp;Decoder&nbsp;Configuration&nbsp;Record        |<sub>`hevc_nalu`</sub>|
|`hevc_nalu`           |H.265/HEVC&nbsp;Network&nbsp;Access&nbsp;Layer&nbsp;Unit      |<sub></sub>|
|`icc_profile`         |International&nbsp;Color&nbsp;Consortium&nbsp;profile         |<sub></sub>|
|`icmp`                |Internet&nbsp;Control&nbsp;Message&nbsp;Protocol              |<sub></sub>|
|`id3v1`               |ID3v1&nbsp;metadata                                           |<sub></sub>|
|`id3v11`              |ID3v1.1&nbsp;metadata                                         |<sub></sub>|
|`id3v2`               |ID3v2&nbsp;metadata                                           |<sub>`image`</sub>|
|`ipv4_packet`         |Internet&nbsp;protocol&nbsp;v4&nbsp;packet                    |<sub>`udp_datagram` `tcp_segment` `icmp`</sub>|
|`jpeg`                |Joint&nbsp;Photographic&nbsp;Experts&nbsp;Group&nbsp;file     |<sub>`exif` `icc_profile`</sub>|
|`json`                |JSON                                                          |<sub></sub>|
|`matroska`            |Matroska&nbsp;file                                            |<sub>`aac_frame` `av1_ccr` `av1_frame` `avc_au` `avc_dcr` `flac_frame` `flac_metadatablocks` `hevc_au` `hevc_dcr` `image` `mp3_frame` `mpeg_asc` `mpeg_pes_packet` `mpeg_spu` `opus_packet` `vorbis_packet` `vp8_frame` `vp9_cfm` `vp9_frame`</sub>|
|`mp3`                 |MP3&nbsp;file                                                 |<sub>`id3v2` `id3v1` `id3v11` `apev2` `mp3_frame`</sub>|
|`mp3_frame`           |MPEG&nbsp;audio&nbsp;layer&nbsp;3&nbsp;frame                  |<sub>`xing`</sub>|
|`mp4`                 |MPEG-4&nbsp;file&nbsp;and&nbsp;similar                        |<sub>`aac_frame` `av1_ccr` `av1_frame` `flac_frame` `flac_metadatablocks` `id3v2` `image` `jpeg` `mp3_frame` `avc_au` `avc_dcr` `mpeg_es` `hevc_au` `hevc_dcr` `mpeg_pes_packet` `opus_packet` `protobuf_widevine` `pssh_playready` `vorbis_packet` `vp9_frame` `vpx_ccr`</sub>|
|`mpeg_asc`            |MPEG-4&nbsp;Audio&nbsp;Specific&nbsp;Config                   |<sub></sub>|
|`mpeg_es`             |MPEG&nbsp;Elementary&nbsp;Stream                              |<sub>`mpeg_asc` `vorbis_packet`</sub>|
|`mpeg_pes`            |MPEG&nbsp;Packetized&nbsp;elementary&nbsp;stream              |<sub>`mpeg_pes_packet` `mpeg_spu`</sub>|
|`mpeg_pes_packet`     |MPEG&nbsp;Packetized&nbsp;elementary&nbsp;stream&nbsp;packet  |<sub></sub>|
|`mpeg_spu`            |Sub&nbsp;Picture&nbsp;Unit&nbsp;(DVD&nbsp;subtitle)           |<sub></sub>|
|`mpeg_ts`             |MPEG&nbsp;Transport&nbsp;Stream                               |<sub></sub>|
|`msgpack`             |MessagePack                                                   |<sub></sub>|
|`ogg`                 |OGG&nbsp;file                                                 |<sub>`ogg_page` `vorbis_packet` `opus_packet` `flac_metadatablock` `flac_frame`</sub>|
|`ogg_page`            |OGG&nbsp;page                                                 |<sub></sub>|
|`opus_packet`         |Opus&nbsp;packet                                              |<sub>`vorbis_comment`</sub>|
|`pcap`                |PCAP&nbsp;packet&nbsp;capture                                 |<sub>`link_frame` `tcp_stream` `ipv4_packet`</sub>|
|`pcapng`              |PCAPNG&nbsp;packet&nbsp;capture                               |<sub>`link_frame` `tcp_stream` `ipv4_packet`</sub>|
|`png`                 |Portable&nbsp;Network&nbsp;Graphics&nbsp;file                 |<sub>`icc_profile` `exif`</sub>|
|`protobuf`            |Protobuf                                                      |<sub></sub>|
|`protobuf_widevine`   |Widevine&nbsp;protobuf                                        |<sub>`protobuf`</sub>|
|`pssh_playready`      |PlayReady&nbsp;PSSH                                           |<sub></sub>|
|`raw`                 |Raw&nbsp;bits                                                 |<sub></sub>|
|`sll2_packet`         |Linux&nbsp;cooked&nbsp;capture&nbsp;encapsulation&nbsp;v2     |<sub>`ether8023_frame`</sub>|
|`sll_packet`          |Linux&nbsp;cooked&nbsp;capture&nbsp;encapsulation             |<sub>`ether8023_frame`</sub>|
|`tar`                 |Tar&nbsp;archive                                              |<sub>`probe`</sub>|
|`tcp_segment`         |Transmission&nbsp;control&nbsp;protocol&nbsp;segment          |<sub></sub>|
|`tiff`                |Tag&nbsp;Image&nbsp;File&nbsp;Format                          |<sub>`icc_profile`</sub>|
|`udp_datagram`        |User&nbsp;datagram&nbsp;protocol                              |<sub>`udp_payload`</sub>|
|`vorbis_comment`      |Vorbis&nbsp;comment                                           |<sub>`flac_picture`</sub>|
|`vorbis_packet`       |Vorbis&nbsp;packet                                            |<sub>`vorbis_comment`</sub>|
|`vp8_frame`           |VP8&nbsp;frame                                                |<sub></sub>|
|`vp9_cfm`             |VP9&nbsp;Codec&nbsp;Feature&nbsp;Metadata                     |<sub></sub>|
|`vp9_frame`           |VP9&nbsp;frame                                                |<sub></sub>|
|`vpx_ccr`             |VPX&nbsp;Codec&nbsp;Configuration&nbsp;Record                 |<sub></sub>|
|`wav`                 |WAV&nbsp;file                                                 |<sub>`id3v2` `id3v1` `id3v11`</sub>|
|`webp`                |WebP&nbsp;image                                               |<sub>`vp8_frame`</sub>|
|`xing`                |Xing&nbsp;header                                              |<sub></sub>|
|`zip`                 |ZIP&nbsp;archive                                              |<sub>`probe`</sub>|
|`image`               |Group                                                         |<sub>`gif` `jpeg` `mp4` `png` `tiff` `webp`</sub>|
|`link_frame`          |Group                                                         |<sub>`bsd_loopback_frame` `ether8023_frame` `sll2_packet` `sll_packet`</sub>|
|`probe`               |Group                                                         |<sub>`adts` `ar` `bzip2` `elf` `flac` `gif` `gzip` `jpeg` `json` `matroska` `mp3` `mp4` `mpeg_ts` `ogg` `pcap` `pcapng` `png` `tar` `tiff` `wav` `webp` `zip`</sub>|
|`tcp_stream`          |Group                                                         |<sub>`dns`</sub>|
|`udp_payload`         |Group                                                         |<sub>`dns`</sub>|

[#]: sh-end

![alt text](formats.svg "Format diagram")
