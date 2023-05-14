# 0.5.0

Mostly a bug fix release but adds `-V` for easy JSON output.

## Changes

- Add `-V` argument to default output JSON instead of decode tree in case of decode value. #385 Thanks @peterwaller-arm for reminding me to merge this.
  ```sh
  # default in case of decode value is to show a hexdump tree
  $ fq '.headers | grep_by(.id=="TSSE").text' file.mp3
      │00 01 02 03 04 05 06 07 08 09 0a 0b│0123456789ab│
  0x0c│                           4c 61 76│         Lav│.headers[0].frames[0].text: "Lavf58.76.100"
  0x18│66 35 38 2e 37 36 2e 31 30 30 00   │f58.76.100. │

  # with -V an implicit "tovalue" is done
  $ fq -V '.headers | grep_by(.id=="TSSE").text' file.mp3
  "Lavf58.76.100"

  # and in combination with -r will for strings output a "raw string" without quotes
  # for other types like number, object, array etc -r makes not difference (same as jq)
  $ fq -Vr '.headers | grep_by(.id=="TSSE").text' file.mp3
  Lavf58.76.100
  ```

  As a side note `-V` can be used with binary type also. Then the binary data will be interpreted as UTF-8 and turned into a string.
  ```sh
  # trailing null terminator ends up as codepoint zero `\u0000`
  $ fq -V '.headers | grep_by(.id=="TSSE").text | tobytes' file.mp3
  "Lavf58.76.100\u0000"

  # with -r null terminator and a new line is outputted
  $ fq -Vr '.headers | grep_by(.id=="TSSE").text | tobytes' file.mp3 | hexdump -C
  00000000  4c 61 76 66 35 38 2e 37  36 2e 31 30 30 00 0a     |Lavf58.76.100..|
  0000000f

  # in contrast raw binary output has no new line separator
  $ fq '.headers | grep_by(.id=="TSSE").text | tobytes' doc/file.mp3 | hexdump -C
  00000000  4c 61 76 66 35 38 2e 37  36 2e 31 30 30 00        |Lavf58.76.100.|
  0000000e
  ```
- Fix issue using decode value in object passed as argument to internal function. #638
  ```sh
  # this used to fail but now works
  fq '.tracks[0].samples[10] | avc_au({length_size: <decode value>})' file.mp4
  ```
- Some typo fixes. Thanks @retokromer and @peterwaller-arm

## Decoder changes

- `aiff` Basic AIFF decoder added. #614
- `matroska` Update to latest specification. #640
- `msgpack` Fix bug decoding some fixstr lengths. #636 Thanks @schmee for reporting.

## Changelog

* 4ad1cced Update docker-golang to 1.20.3 from 1.20.2
* f7dca477 Update github-go-version to 1.20.3 from 1.20.2
* c9608939 Update github-golangci-lint to 1.52.0 from 1.51.2
* 0a6b46c8 Update github-golangci-lint to 1.52.1 from 1.52.0
* c4eb67d9 Update github-golangci-lint to 1.52.2 from 1.52.1
* 19140a6f Update gomod-creasty-defaults to 1.7.0 from 1.6.0
* 6e5df724 Update gomod-golang-x-crypto to 0.8.0 from 0.7.0
* 6c4aebfe Update gomod-golang-x-net to 0.9.0 from 0.8.0
* f13cc979 Update gomod-golang/text to 0.9.0 from 0.8.0
* e2af57ee Update gomod-gopacket to 1.1.0 from 1.0.0
* a63fd684 Update make-golangci-lint to 1.52.0 from 1.51.2
* d3d1f0e8 Update make-golangci-lint to 1.52.1 from 1.52.0
* f0b08457 Update make-golangci-lint to 1.52.2 from 1.52.1
* dc4a82ee aiff: Add basic decoder
* c5f6809b decode,fuzz,dev: Move recoverable error check to recoverfn.Run
* 980ecdba decode: Add float 80 reader
* a6c4db75 decode: Cleanup old unused help system code
* 87e5bb14 fix typo
* 0b6ef2a9 golangci-lint: Disable revive unused-parameter and update for new default config
* 427ce78d interp: Add --value-output/-V option to do tovalue before output
* 9a1ef84c interp: Allow and convert JQValues:s (ex decode value) in function arg objects
* 3dd2c61d interp: Fix input completion regression in sub-REPLs
* 5415bfca interp: Make completion work again
* 2a2b64dd matroska: Update ebml specification
* 82da99c9 msgpack: Add str, array and object type tests
* 97360d6f msgpack: fixstr length field is 5 bits
* ffc66db0 readline: remove direct access to (*Instance).Config
* e1b02312 wav: Cleanup avi leftovers

# 0.4.0

TLS decode and decryption, better streaming matroska/webm support, support raw IP in PCAP and bug fixes.

## Changes

- Fix panic when interrupting big JSON output. #573
- Support passing options (`-o name=value`) to nested decoders. #589
  - Allows for example to pass keylog to a TLS decoder inside a PCAP file or to tell a container decoders to not decode samples inside a ZIP file etc.
- Exit with error if `-o name=@path` fails to read file at `path`. #597

## Decoder changes

- `id3v2` Properly decode CTOC subframes. #606
- `matroska`
  - Now supports streaming matroska and webm better (master elements with unknown size). #576 #581
  - Add `decode_samples` option. #574
  - Spec update and clean up of symbols and descriptions. #580
- `pcap,pcapng` Support raw IPv4 and IPv6 link frames. #599 #590
- `tls` Add Transport layer security decoder and decryption. #603
  - Supports TLS 1.0, 1.1, 1.2 and some SSL 3.0.
  - Decodes records and most messages and extensions.
  - Can decrypt most common cipher suites if a keylog is provided. See documentation for list of supported ciphers suites.
  ```sh
  # show first 50 bytes of decrypted client/server TLS application data stream
  # -o keylog=@file.pcap.keylog is used to read keylog from a file
  # first .stream is TCP stream, second .stream the application data stream
  $ fq -o keylog=@file.pcap.keylog '.tcp_connections[0].["client", "server"].stream.stream | tobytes[0:50] | dd' file.pcap
      │00 01 02 03 04 05 06 07 08 09 0a 0b│0123456789ab│
  0x00│47 45 54 20 2f 64 75 6d 70 2f 6c 6f│GET /dump/lo│.: raw bits 0x0-0x31.7 (50)
  0x0c│67 20 48 54 54 50 2f 31 2e 31 0d 0a│g HTTP/1.1..│
  0x18│48 6f 73 74 3a 20 69 6e 77 61 64 65│Host: inwade│
  0x24│72 2e 63 6f 6d 0d 0a 55 73 65 72 2d│r.com..User-│
  0x30│41 67                              │Ag          │
      │00 01 02 03 04 05 06 07 08 09 0a 0b│0123456789ab│
  0x00│48 54 54 50 2f 31 2e 31 20 32 30 30│HTTP/1.1 200│.: raw bits 0x0-0x31.7 (50)
  0x0c│20 4f 4b 0d 0a 41 63 63 65 70 74 2d│ OK..Accept-│
  0x18│52 61 6e 67 65 73 3a 20 62 79 74 65│Ranges: byte│
  0x24│73 0d 0a 43 6f 6e 74 65 6e 74 2d 4c│s..Content-L│
  0x30│65 6e                              │en          │

  # show first TLS record from server
  $ fq '.tcp_connections[0].server.stream.records[0] | d' file.pcap
      │00 01 02 03 04 05 06 07 08 09 0a 0b│0123456789ab│.tcp_connections[1].server.stream.records[0]{}: record
  0x00│16                                 │.           │  type: "handshake" (22) (valid)
  0x00│   03 03                           │ ..         │  version: "tls1.2" (0x303) (valid)
  0x00│         00 40                     │   .@       │  length: 64
      │                                   │            │  message{}:
  0x00│               02                  │     .      │    type: "server_hello" (2)
  0x00│                  00 00 3c         │      ..<   │    length: 60
  0x00│                           03 03   │         .. │    version: "tls1.2" (0x303)
      │                                   │            │    random{}:
  0x00│                                 86│           .│      gmt_unix_time: 2249760024 (2041-04-16T21:20:24Z)
  0x0c│18 9d 18                           │...         │
  0x0c│         19 92 33 c2 21 ce 4f 97 30│   ..3.!.O.0│      random_bytes: raw bits
  0x18│28 98 b3 fd 1e 15 f4 36 bb e9 14 f4│(......6....│
  0x24│67 61 66 79 d5 3f 06               │gafy.?.     │
  0x24│                     00            │       .    │    session_id_length: 0
      │                                   │            │    session_id: raw bits
  0x24│                        c0 2f      │        ./  │    cipher_suit: "TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256" (0xc02f)
  0x24│                              00   │          . │    compression_method: "null" (0x0)
  0x24│                                 00│           .│    extensions_length: 20
  0x30│14                                 │.           │
      │                                   │            │    extensions[0:2]:
      │                                   │            │      [0]{}: extension
  0x30│   ff 01                           │ ..         │        type: "renegotiation_info" (65281)
  0x30│         00 01                     │   ..       │        length: 1
  0x30│               00                  │     .      │        data: raw bits
      │                                   │            │      [1]{}: extension
  0x30│                  00 10            │      ..    │        type: "application_layer_protocol_negotiation" (16)
  0x30│                        00 0b      │        ..  │        length: 11
  0x30│                              00 09│          ..│        serer_names_length: 9
      │                                   │            │        protocols[0:1]:
      │                                   │            │          [0]{}: protocol
  0x3c│08                                 │.           │            length: 8
  0x3c│   68 74 74 70 2f 31 2e 31         │ http/1.1   │            name: "http/1.1"

  # use ja3.jq to calculate ja3 TLS fingerprint
  # https://github.com/wader/fq/blob/master/format/tls/testdata/ja3.jq
  $ fq -L path/to/ja3 'include "ja3"; pcap_ja3' file.pcap
  [
    {
      "client_ip": "192.168.1.193",
      "client_port": 64126,
      "ja3": "771,4866-4867-4865-49196-49200-159-52393-52392-52394-49195-49199-158-49188-49192-107-49187-49191-103-49162-49172-57-49161-49171-51-157-156-61-60-53-47-255,0-11-10-16-22-23-49-13-43-45-51-21,29-23-30-25-24,0-1-2",
      "ja3_digest": "bc29aa426fc99c0be1b9be941869f88a",
      "server_ip": "46.101.135.150",
      "server_port": 443
    }
  ]
   ```
- `toml` Fail faster to speed up probe. Could in some cases read the whole file before failing. Thanks @0-wiz-0 for report. #594
- `zip` Properly decode EOCD record in zip64 files. Thanks @0-wiz-0 for report and spec interpretation. #586 #596
- `xml` Fail faster to speed up probe. Could in some cases read the whole file before failing. Thanks @0-wiz-0 for report. #594

## Changelog

* 0581ecea Update docker-golang to 1.20.1 from 1.20.0
* 72870a5a Update docker-golang to 1.20.2 from 1.20.1
* 02e573a9 Update github-go-version to 1.20.1 from 1.20.0, 1.20.0, 1.20.0
* c5130887 Update github-go-version to 1.20.2 from 1.20.1
* ce263726 Update github-golangci-lint to 1.51.1 from 1.51.0
* 75bfdda3 Update github-golangci-lint to 1.51.2 from 1.51.1
* b1d9306b Update gomod-golang-x-crypto to 0.6.0 from 0.5.0
* c03d3ccd Update gomod-golang-x-crypto to 0.7.0 from 0.6.0
* 2430fba7 Update gomod-golang-x-net to 0.6.0 from 0.5.0
* dd8ab799 Update gomod-golang-x-net to 0.7.0 from 0.6.0
* 80a07446 Update gomod-golang-x-net to 0.8.0 from 0.7.0
* 97643b98 Update gomod-golang/text to 0.7.0 from 0.6.0
* e7168b99 Update gomod-golang/text to 0.8.0 from 0.7.0
* 36df57eb Update make-golangci-lint to 1.51.1 from 1.51.0
* 70e08faa Update make-golangci-lint to 1.51.2 from 1.51.1
* 50d26ec7 colorjson: Handle encoding error value
* 5c8e1151 colorjson: Refactor to option struct
* 8e0dde03 decode: Support multiple format args and some rename and refactor
* a1bb630a doc,fq: Improve cli help and some cleanup
* 156aeeca doc: Add FOSDEM 2023 talk
* 3e0ebafa doc: Run make doc
* 3cc83837 gojq: Update fq fork
* dec433fc help,markdown: Fix double line breaks when converting to text
* c75a83c8 help: Show default option value as JSON
* cc52a441 id3v2: Decode subframes for CTOC and add struct for headers
* dc79a73b interp,json: Move error handling to colorjson
* 73db6587 interp: Exit with error if -o name=@path fails to be read, also document
* c8666eeb ipv4_packet,ipv6_packet,sll_packet,sll2_packet: Support ipv4/ipv6 link frames and pass correct in arg
* b60aceca matroska: Add decode_samples option
* 9aaf2ddf matroska: Add unknown size test and add description to ebml header
* a8d0bf4d matroska: Assume master with unknown size has ended if a valid parent is found
* 0d14d7b4 matroska: Handle unknown size for non-master types a bit better
* c890a289 matroska: Update spec and make refs in descriptions look nicer
* 6c032455 pcap,pcapng,ipv4,ipv6: Support raw link type (ipv4 or ipv6)
* d4ea6632 pcap: Add ipv4 fragments tcp test
* f50bd6ee readline: Update fq fork
* 9852f56b tls: Add TLS 1.0, 1.1, 1.2 decode and decryption
* 56edb59e toml,xml: Fail fast on invalid content
* 5228fdd6 zip: Correctly look for and decode both zip32/64 EOCD record
* bdd6718d zip: Correctly peek for zip64 EOCD


# 0.3.0

Bug fix release, no new features mostly due to holidays and busy with other things (some jq related!).

Also been preparing for a [talk about fq](https://fosdem.org/2023/schedule/event/bintools_fq/) this weekend at [FOSDEM 2023](https://fosdem.org/2023/).

## Changes

* TCP reassembly is now less strict about invalid TCP options. Turns out some options might end up wrong in packet captures due to hardware acceleration etc. For example it seems to be common that TCP segments end up larger than configured connection MSS. Now PCAP:s with those kinds of TCP segments should be reassembled correctly.
* REPL now handles the del key properly. Before it could in some cases cause the output to be ignored.

## Decoder changes

- `mp3` Add option for max unknown bits to handle more mis-probing. Default to 50%
- `mp4`
  - `ftyp` set minor description to date for "qt" files
  - `tkhd` decode enabled, preview, etc flags
  - `udta` Handle case with box type is key and value rest of box
  - `sgpd`,`sbgp` Change grouping type to a string as it seems to be what it is in practice.
- `tcp_segment` Decode all standard options, MSS, Window scale, timestamp etc. Rename "maxseg" to "mss".

## Changelog

* 8702e1d1 Update docker-golang to 1.19.5 from 1.19.4
* a7f37d73 Update docker-golang to 1.20.0 from 1.19.5
* 826d9a52 Update github-go-version to 1.19.5 from 1.19.4, 1.19.4, 1.19.4
* d338c8b7 Update github-go-version to 1.20.0 from 1.19.5, 1.19.5, 1.19.5
* ad4919a8 Update github-golangci-lint to 1.51.0 from 1.50.1
* e8ecbf95 Update gomod-golang/text to 0.6.0 from 0.5.0
* f1057b9b Update make-golangci-lint to 1.51.0 from 1.50.1
* ca27e426 doc: Add _parent for decode values and clenaup doc a bit
* b04a650b flac_picture,mpeg: Fix trailing ")" typo in map sym and description
* 57144b2f github-action: Use quotes because yaml (1.20 -> 1.2)
* 0aa6e3e2 gojq: Update rebased fq fork
* 7855b359 gomod: Update non-bump tracked mods and add bump config
* 6e17de36 goreleaser: Use name_template instead of deprecated archive replacements
* 8b49b42f interp: Wrap Binary in decodeValue to fix prompt issue with bits/bytes format
* 2d82c05f mp3: Add max_unknown option to fail decode if too much unknown bits
* f386a515 mp4: Decode qt minor version as YYYY.MM description
* 3555dc67 mp4: Decode tkhd flags
* c3e3b3e9 mp4: Decode udta metadata boxes without meta box
* c49012db mp4: sgpd,sbgp: Change grouping_type to a string
* 63403658 mp4: udta: Handle box with value rest of box
* 55ef7a4b readline: Update fq fork to fix draw issue when using del key
* 1eb5e502 tcp: Ignore TCP option check for now as it seems unreliable in dumps
* 62e2cef5 tcp_segment: Decode standard options and rename maxseg to mss

# 0.2.0

This ended up being a release to cleanup old sins in the decoder internals and change some defaults how binary values work with JSON and string functions.

It also adds a new Time Zone Information Format decoder `tzif` (Thanks Takashi Oguma @bitbears-dev) and a new Apple BookmarkData decoder `apple_bookmark` decoder (Thanks David McDonald @dgmcdona). Also a new function `from_ns_keyed_archiver` was added to convert NSKeyedArchiver encoded objects into JSON.

A possible breaking change is that now all `from`/`to` prefix functions now has a `from_`/`to_` prefix, ex: `from_mp3` instead of `frommp3`. There are some few exceptions to this. Note that the functions named just be the format name, ex `mp3` are still around.

In other fq related news [jq-lsp](https://github.com/wader/jq-lsp) got some fixed and additions and seems to work fine with neovim. It's also possible to use jq-lsp with vscode using [vscode-jq](https://github.com/wader/vscode-jq).

## Changes

- All functions that had a `from`/`to` prefix now has the prefix `from_`/`to_`. This is to be easier to read and more consistent, there are still some exceptions like `tovalue`, `torepr`, `tobytes` etc but in general anything that does not deal with primitive types is now `snake_case`. #535
- Change default `bit_formats` option value (how raw bits values are represented in JSON) from `snippet` to `string`. `snippet` meant truncated bits as base64. Now all bits are included as a UTF-8 string. The string will be binary safe (not lose any data) when used internally in fq but will lose data when represented in JSON as some bytes can't be encoded as UTF-8. #499
- Don't auto convert to binary for string/regexp functions, turned out this is very confusing. Now you have to manually use `tobytes` etc to convert to binary value. #540
  ```sh
  # This used to not work as test/1 would convert decode values to the source bytes
  # (0x00 0x00 0x00 0x01) in this case. Now the jq value (symbolic in this case) will
  # be used instead. You can do ".test | tobytes" to get old behavior.
  #
  # find all types with a "mdta." prefix
  $ fq -o line_bytes=10 'grep_by(.type | test(`^mdta\.`))' file.mp4
       │00 01 02 03 04 05 06 07 08 09│0123456789│.boxes[3].boxes[2].boxes[0].boxes[2].boxes[0]{}: box
  0x528│      00 00 00 1c            │  ....    │  size: 28
  0x528│                  00 00 00 01│      ....│  type: "mdta.title" ("\x00\x00\x00\x01")
  0x532│00 00 00 14 64 61 74 61 00 00│....data..│  boxes[0:1]:
  0x53c│00 01 00 00 00 00 74 65 73 74│......test│
       │00 01 02 03 04 05 06 07 08 09│0123456789│.boxes[3].boxes[2].boxes[0].boxes[2].boxes[1]{}: box
  0x546│00 00 00 25                  │...%      │  size: 37
  0x546│            00 00 00 02      │    ....  │  type: "mdta.encoder" ("\x00\x00\x00\x02")
  0x546│                        00 00│        ..│  boxes[0:1]:
  0x550│00 1d 64 61 74 61 00 00 00 01│..data....│
  0x55a│00 00 00 00 4c 61 76 66 35 39│....Lavf59│
  0x564│2e 32 37 2e 31 30 30│        │.27.100│  │
  ```
- Fix panic when cancel (ctrl-c etc) before interpreter is executing. Thanks @pldin601 for reporting. #495
- Fix error using JQValue:s in assign/update paths, ex `.[<JQValue here>] = 123` #509
- Rename fields added for bit-ranges not used by a decoder from `unknown#` to `gap#`. "unknown" is probably a useful field name in some formats and "gap" describe better what it is. #500
- Big decode API internals refactor to split scalars types into their own go types so they can store per type specific values. This also opens up for more ways to make fq both faster and more memory efficient. It also makes the decode API more type safe and makes it possible to experiment with decode DLS that uses chained methods etc. #523

## Decoder changes

- `apple_bookmark` New Apple BookmarkData decoder. Thanks David McDonald @dgmcdona. #493
- `bplist`
  - Fix decoding of UID types
  - Adds a `lost_and_found` array with unused values
  - Fix an endian issue for unicode strings
  - Add NSKeyedArchiver to JSON helper function `from_ns_keyed_archiver`, see `bplist` docs for details on how to use it. Thanks David McDonald @dgmcdona. #493
  ```
  # decode bplist, from_ns_keyed_archiver converts to JSON plist and then into object data as JSON, find app bookmarks keys and expand them as bookmark data and convert to represented JSON, and finally build path to applications
  $ fq -r 'from_ns_keyed_archiver | (.. | .Bookmark? // empty) |= (apple_bookmark | torepr) | .. | .target_path? // empty | join("/")' recentapps.sfl2
  System/Applications/Utilities/Terminal.app
  Applications/Spotify.app
  System/Applications/Calculator.app
  System/Applications/Preview.app
  Applications/Alacritty.app
  Applications/DB Browser for SQLite.app
  System/Applications/System Preferences.app
  System/Library/CoreServices/Applications/Directory Utility.app
  System/Applications/Utilities/Activity Monitor.app
  Applications/Safari.app
  ```
- `tzif` new Time Zone Information Format decoder. Thanks Takashi Oguma @bitbears-dev. #498
- `mp4`
  - Map `mdta` metadata namespace and key names for `ilst` child boxes. #521
  ```sh
  $ fq 'grep_by(.type=="ilst").boxes | map({key: .type, value: .boxes[0].data}) | from_entries' file.mp4
  # create object with all ilst key/value pairs
  {
    "mdta.encoder": "Lavf59.27.100",
    "mdta.title": "test"
  }
  # query specific value
  $ fq -r 'grep_by(.type=="mdta.encoder").boxes[0].data | tovalue' file.mp4
  Lavf59.27.100
  ```
  - Support `sidx` version 1. #506
  - Add description and symbolic values for traf sample flags, makes it easier to see and query for I-frames etc. #514
  ```
  # which boxes has depends_on flags
  $ fq 'grep_by(.sample_depends_on) | parent.type' fragmented.mp4
  ```
  - Support PNG codec mapping. #492
  - Decode `pdin` boxes. #524
  - Decode `hnti` boxes. #513
- `mp3_tags` Add VBRI support and split into into `mp3_frame_xing` and `mp3_frame_vbri` decoders. #525

## Changelog

* 7fa8b635 Add related file format projects to README
* 4fdb7362 Update docker-golang to 1.19.4 from 1.19.3
* 519eff6c Update github-go-version to 1.19.4 from 1.19.3, 1.19.3, 1.19.3
* 2a91d293 Update gomod-golang/text to 0.5.0 from 0.4.0
* cb15b371 added checks to prevent infinite looping and recursion
* c2445335 added some sfl2 test files to bplist package
* 7d13cf73 adds flag parsing to applebookmark
* 71b17d03 apple bookmarkdata decoder initial commit
* 8f39ef63 bplist: Harmonize ns_keyed_archive jq style a bit
* cba72dbd bplist: added overload for from_ns_keyed_archiver jq func
* 129b4b70 bplist: doc: update docs to reflect changes to ns_keyed_archiver
* 9dab3c60 bplist: minor fix to from_ns_keyed_archiver
* 448c3efb bplist: update docs with from_ns_keyed_archiver reference, add error case to function
* a9047c02 bplist: updates from_ns_keyed_archiver to do automatic torepr based on format detection
* 4a28e44f changes decoder package name from bookmark to apple_bookmark
* d0b044c2 converts to snake_case and refactors decode helper
* d199793a created stack type
* e77f7769 decode,interp: Rename unknown gap fields from "unknown#" to "gap#"
* a85da295 decode: Make FieldFormat usage more consistent
* 9b81d4d3 decode: More type safe API and split scalar into multiple types
* 3ec0ba3f decode: add ns_keyed_archiver, restructure apple decoder into apple package
* 330d5f7f decode: apple_bookmark: simplifies flag decoding
* 93f2aa5d decode: change PosLoopDetector to use generics
* 7e98b538 decode: fix type on defer function call, test: add loop.fqtest
* a873819e decode: fixes endian of unicode strings
* f747873d decode: implements lost and found for unreferenced objects
* b45f9fa6 decode: improve stack push/pop
* a162e07b decode: minor change to method receiver name
* 3232f9cc decode: moves PosLoopDetector into its own package
* 7c9504c7 decode: moves macho decoder to apple package
* 70834678 decode: remove dead code from ns_keyed_archiver
* 7ab44662 decode: remove unused field from decoder, unused parens from torepr
* bdb81662 decode: removed unnecessary struct
* 98eab8cb decode: rename parameter for consistency
* 04379df8 decode: revert decode.D back, place posLoopDetector in apple_bookmark
* 7fb674b5 decode: unexport methods
* fa368bb7 decode: updates all.go with correct macho path
* 0287ffa4 decoding well but torepr needs work
* 42debe58 dev,doc,make: Cleanup makefile and have proper targets for *.md and *.svg
* 423bab9e dev,test: Use jqtest code from jqjq for jq tests
* 6fc84a88 doc,dev: Add more usage and dev tips
* 2fc16ae2 doc: Add some padding margin to formats table to make it less likely to cause git conflicts
* 62f377c2 doc: fixes snippet for recursive bookmark searching
* 22064f50 doc: remake
* 4aad2fde doc: remake
* b872b1a3 doc: remake
* 1e1fc551 fixed one more snake_case letter
* d0b76cae fixes broken test and removes long link from markdown body
* 5146f28d fixes broken test for all.fqtest
* 253033cc fixes broken uid parsing in plist decoder
* f535ad3d fixes spacing in jq files
* 64351e8b fixes tests and adds torepr test
* c7d00b87 fixes unknown bit ranges
* 8f930aac forgot to add bookmark.jq in last commit
* 164e527b gojq: Update rebased fq fork
* 6c869451 gojq: Update rebased fq fork
* 578b84d4 interp,display: Add workaround for go 1.18 when escaping 0x7f
* 42d9f2c2 interp,help: Properly count line length when breaking on whole words
* 8d69f1fb interp: Change default bits_format=string
* 6c229d73 interp: Don't auto convert to binary for string functions, is just confusing
* 568afff3 interp: Fix panic when trigger before any context has been pushed
* e3ae1440 interp: Rename to/from<format> functions to to_/from_<format>
* ba88a684 interp: mimic jq: if expr arg is given read stdin even if tty
* 9bd65f93 migrates tests to per-sample files
* f7d7a49f missed a letter on last commit - converting to snake_case
* 2f37cb55 mod: Update modules not tracked with bump
* 55f4f1aa moved a flag bit fields into correct positions
* 9e5a072e mp3_frame_tags: Convert to decode group and split to mp3_frame_{xing,vbri} decoders
* 48522e3c mp3_tags,mp3: Add VBRI header support and rename tags to tag as there is only one
* 83ccedc5 mp4,decode: Properly decode ilst items (both mdta and mdir)
* 1dea40e6 mp4,doc: Add JSON box tree example and reorder a bit
* b1b3b63d mp4: Add namespace to mdta ilst boxes
* 7b60b24a mp4: Add pdin box support
* ef2d5232 mp4: Add png mapping
* 5fb81a14 mp4: Add sym and description for traf sample flags
* 1d6ce2c0 mp4: Decode hint and hnti child boxes
* 9ac453a1 mp4: Fix typo in sample flags sample_is_depended_on description
* a23fe618 mp4: sidx version 1 segment_duration is s64
* 3942db79 pkg/decode/D: Adds PushAndPop, Push, Pop methods. doc: adds help_applebookmark.fqtestdecode: converts applebookmark to use new d.PushAndPop method
* 0c216dff refactors some decoder logic in apple_bookmark for better querying
* 34db9d7f regenerated docs, added tests, fixed torepr
* 0a72635a remade documentation
* 1352598a removed commented out line
* 81269430 removed unnecessary conversions
* 5b1455e7 removed unused function
* 63a3ca20 removes underscore from apple_bookmark package name
* a351c346 removes unused function
* 2ee6360b support tzif (time zone information format)
* 8d5dcff8 test: applebookmark: adds problematic test case
* 63a4e80c test: fixed doc test
* 47a568e0 text,test: Unbeak base64 tests
* 44c91d82 tweaks apple_bookmark markdown documentation
* fd22426b tzif: add help_tzif.fqtest
* c4e7fc79 tzif: moved document to tzif.md
* abde823a tzif: use PeekFindByte() to find end of the string
* 4481a77a tzif: use scalar.Fn() to define a mapper ad hoc
* dbc6fccd updated doc with apple reference
* f5e25fca updated docs
* 6f4d1cb1 updated documentation
* b2aeac6a updates bplist fq tests
* a23ac8f5 updates fqtest for torepr in apple_bookmarkdata


# 0.1.0

Adds `avi` decoder and replace `raw` with more convenient `bits` and `bytes` format. Otherwise mostly small updates and bug fixes.

Increase minor version. fq does not use semantic versioning (yet) but it's probably a good idea to start increase minor version when adding features to be able to do patch releases.

In other fq related news:
- I gave a [talk about fq](https://www.youtube.com/watch?v=-Pwt5KL-xRs&t=1450s) at [No Time To Wait 6](https://mediaarea.net/NoTimeToWait6) a conference about open media, standardization, and audiovisual preservation.
- While prototyping writing decoders directly in jq for fq I ended up [implementing jq in jq](https://github.com/wader/jqjq). Still thinking and working on how to do decoders in jq.

## Changes

- Replace `raw` format with `bits` and `bytes` format that decodes directly to a binary with different unit size.
  ```sh
  $ echo -n 'hello' | fq -d bytes '.[-3:]' > last_3_bytes
  $ echo 'hello' | fq -d bytes '.[1]'
  101
  $ echo 'hello' | fq -c -d bits '[.[range(8)]]'
  [0,1,1,0,1,0,0,0]
  ```

## Decoder changes

- `avc_au` Support annexb format (used in AVI). #476
- `avi` Add AVI (Audio Video Interleaved) decoder. #476
  ```sh
  # extract samples for stream 1
  $ fq '.streams[1].samples[] | tobytes' file.avi > stream01.mp3
  ```
- `bits` Replaces `raw` but is a binary using bit units. #485
- `bytes` Replaces `raw` but is a binary using byte units. #485
- `bplist`
  - Fix signed integer decoding. #451 @dgmcdona
  - Use correct size for references and check for infinite loops. #454 @dgmcdona
- `flac_frame` Correctly decode zero escape sample size. #461
- `id3v2` Fix decoding of COMM and TXXX with missing null terminator. #468
- `matroska` Updated to latest specification. #455
- `mp3_frame` Use frame size calculation from spec instead of own as it seems to not work in some cases. #480
- `mp3_frame_tags` Replaces `xing` and also decodes "lame extensions" for both Xing and Info. #481
- `raw` Removed. #485
- `wav` More codec symbol names and now shares RIFF code with AVI decoder. #476
- `yaml` Fix type panic for large integers. #462

## Changelog

* 7b6492ee Improve README.md a bit, one more demo and move up usage
* 4e069625 Update docker-golang to 1.19.2 from 1.19.1
* e0334497 Update docker-golang to 1.19.3 from 1.19.2
* f3f2648b Update github-go-version to 1.19.2 from 1.19.1, 1.19.1, 1.19.1
* 003197eb Update github-go-version to 1.19.3 from 1.19.2, 1.19.2, 1.19.2
* 453963dd Update github-golangci-lint to 1.50.1 from 1.50.0
* 56dcb3a0 Update gomod-BurntSushi/toml to 1.2.1 from 1.2.0
* 101b2806 Update gomod-golang/text to 0.3.8 from 0.3.7
* d80f12c7 Update gomod-golang/text to 0.4.0 from 0.3.8
* 753927ba Update make-golangci-lint to 1.50.1 from 1.50.0
* 4d8dd5c5 adds check for recursion in decodeReference, adds test to verify fix
* b7c4576c adds necessary cast
* 46b7ab32 adds test to verify fix
* 4ee7dd8a changes Errorf to Fatalf on infinite loops
* 41b2d1ad cli: Better decode error help
* 7254b0f9 decode,elf,fuzz: TryBytesRange error on negative size
* bafd1f56 decode,fuzz: Signed integer (S) read require at least one bit
* 2a86d323 doc,rtmp,pcap,markdown: Add more examples
* 0b9b0173 doc: Add gomarkdown to license/dependencies
* 4bfd9d81 doc: Add link to nttw6 presentation video and slides
* fb1a91ac drop indented else block per lint
* 4dd594c1 fixes bad path in test output
* f9a1d3f4 fixes calculation of floating point lengths
* 236fbc17 fixes reference calculation to use reference size from trailer
* ac86f931 fixes signed integer parsing
* fb2a2b94 flac,fuzz: Fatal error on negative partition sample count
* 7859be1e flac_frame: Properly decode zero escape sample size
* 7cb2a6c9 fuzz: gotip not needed anymore
* cef4245b fuzz: make fuzz GROUP=mp4 to fuzz one group
* 413d4250 gofmt
* 349d9497 gojq: Update rebased fq fork
* 450f5844 gojq: Update rebased fq fork
* d8641ab1 gomod: Update modules that lack bump config
* f66e2244 id3v2: In the wild COMM and TXXX frame might not have a null terminator
* b09d6116 makes dictionary key type checking more sensible
* d07b2eec markdown,fuzz: Update gomarkdown
* 646f32d5 matroska: Fix path tests and make _tree_path more robust
* e748079e matroska: Update spec and regenerate
* 1c7d3252 mod: Update ones without bump config
* 2de87539 mp3_frame: Fix issue calc frame size for some configs
* c3a0686c mp3_frame_tags: Refactor and rename xing format to mp3_frame_tags
* d75748d8 mp4: Decode more sample flags
* c93301fc raw,bits,bytes: Replace raw format with bits and bytes format that decode to a binary
* b08e25ce removes unnecessary cast
* 2b3adbe8 renames test data file
* 0cf46e11 wav,avi,avc_au: Add avi decoder and refactor wav decoder
* 26069167 yaml,fuzz: gojq.Normalize value to fix type panic

# 0.0.10

## Changes

- Add `bplist` Apple Binary Property List decoder. Thanks David McDonald @dgmcdona #427
- Add `markdown` decoder. #422
- Fix panic when interrupting (ctrl-c) JSON output (`fq tovalue file ` etc), #440
- Fix issue using `debug` (and some other native go iterator functions) inside `path(...)`, which is used by assign (`... = ...`) expressions etc. #439
- Fix issue and also make `toactual` and `tosym` work more similar to `tovalue`. #432
- Fix issue with unknown fields (gaps found after decoding) where one continuous gap could end up split into two of more unknown fields. #431
- More format documentation and also nicer help output. Also now all documentation is in markdown format. #430 #422
  ```
  # or help(matroska) in the REPL
  $ fq -h matroska
  matroska: Matroska file decoder

  Decode examples
  ===============

    # Decode file as matroska
    $ fq -d matroska . file
    # Decode value as matroska
    ... | matroska

  Lookup element using path
  =========================

    $ fq 'matroska_path(".Segment.Tracks[0)")' file.mkv

  Get path to element
  ===================

    $ fq 'grep_by(.id == "Tracks") | matroska_path' file.mkv

  References
  ==========
  - https://tools.ietf.org/html/draft-ietf-cellar-ebml-00
  - https://matroska.org/technical/specs/index.html
  - https://www.matroska.org/technical/basics.html
  - https://www.matroska.org/technical/codec_specs.html
  - https://wiki.xiph.org/MatroskaOpus
  ```

## Decoder changes

- `ar` Allow empty integer strings. For example owner id can be an empty string. #428
- `bitcoin_blkdat` Assert that there is a header. As the format is part of the probe group this speeds up probing. #402
- `bplist` Add Apple Binary Property List decoder.
  ```sh
  $ fq '.objects.entries[0] | .key, .value' Info.plist
      │00 01 02 03 04 05 06 07 08 09│0123456789│.objects.entries[0].key{}:
  0x32│               5c            │     \    │  type: "ascii_string" (5) (ASCII encoded string)
  0x32│               5c            │     \    │  size_bits: 12
      │                             │          │  size: 12
  0x32│                  43 46 42 75│      CFBu│  value: "CFBundleName"
  0x3c│6e 64 6c 65 4e 61 6d 65      │ndleName  │
      │00 01 02 03 04 05 06 07 08 09│0123456789│.objects.entries[0].value{}:
  0x1ea│         5f                  │   _      │  type: "ascii_string" (5) (ASCII encoded string)
  0x1ea│         5f                  │   _      │  size_bits: 15
  0x1ea│            10               │    .     │  large_size_marker: 1 (valid)
  0x1ea│            10               │    .     │  exponent: 0
  0x1ea│               18            │     .    │  size_bigint: 24
      │                             │          │  size: 24
  0x1ea│                  41 70 70 6c│      Appl│  value: "AppleProResCodecEmbedded"
  0x1f4│65 50 72 6f 52 65 73 43 6f 64│eProResCod│
  0x1fe│65 63 45 6d 62 65 64 64 65 64│ecEmbedded│
  ```
  - Supports `torepr`
  ```sh
  $ fq torepr.CFBundleName Info.plist
  "AppleProResCodecEmbedded"
  ```
- `elf`
  - More robust decoding when string tables are missing. #417<br>
    ```sh
    # extract entry opcodes and disassemble with ndisasm
    $ fq -n '"f0VMRgIBAQAAAAAAAAAAAAIAPgABAAAAeABAAAAAAABAAAAAAAAAAAAAAAAAAAAAAAAAAEAAOAABAAAAAAAAAAEAAAAFAAAAAAAAAAAAAAAAAEAAAAAAAAAAQAAAAAAAAAAAAAEAAAAAAAAAAQAAAAAAIAAAAAAAsDxmvwYADwU=" | frombase64 | . as $b | elf | $b[.header.entry-.program_headers[0].vaddr:]' \
    | ndisasm -b 64 -
    00000000  B03C              mov al,0x3c
    00000002  66BF0600          mov di,0x6
    00000006  0F05              syscall
    ```
  - Now decodes program header notes. #421
- `markdown` Add decoder. Is used in fq to render CLI help. #422
  ```sh
  # array with all level 1 and 2 headers
  $ fq -d markdown '[.. | select(.type=="heading" and .level<=2)?.children[0]]' README.md
  [
    "fq",
    "Usage",
    "Presentations",
    "Install",
    "TODO and ideas",
    "Development and adding a new decoder",
    "Thanks and related projects",
    "License"
  ]
  ```
- `matroska` Add support for sample lacing. Used by FLAC samples etc. #404
- `mp4` More codec names and also use official names from mp4ra.org. #424<br>
  ```sh
  # show details of first two track in file
  $ fq -o line_bytes=10 '.tracks[0,1]' big_buck_bunny.mp4
          │00 01 02 03 04 05 06 07 08 09│0123456789│.tracks[0]{}: track
  0x00910a│20 68 10 01 a0 40 0e 20 8c 1b│ h...@. ..│  samples[0:1295]:
  0x009114│c2 2b 99 09 84 42 60 a8 c4 60│.+...B`..`│
  *       │until 0x541697.7 (5473678)   │          │
          │                             │          │  id: 1
          │                             │          │  data_format: "mp4a" (MPEG-4 Audio)
          │00 01 02 03 04 05 06 07 08 09│0123456789│.tracks[1]{}: track
  0x00a5e6│                           00│         .│  samples[0:1440]:
  0x00a5f0│00 00 0c 06 00 07 8b 71 b0 00│.......q..│
  0x00a5fa│00 03 00 40 80 00 00 00 15 06│...@......│
  *       │until 0x540959.7 (5464939)   │          │
          │                             │          │  id: 2
          │                             │          │  data_format: "avc1" (Advanced Video Coding)

  ```
- `html` Handle leading doc type and processing directives. #414

## Changelog

* a77cec92 Added documentation and tests, fixed bad date parsing
* d784db69 Adds support for Apple Binary Plist, version 00
* 5711f290 Code fixes from PR, still need to add tests and testdata
* 6b04f2de Documentation cleanup
* bcccde23 Fixes and embeds documentation
* ebae938d Fixes bug in integer parsing
* 368d183b Size check on nBits to save memory
* 84ca1010 Update docker-golang from 1.19.0 to 1.19.1
* c47c3866 Update github-go-version from 1.19.0, 1.19.0, 1.19.0 to 1.19.1
* 816169b6 Update github-golangci-lint to 1.50.0 from 1.49.0
* 21f2980e Update make-golangci-lint to 1.50.0 from 1.49.0
* 5f619940 adds function for decoding fixed sized arrays
* f08f44f1 ar: Integer strings might be empty
* 004406de bitcoin_blkdat,bitcoin_block: Make sure there is a header if blkdat
* 421b2b30 bplist: Fix unknown field for singletons and add torepr tests
* 16b01211 bplist: Make torepr convert to values
* fe64530e csv: Add tsv and header example
* cb3dc802 decode,tar: Add scalar description and Try* helpers
* a6429ffe decode: Remove RangeSorted flag as we can decide on array/struct instead
* a468684a deps: Manual update ones not using bump
* a7a101ca doc,help: Nicer format help and move help tests into each format
* 725ab1b1 doc,html,xml: Add more documentation and examples
* abd19ed8 doc: Fix format sections a bit
* 0fdc03a4 doc: Fix some incorrect example prompts
* 5382d46a elf: Basic program header notes decoding
* 12105d8c elf: Treat missing string tables as empty to be more robust
* 3deceeeb fixes from PR comments
* 226a9a3e generics: Use more from x/exp
* 404b1704 gojq: Update fq fork
* 376f0ebb gojq: Update rebased fq fork
* 87b2c6c1 help,doc: Use markdown for format documentation again
* 8016352b html: Handle html with leading doctype etc
* 768df301 interp,decode: For struct use map to lookup field
* c4219d69 interp: Fix interrupt panic for cli eval
* 00ee10a1 interp: Make to{actual,sym} behave similar to tovalue
* 00a50662 markdown: Add decoder
* 7749e1b5 matroska: Add proper lacing support
* 20a15372 mp4: Fix data_format typo
* 2655ba09 mp4: More codec names (from mp4ra.org)
* 7cd43b49 performance: increase performance by map usage
* 6a6fec54 range,decode: Use own range sort impl to speed up a bit
* 0f35fe48 ranges,decode: Correctly skip empty ranges when adding unknown fields
* ea81efec readline: Update rebased fq fork
* 369f4016 removed unnecessary type conversions
* 3198602d removed unused return type
* 7d865343 sortex: Package with type safe sort helpers
* 808202fa test: Skip go test with -race by default
* 12836abe updates fqtest
* 1e47f4f2 updates tests post integer-bug fix
* 3d8ea1de updates torepr for data type
* 1385b5d0 wasm: Add some documentation
* d6316d5c wav: Decode smpl chunk

# 0.0.9

## Changes

- New `wasm` WebAssembly Binary Format decoder by Takashi Oguma @bitbears-dev<br>
  ```sh
  # show part of code_section
  $ fq '.sections[4].content.code.x[0].code.e | d' add.wasm
      │00 01 02 03 04 05 06 07 08 09│0123456789│.sections[4].content.code.x[0].code.e[0:4]:
      │                             │          │  [0]{}: in
  0x3c│                           20│          │    opcode: "local.get" (0x20)
  0x46│01                           │.         │    x: 1 (valid)
      │                             │          │  [1]{}: in
  0x46│   20                        │          │    opcode: "local.get" (0x20)
  0x46│      00                     │  .       │    x: 0 (valid)
      │                             │          │  [2]{}: in
  0x46│         6a                  │   j      │    opcode: "i32.add" (0x6a)
      │                             │          │  [3]{}: in
  0x46│            0b               │    .     │    opcode: "end" (0xb)
  ```
  ```sh
  # count opcode usage
  $ fq '.sections[] | select(.id == "code_section") | [.. | .opcode? // empty] | count | map({key: .[0], value: .[1]}) | from_entries' add.wasm
  {
    "end": 1,
    "i32.add": 1,
    "local.get": 2
  }
  ```
  ```sh
  # list exports and imports
  $ fq '.sections | {import: map(select(.id == "import_section").content.im.x[].nm.b), export: map(select(.id == "export_section").content.ex.x[].nm.b)}' add.wasm
  {
    "export": [
      "memory",
      "add"
    ],
    "import": []
  }
  ```
- Decode value display now shows address bar on new format or buffer. Should make it easier to spot changes and read hex and ASCII view. #365<br>
  Examples of PCAP with different formats and TCP stream buffers:
  <pre>
  ...
         │<ins>00 01 02 03 04 05 06 07 08 09</ins>│0123456789</ins>│      packet{}: (ether8023_frame)
  0x00668│   00 0a 95 67 49 3c         │ ...gI<   │        destination: "00:0a:95:67:49:3c" (0xa9567493c)
  0x00668│                     00 c0 f0│       ...│        source: "00:c0:f0:2d:4a:a3" (0xc0f02d4aa3)
  0x00672│2d 4a a3                     │-J.       │
  0x00672│         08 00               │   ..     │        ether_type: "ipv4" (0x800) (Internet Prot...
         │<ins>00 01 02 03 04 05 06 07 08 09</ins>│0123456789</ins>│        payload{}: (ipv4_packet)
  0x00672│               45            │     E    │          version: 4
  ...
         │<ins>00 01 02 03 04 05 06 07 08 09</ins>│0123456789</ins>│          payload{}: (tcp_segment)
  0x00686│               00 50         │     .P   │            source_port: "http" (80) (World Wide ...
  ...
         │                             │          │  ipv4_reassembled[0:0]:
         │                             │          │  tcp_connections[0:1]:
         │                             │          │    [0]{}: tcp_connection
         │                             │          │      client{}:
         │                             │          │        ip: "192.168.69.2"
         │                             │          │        port: 34059
         │                             │          │        has_start: true
         │                             │          │        has_end: true
         │                             │          │        skipped_bytes: 0
         │<ins>00 01 02 03 04 05 06 07 08 09</ins>│<ins>0123456789</ins>│
    0x000│47 45 54 20 2f 74 65 73 74 2f│GET /test/│        stream: raw bits
    0x000│65 74 68 65 72 65 61 6c 2e 68│ethereal.h│
    *    │until 0x1bc.7 (end) (445)    │          │
  ...
  </pre>
- Add `--unicode-output`/`-U` argument to force use of Unicode characters to improve output readability. #377
  - For example useful when piping to less and you want fancy unicode and colors:<br>
  `fq -cU d file | less -r`
- `to_entries` now preserves struct field order. #340
- Experimental <code>&#96;raw string&#96;</code> literal support. Work the same as golang raw string literals. Useful in REPL when pasting things etc but should probably be avoided in jq scripts. #371
- Properly fail lexing of invalid binary, octal and hex number literals. #371
- REPL completion now include all functions. Before some functions with multiple argument counts were skipped. #375
- Switch to new gopacket fork with speedup and bug fixes. Remove SLL2 workarounds in fq. #386

## Decoder changes

- `csv` Correctly handle decode values when `tocsv` normalize to strings. Before array and object ended up being JSON serialized to strings. #341
  - Normalize to strings is done so that non-string scalars can be used:
    ```
    $ fq -n '[[1,true,null,"a"]] | tocsv'
    "1,true,,a\n"
    ```
- `dns` DNS over UDP format was accidentally used to probe TCP streams #366
- `elf` Remove redundant `program_header` struct
- `flac`
  - Add 32 bit samples support. #378 Thanks @ktmf01
  - Properly decode/checksum samples in partitions with raw samples. #379 Thanks @ktmf01<br>
    Now successfully decodes all test cases from https://github.com/ietf-wg-cellar/flac-test-files
- `jsonl` Add decoder. Decodes JSON lines. There is also `fromjsonl` and `tojsonl`. #374
- `macho`
  - Split FAT Macho decoding into `macho_fat` format which also fixed handling of file offsets in sections. #362
  - Decode symbol and string sections. #352
- `matroska` Remove new lines in descriptions. Messes up tree. #369
- `mp3_frame`
  - Support LSF (low sampling frequency) frames. #376
  - Skip trying to figure out what main data is for current frame and not. Was incorrect and doing it properly probably requires hoffman decoding. #376
- `pcap` Support files with nanosecond precision. Has a different magic. #382
- `prores_frame` Add basic decoder. Decodes container and fram header. #396 Thanks @Doctor-love for test files
- `tar` Fix regression when decoding number fields. Now ok again to be empty string. #372
- `wasm` Add WebAssembly Binary Format decoder. #383 Thanks to Takashi Oguma @bitbears-dev
  - Decodes to a tree following the [WASM binary grammar specification](https://webassembly.github.io/spec/core/binary/index.html)
- `yaml` Fail on trailing data. Before it succeeded with the last value. #373
- `zip`
  - Don't require PK header as there seems to be zip files with prepended data. #359
  - Correctly limit amount of backwards search for EOCD (end of content directory). #363
- `xml` Correctly handle decode values when `toxml` normalize to strings. Before array and object ended up being JSON serialized to strings. #341
- `xml`
  - Change attribute prefix to `@` instead of `-` and make it an option `attribute_prefix`. #401
  - Skip default namespace in element names. #389
  - Always include attributes and children even when empty in array mode. Makes it a lot easier to work with as you can assume `.[1]` will be attributes and so on. #357
  - Normalize to strings is done so that non-string scalars can be used:
    ```
    $ fq -nr '{a: {"-boolean": true, "-number": 123, "-null": null}} | toxml'
    <a boolean="true" null="" number="123"></a>
    ```
  - Allow and ignore trailing `<?procinstr?>` and improve trailing data error message. #368
  - Correctly sort if any `#seq` is found and also properly sort negative `#seq`. #384

## Changelog

* 0cd846a1 *extra: Rename <pkg>extra to just <pkg>ex and refactor to use generics
* fb583e2c Add 32 bps FLAC to test
* c1d5b2b1 Add sample size entry to list for 32bps flac streams
* 3f209c46 Fix decoding of FLAC raw entropy partition
* 25061aca Update docker-golang from 1.18.4 to 1.18.5
* 0de2c906 Update docker-golang from 1.18.5 to 1.19.0
* 7b8d95bf Update github-go-version from 1.18.4, 1.18.4, 1.18.4 to 1.18.5
* 103991f7 Update github-go-version from 1.18.5, 1.18.5, 1.18.5 to 1.19.0
* 4255b87a Update github-golangci-lint from 1.47.2 to 1.47.3
* 198305ec Update github-golangci-lint from 1.47.3 to 1.48.0
* fa9fec30 Update github-golangci-lint from 1.48.0 to 1.49.0
* f579e9c3 Update make-golangci-lint from 1.47.2 to 1.47.3
* c8069d22 Update make-golangci-lint from 1.47.3 to 1.48.0
* 004eb564 Update make-golangci-lint from 1.48.0 to 1.49.0
* abcc7366 add ULEB and SLEB to known words for spell check
* 9238251b ci: Skip -race for windows and macos
* 913f5780 columnwriter,dump: Add Column interface and refactor into BarColumn and MultiLineColumn
* 5d9ffead decode,scalar: Map empty string also else sym might ends up nil
* 326dada7 decode: Add LEB128 readers
* 502f451c decode: Refactor to use scalar type assert helper
* 840292ba decode: Simplify compound range sort behaviour
* 15f7c67a dev,fuzz: Add some useful retrigger snippets
* 46dca8cd dns: Don't use dns (udp) format for tcp also
* c233215a dns: Rename isTCP to hasLengthHeader
* ed424783 doc,interp: Update and add more examples
* f247edb5 doc: Update README demo a bit with new features
* 3613b6d4 elf: Remove redundant program_header struct
* 8a19978b flac: Make gen script generate correct fqtest files
* 2bfbe9a9 flac_frame: Cleanup some dev lefterovers and todos
* 64b23659 fqtest: Run tests in parallel
* af35b284 gojq: Preserve keys order for to_entries when used with JQValue
* 804ad1e2 gojq: Update fq fork
* add3dcfd gojq: Update fq fork, fix scope argcount issue
* d898732c gojq: Update fq fork, new scope function, rawstring, stricter integers
* 394717ca gopacket: Switch/update to new fork, remove SLL2 hack
* 4eae7ffd interp,doc: Add -R raw string slurp hint to -s help
* d8792fd1 interp,dump: Correctly flush columns if data will be shown
* 29005c70 interp,dump: Show address bar for root, nested roots and on format change
* c7559b59 interp: Add --unicode-output/-U to force use of unicode
* 9e447c9a interp: Use RegisterFS instead of format files
* 701c67c1 jsonl: Add decoder, also tojsonl encoder
* bc6cffde lint,decode,fuzz:: Fix nilerr warnings, one real one should be ignored for now
* 3c21b058 lint: Fix ioutil deprecation, reformat for new doc standard
* b2d4e6d9 macho: Decode cmd symtab symbols
* 725c8e83 macho: Split into macho/macho_fat, fix offset issue and add string decoding
* 2e407386 matroska: Strip newlines in description
* cf15661e mp3_frame: Add LSF support and fix incorrect main data handling
* 74c7dc4e pcap: Add ns support and add header field
* 8fc43533 prores_frame: Add basic container and frame header decoder
* dc32ac08 script: Use strings.Builder to collect output
* 0d44b937 tar: Some number fields can be empty
* 545dac8c test: Update tests, go 1.19 uses \xff instead of \u00ff
* ce438872 wasm: `make doc`
* 074c22c9 wasm: add `-timeout 20m` for go test to workaround ci test fail
* cd037c51 wasm: add comment to clarify lazy initialization
* f73965d2 wasm: add wasm to probe list
* 00869b37 wasm: avoid race condition
* db8021c9 wasm: define and use constants for some insturctions
* bcc0dfd9 wasm: fix comment format
* 289ddf59 wasm: fix lint issues
* 3fca7cc0 wasm: fix lint issues
* cbb5a8ed wasm: further simplification
* 934ed9a8 wasm: initial version
* e5cf1731 wasm: make the godoc formatter happy
* b0f3fec8 wasm: remove nolint:unparam which is no longer needed
* e1691dec wasm: remove unused function
* ae4529c4 wasm: run `golangci-lint run --fix`
* fead68de wasm: tidy up
* 3298d181 wasm: to be able to probe
* 2eb17505 wasm: update tests
* d5d9e738 wasm: use FieldULEB128() / FieldSLEB128() defined in the upstream
* 7401d141 wasm: use WRITE_ACTUAL=1 to generate .fqtest files
* 2037b86a wasm: use map, not switch
* ae08bf70 wasm: use s.ActualU() instead of s.Actual.(uint64)
* 63f4a726 wasm: use scalar.UToSymStr for simplicity
* 0ad5a8ec wasm: use underscores for symbol values
* fa20c74c xml,csv,interp: Handle JQValue when string normalizing
* f4e01372 xml,html: Always include attrs and children in array mode
* 9a5fcc89 xml: Allow trailing <?procinstr?>
* 71900c2a xml: Correctly sort if one #seq is found and allow negative seq numbers
* 716323ce xml: Even more namespace fixes
* f24d685a xml: Keep track of default namespace and skip it element names
* 095e1161 xml: Switch from "-" to "@" as attribute prefix and make it an option
* 3623eac3 yaml: Error on trailing yaml/json
* d607bee1 zip: Correctly limit max EOCD find
* 19b70899 zip: Skip header assert as there are zip files with other things appended
# 0.0.8

## Changes

- Add support for some common structured serialization formats: #284 #335
  - XML, `toxml`, `fromxml` options for indent, jq mapping variants (object or array) and order preservation
  - HTML, `fromhtml` options for indent, jq mapping variants (object or array) and order preservation
  - TOML, `totoml`, `fromtoml`
  - YAML, `toyaml`, `fromyaml`
  - jq-flavored JSON (optional key quotes and trailing comma) `tojq`, `fromjq` options for indent #284
    ```sh
    # query a YAML file
    $ fq '...' file.yml

    # convert YAML to JSON
    # note -r for raw string output, without a JSON string with JSON would outputted
    $ fq -r 'tojson({indent:2})' file.yml

    $ fq -nr '{hello: {world: "test"}} | toyaml, totoml, toxml, tojq({indent: 2})'
    hello:
        world: test

    [hello]
      world = "test"

    <hello>
      <world>test</world>
    </hello>
    {
      hello: {
        world: "test"
      }
    }
    $ echo '<doc><element a="b"></doc>' | fq -r '.doc.element."-a"'
    b
    $ echo '<doc><element a="b"></doc>' | fq -r '.doc.element."-a" = "<test>" | toxml({indent: 2})'
    <doc>
      <element a="&lt;test&gt;"></element>
    </doc>
    ```
  - CSV, `tocsv`, `fromcsv` options for separator and comment character
    ```sh
    $ echo -e  '1,2\n3,4' | fq -rRs 'fromcsv | . + [["a","b"]] | tocsv'
    1,2
    3,4
    a,b
    ```
- Add support for binary encodings
  - Base64. `tobase64`, `frombase64` options for encoding variants.
    ```sh
    $ echo -n hello | base64 | fq -rRs 'frombase64 | tostring'
    hello
    ```
  - Hex string. `tohex`, `fromhex`
- Add support for text formats
  - XML entities `toxmlentities`, `fromxmlentities`
  - URL `tourl`, `fromurl`
    ```sh
    $ echo -n 'https://host/path/?key=value#fragment' | fq -Rs 'fromurl | ., (.host = "changed" | tourl)'
    {
      "fragment": "fragment",
      "host": "host",
      "path": "/path/",
      "query": {
        "key": "value"
      },
      "rawquery": "key=value",
      "scheme": "https"
    }
    "https://changed/path/?key=value#fragment"
    ```
  - URL path encoding `tourlpath`, `fromurlpath`
  - URL encoding `tourlencode`, `fromurlencode`
  - URL query `tourlquery`, `fromurlquery`
- Add support for common hash functions:
  - MD4 `tomd4`
  - MD5 `tomd5`
    ```sh
    $ echo -n hello | fq -rRs 'tomd5 | tohex'
    5d41402abc4b2a76b9719d911017c592
    ```
  - SHA1 `tosha1`
  - SHA256 `tosha256`
  - SHA512 `tosha512`
  - SHA3 224 `tosha3_224`
  - SHA3 256 `tosha3_256`
  - SHA3 384 `tosha3_384`
  - SHA3 512 `tosha3_512`
- Add support for common text encodings:
  - ISO8859-1 `toiso8859_1`, `fromiso8859_1`
  - UTF8 `tutf8`, `fromutf8`
  - UTF16 `toutf16`, `fromutf16`
  - UTF16LE `toutf16le`, `fromutf16le`
  - UTF16BE `toutf16be`, `fromutf16be`
    ```sh
    $ echo -n 00680065006c006c006f | fq -rRs 'fromhex | fromutf16be'
    hello
    ```
- Add `group` function, same as `group_by(.)` #299
- Update/rebase readline dependency (based on @tpodowd  https://github.com/chzyer/readline/pull/207) #305 #308
  - Less blinking/redraw in REPL
  - Lots of small bug fixes
- Update/rebase gojq dependency #247
  - Fixes JQValue destructing issue (ex: `<some object JQValue> as {$key}`)
- Major rewrite/refactor how native function are implemented. Less verbose and less error-prone as now shared code takes care of type casting and some argument errors. #316
- Add `tojson($opts)` that support indent option. `tojson` still works as before (no indent).
  ```sh
  $ echo '{a: 1}' | fq -r 'tojson({indent: 2})'
  {
    "a": 1
  }
  ```
- Rename `--decode-file` (will still work) to `--argdecode` be be more consistent with existing `--arg*` arguments. #309
- On some decode error cases fq can now keep more of partial tree making it easier to know where it stopped #245
- Build with go 1.18 #272

## Decoder changes

- `bitcoin` Add Bitcoin blkdat, block, transcation and script decoders #239
- `elf` Use correct offset to dynamic linking string table #304
- `tcp` Restructure into separate client/server objects and add `skipped_bytes` (number of bytes with known missing ACK), `has_start` (has first byte in stream) and `has_end` (has last byte in stream) per direction #251
  - Old:
  ```
        │00 01 02 03 04 05 06 07│01234567│.tcp_connections[0]{}: tcp_connection
        │                       │        │  source_ip: "192.168.69.2"
        │                       │        │  source_port: 34059
        │                       │        │  destination_ip: "192.168.69.1"
        │                       │        │  destination_port: "http" (80) (World Wide Web HTTP)
        │                       │        │  has_start: true
        │                       │        │  has_end: true
   0x000│47 45 54 20 2f 74 65 73│GET /tes│  client_stream: raw bits
   0x008│74 2f 65 74 68 65 72 65│t/ethere│
   *    │until 0x1bc.7 (end) (44│        │
   0x000│48 54 54 50 2f 31 2e 31│HTTP/1.1│  server_stream: raw bits
   0x008│20 32 30 30 20 4f 4b 0d│ 200 OK.│
   *    │until 0x191.7 (end) (40│        │
  ```
  - New:
  ```
        │00 01 02 03 04 05 06 07│01234567│.tcp_connections[0]{}: tcp_connection
        │                       │        │  client{}:
        │                       │        │    ip: "192.168.69.2"
        │                       │        │    port: 34059
        │                       │        │    has_start: true
        │                       │        │    has_end: true
        │                       │        │    skipped_bytes: 0
   0x000│47 45 54 20 2f 74 65 73│GET /tes│    stream: raw bits
   0x008│74 2f 65 74 68 65 72 65│t/ethere│
   *    │until 0x1bc.7 (end) (44│        │
        │                       │        │  server{}:
        │                       │        │    ip: "192.168.69.1"
        │                       │        │    port: "http" (80) (World Wide Web HTTP)
        │                       │        │    has_start: true
        │                       │        │    has_end: true
        │                       │        │    skipped_bytes: 0
   0x000│48 54 54 50 2f 31 2e 31│HTTP/1.1│    stream: raw bits
   0x008│20 32 30 30 20 4f 4b 0d│ 200 OK.│
   *    │until 0x191.7 (end) (40│        │
  ```
- `zip` Add 64-bit support and add `uncompress` option #278
- `matroska` Update and regenerate based on latest spec and also handle unknown ids better #291
- `mp4` Changes:
  - Fix PSSH decode issue #283
  - Add track for track_id references without tfhd box
  - Makes it possible to see samples in fragments without having an init segment.
    Note it is possible to decode samples in a fragment file by concatenating the init and fragment file ex: `cat init frag | fq ...`.
  - Add `senc` box support #290
  - Don't decode encrypted samples #311
  - Add `track_id` to tracks #254
  - Add fairplay PSSH system ID #310
  - Properly handle `trun` data offset #294
  - Skip decoding of individual PCM samples for now #268
  - Add `mvhd`, `tkhd`, `mdhd` and `mehd` version 1 support #258
  - Make sure to preserve sample table order #330
- `fairplay_spc` Add basic FairPlay Server Playback Context decoder #310
- `avc_pps` Correctly check for more rbsp data

## Changelog

* 210940a4 Update docker-golang from 1.18.1 to 1.18.2
* fbeabdc3 Update docker-golang from 1.18.2 to 1.18.3
* 51a414db Update docker-golang from 1.18.3 to 1.18.4
* 3017e8b4 Update github-go-version from 1.18.1, 1.18.1, 1.18.1 to 1.18.2
* c597f7f7 Update github-go-version from 1.18.2, 1.18.2, 1.18.2 to 1.18.3
* dd283923 Update github-go-version from 1.18.3, 1.18.3, 1.18.3 to 1.18.4
* d10a3616 Update github-golangci-lint from 1.45.2 to 1.46.0
* 75b5946c Update github-golangci-lint from 1.46.0 to 1.46.1
* 3ffa9efb Update github-golangci-lint from 1.46.1 to 1.46.2
* 4be8cb91 Update github-golangci-lint from 1.46.2 to 1.47.0
* 1b8f4be8 Update github-golangci-lint from 1.47.0 to 1.47.1
* fc596a7a Update github-golangci-lint from 1.47.1 to 1.47.2
* 62be9223 Update gomod-BurntSushi/toml from 1.1.0 to 1.2.0
* 5db7397a Update make-golangci-lint from 1.45.2 to 1.46.0
* 456742ea Update make-golangci-lint from 1.46.0 to 1.46.1
* 06757119 Update make-golangci-lint from 1.46.1 to 1.46.2
* 3d69e9d0 Update make-golangci-lint from 1.46.2 to 1.47.0
* 2170925d Update make-golangci-lint from 1.47.0 to 1.47.1
* c4199c0f Update make-golangci-lint from 1.47.1 to 1.47.2
* 02f00be9 Update usage.md
* 75169a65 asn1: Add regression test for range decode fix ##330
* b0096bc1 avc_pps: Correct check if there is more rbsp data
* 5d67df47 avro_ocf: Fix panic on missing meta schema
* 417255b7 bitcoin: Add blkdat, block, transcation and script decoder
* a6a97136 decode: Cleanup Try<f>/<f> pairs
* 3ce660a2 decode: Keep decode tree on RangeFn error
* c4dd518e decode: Make compound range sort optional
* 8bb4a6d2 decode: Range decode with new decoder to preserve bit reader
* 342612eb dev: Cleanup linters and fix some unused args
* 78aa96b0 dev: Cleanup some code to fix a bunch of new linter warnings
* 3570f1f0 doc: Add more related tools
* 7aff654a doc: Clarify decode, slurp and spew args
* 0863374f doc: Correct bencode spec URL
* 10cc5518 doc: Improve and cleanup text formats
* b1006119 doc: Typos and add note about Try* functions
* c27646a6 doc: Update and shorten README.md a bit
* b0388722 doc: Use singular jq value to refer to jq value
* a980656c doc: go 1.18 and improve intro text a bit
* a64c28d9 dump: Skip JQValueEx if there are not options
* 40481f66 elf,fuzz: Error on too large string table
* f66a359c elf: Use correct offset to dynamic linking string table
* 64f3e5c7 fairplay: Add basic SPC decoder and PSSH system id
* cae288e6 format,intepr: Refactor json, yaml, etc into formats also move out related functions
* e9d9f8ae fq: Use go 1.18
* 377af133 fqtest: Cleanup path usage
* 2464ebc2 fuzz: Replace built tag with FUZZTEST env and use new interp api
* 0f78687b gojq: Fix JQValue index and destructuring issue and rebase fq fork
* 59c7d0df gojq: Rebase fq fork
* c57dc17d gojq: Rebase fq fork
* 9a7ce148 gojq: Update rebased fq fork
* c1a0cda5 gojq: Update rebased fq fork
* 32361dee gojqextra: Cleanup gojq type cast code
* 9b2e474e gojqextra: Simplify function type helpers
* fd302093 hevc_vps,fuzz: Error on too many vps layers
* efa5e23a icc_profile: Correctly clamp align padding on EOF
* 1ddea1ad interp,format: Refactor registry usage and use function helpers
* a3c33fc1 interp: Add group/0
* 95e61965 interp: Add internal _is_<type> helpers
* 3b717c3b interp: Add to/from<encoding> for some common serialzations, encodings and hashes
* 6b088000 interp: Cast jq value to go value properly for encoding functions
* f5be5180 interp: Cleanup and clarify some format naming
* c7701851 interp: Extract to/from map/struct to own package
* 8dde3ef5 interp: Fix crash when including relatve path when no search paths are set
* 735c443b interp: Improve type normalization and use it for toyaml and totoml
* 81a014ce interp: Make empty _finally fin error on error
* 2dc509ab interp: Refactor dump and revert #259 for now
* ab8c728a interp: Rename --decode-file to --argdecode to be more consistent
* dff3cc11 interp: dump: Fix column truncate issue with unicode bars
* 5109df4a interp: dump: Show address bar for nested roots
* 80214921 interp: help: Fix incorrect options example
* 76714349 mapstruct: Handle nested values when converting to camel case
* c92f4f13 matroska: Update ebml_matroska.xml and allow unknown ids
* c2a359bd mod: Update golang.org/x/{crypto,net}
* 3780375d mp3: Use d.FieldValueU and some cleanup
* 7b27e506 mp4,bitio: Fix broken pssh decoding and add proper reader cloning to bitio
* 6b00297e mp4,senc: Refactor current track/moof tracking and add senc box support
* 8228ecae mp4: Add track id field and add track for tfhd with unseen track_id
* ea2cc3c2 mp4: Don't  decode encrypted samples
* c6d0d89c mp4: Don't range sort samples, keep sample table order
* 7d25fbfd mp4: Properly use trun data offset
* ba844eb0 mp4: Skip fields for pcm samples for now
* 0e02bb66 mp4: iinf: Only assume sub boxes for version 0
* 2e328180 mp4: mvhd,tkhd,mdhd,mehd: Add version 1 support
* 44bab274 readline: Rebase on top of tpodowd's redraw/completion fixes PR
* a5122690 readline: Rebase on top of tpodowd's update PR
* 54dcdce9 readline: Update fq fork
* 6e7267d2 readme: add MacPorts install details
* 76161a1b scalar,mp4,gzip,tar: Add timestamp to description
* 9133f0e5 scalar: Add *Fn type to map value and clearer naming
* 34cf5442 tcp: Split into client/server structs and add skipped_bytes and has_start/end per direction
* 1aaaefb0 wav,bencode,mpeg_ps_packet,id3v1: Random fixes
* 47350e46 zip: Add uncompress=false test and some docs
* e6412744 zip: Add zip64 support and uncompress option
* aa694e3f zip: s/Decompress/Uncompress/


# 0.0.7

## Changes

- Format specific options
  - Formats can now have own options
  - Example to skip decoding of samples in a mp4 file use:
  - `fq -d decode_samples=false d file.mp4` or `... | mp4({decode_samples: false}})`
  - To see supported options for a formats see formats documentation, use `fq -h mp4` or `help(mp4)` in the REPL.
- gojq fork rebase:
  - Many performance improvements from upstream
  - Assign to a JQValue will now shallowly turn it into a jq value and then be assigned.
  - Refactor and rewrote large parts to make it easier to rebase and maintain in the future.

## Decoder changes

- `amf0` Add Action Message Format 0 decoder #214
- `hevc_pps` Add H.265/HEVC Picture Parameter Set decoder #210
- `hevc_sps` Add H.265/HEVC Sequence Parameter Set decoder #210
- `hevc_vpc` Add H.265/HEVC Video Parameter Set decoder #210
- `mp3` Add `max_unique_header_config` and `max_sync_seek` options #242
- `mp4` Simplify granule structure a bit #242
- `mp4` Add `decode_samples` and `allow_truncate` options #242
- `flac_frame` Add `bits_per_sample` option #242
- `icmpv6` Add Internet Control Message Protocol v6 decoder #216
- `id3v2` Add v2.0 PIC support
- `ipv6_packet` Add Internet protocol v6 packet decoder #216
- `macho` Remove redundant arch struct level and cleanup some sym values #226
- `macho` Add raw fields for section and encryption info #238
- `mp4` Add more HEIF boxes support #221
- `mpeg_pes` Support MPEG1 #219
- `rtmp` Add Real-Time Messaging Protocol decoder. Only plain RTMP for now. #214
- `matroska` Symbol name cleanup #220
- `tcp` Better port matching and make it possible to know if byte stream has start/end. #223
- `udp` Better port matching #223

## Changelog

* 010f6430 Update docker-golang from 1.17.8 to 1.18.0
* 05096f50 Update docker-golang from 1.18.0 to 1.18.1
* e5f61e22 Update github-go-version from 1.17.7, 1.17.7, 1.17.7 to 1.18.0
* fdfc5c5b Update github-go-version from 1.18.0, 1.18.0, 1.18.0 to 1.18.1
* 4ea362e3 Update github-golangci-lint from 1.44.2 to 1.45.0
* 2a90485b Update github-golangci-lint from 1.45.0 to 1.45.2
* d9195ac4 Update gomod-mapstructure from 1.4.3 to 1.5.0
* cf88bc11 Update make-golangci-lint from 1.44.2 to 1.45.0
* 3a0799cb Update make-golangci-lint from 1.45.0 to 1.45.2
* 34cbe487 amf0: Decode strings in more detail
* b2a865ea avc_sps: Add chroma format name mapping
* b35b1804 decode,format: Add d.FieldFormatOrRaw(Len)
* f4480c6f decode,interp: Support for format specific options
* 5ff67e4c formats: Sym and field name cleanup to be more jq friendly
* 3c029925 github: Update action versions
* 02a97fa3 gojq: Rebase fq fork
* 2e240447 gojq: Rebase fq fork
* 518f6af4 gojq: Rebase fq fork
* 88f791e0 gojq: Rebase fq fork
* 8c918702 gojq: Rebase fq fork
* adde8c70 gojq: Rebase fq fork
* d79afeb3 gojq: Rebase fq fork
* dd0d97ea gojq: Rebase fq fork (speedup and fix range with JQValue)
* afd724bf gojq: Rebase fq fork. Fixes JQValue path tracking when iterating
* 74978c9d hevc: Add hevc_vps, hevc_sps and hevc_pps decoders
* c0202483 hevc_vpc,hevc_sps: Use same naming for profile as in spec
* 09385c61 id3v2: Add 2.0 PIC support
* 9cb4b57a interp,cli: Handle ctrl-c properly
* 607202bb interp: Don't truncate last display column
* 6f03471d interp: Paths with a array as root was missing start dot
* dabad850 interp: Proper display column truncate
* e8678ca8 interp: Remove opts refactor leftover
* d376520f interp: Remove to*range pad argument and fix stdout padding issue
* 087d1241 interp: Simpler and more efficient hexdump
* 21ad628a interp: dump: Show field name for compound values in arrays
* e8dc7112 ipv6,icmpv6: Add decoder
* d6c31dac macho: Add section and encryption_info raw data fields
* 5424eed7 macho: Cleanup syms and remove redundant fat_arch struct
* f8d79a57 matroska: More sym cleanup
* f34ebd83 mp4: Add more HEIF boxes
* f8fd6b7f mp4: Add more HEIF boxes
* 39ba5c4d mpeg_pes: Support mpeg1 and some cleanup
* d8aaf303 rtmp,amf0: Add decoders
* 788b0ac1 rtmp,amf0: Improve decoders, aac asc, chunk stream interrupt, fix amf0 ecma arrays
* 5d25bbc2 tcp,udp: Refactor and make port matching better

# 0.0.6

Added `macho` decoder (thanks @Akaame), nicer REPL interrupt, error and prompt, add `slurp`/`spew` functions and `explode` for binary.

Added fq talk slides from [Binary Tools Summit 2022](https://binary-tools.net/summit.html) to `README.md`.

## Changes

- Major query rewrite refactor to share code for slurp-ish functions `repl`, `slurp` and future `help` system. #178
- REPL improvements:
  - Much improved eval and output interrupt. Should fix more or less all issue with un-interruptable long outputs. It is still possible to get "hangs" if some decode value ends up being expanded into a huge string etc. #191  #192
  - Prompt paths now has colors support. #181
  - Shows an arrow on parse error.
  - Faster on multi inputs. #165
- Speedup interpreter by skipping redundant includes. #172
- gojq fork rebase: #179
  - Fixes `try ... | ... catch` precedence issue.
  - `tonumber` now supports non-base-10 numbers.
- Add `slurp`/`spew` to collect outputs and outputs them later.
  - `1,2,3 | slurp("a")` collects, later do `spew("a")` to output them. Also a global array `$a` will be available.  #178
- Add `explode` for binary. #188
  - `"a" | tobits | explode` return bits `[0,1,1,0,0,0,0,1]`.
  - `"åäö" | tobytes | explode` return utf8 bytes instead of codepoints `[195,165,195,164,195,182]`.
- Add optional sub topic to `--help`: #177
  - Replace `--formats` with `--help formats`. #181
  - Add `--help options` to see all default option values. #181
- Remove `var`, use `slupr` instead.

## Decoder changes

- `macho` Add decoder. Thanks @Akaame #43
- `mp4` Support `colr` box. #176

## Changelog

* ee5e4718 Update docker-golang from 1.17.7 to 1.17.8
* ca04cc20 Update github-golangci-lint from 1.44.0 to 1.44.1
* 5c6e1d32 Update github-golangci-lint from 1.44.1 to 1.44.2
* 1b8e6936 Update make-golangci-lint from 1.44.0 to 1.44.1
* 9d5ba826 Update make-golangci-lint from 1.44.1 to 1.44.2
* cd2cbef6 decode: Some cleanup
* 9e4f2641 dev: Add .jq-lsp.jq to add additional builtins for jq-lsp
* c6a90cfc doc,asn1_ber: Add more documentation
* c53bd777 doc: Add bts2022 video
* b97776c9 doc: Add fq bts2022 presentation
* d334c2d4 doc: Add href in supported format list
* c95b0d6d doc: Forgot make doc
* a202df9a doc: Improve and fix some typos
* 9ec1d357 doc: Improve project description
* 758b2d0e doc: Regenerate after macho merge
* 920629f5 doc: Regenerate and fix macho section size
* d3397cf9 doc: Tweak format diagram
* d47e04c4 fixup! macho: CPU_SUBTYPE_MULTIPLE and TYPE_ALL are 0xff_ff_ff_ff
* 27e76157 format: Simplify torepr, no need for _f function
* 206dcd02 fuzz: Include more testdata seed files
* be6f0093 gojq: Rebase fq fork and add support for non-10 base for tonumber
* 33efb02a interp,repl: Add path and value colors to prompt
* 41551de3 interp,repl: Improved eval and output interrupt
* dff7e7da interp: Cleanup binary regexp overloading and add explode
* fe8183b5 interp: Color parse in jq
* 6f10745a interp: Fix interrupt regression after query rewrite refactor
* f66f115e interp: Make _finally handle null, call fin once and last
* eeb59152 interp: Make help output less wide
* 9dc59e5d interp: Move _is_decode_value to jq
* 0bc11719 interp: Move opts eval to options.jq
* 3f50bb90 interp: Rework formats and options help
* 03f450f8 interp: Skip redundant includes
* c5918d23 macho: CPU_SUBTYPE_MULTIPLE and TYPE_ALL are 0xff_ff_ff_ff
* 5c974209 macho: TS string to UTC
* 04eae939 macho: add basic docs
* 5e95d1c3 macho: add cpuSubTypes
* 2638f419 macho: add darwin_amd64 test
* 5c5bd879 macho: add fqtest actualization
* bf214d5e macho: add nolint suppression to const defs
* 333a3243 macho: add scalar.Hex mapper to addr fields
* a86e7043 macho: add section type parsing
* 90b94631 macho: adopt plural-singular scheme for FieldStructArrayLoop
* b78ed02f macho: barebones decoder impl
* e199d219 macho: basic impl for ar and fat file parsing
* 66feebc5 macho: change parseFlags impl for ordered results
* b5fe9ce6 macho: change registry description
* 20e5be3f macho: delete ar decoder code
* efdd0bf5 macho: discard lc_ and lowercase command names
* b0911af2 macho: docs review changes
* a29bfca5 macho: expand filetypes and header flags
* fb0654ec macho: fix FieldUTF8NullFixedLen for segname
* d1f093ce macho: fix fat header decode bug
* 0d648928 macho: fix null in segname sectname
* 9eb71dc6 macho: generate doc via make doc
* 3991c51a macho: handle unknown lc_commands better
* ef2919b3 macho: introduce arm and fat tests
* 98c9840d macho: linting changes for ar parse
* 1feb81c9 macho: little-endian to little_endian
* 141a8e84 macho: mach_header_X to header
* 9206d9d8 macho: magicToHex to scalar.Hex
* 2021b054 macho: make actual
* 70b84cde macho: ntools fix LC_MAIN fix
* 78699f3a macho: parse flags individually
* 4016ad0b macho: parse segment section flags
* 5a48cb30 macho: refactor prebound_dylib
* 2e7767cd macho: remake docs
* 33347503 macho: reuse ar decoder
* 228757b9 macho: review fixes
* 5ee9a23c macho: review fixes
* e3daee7d macho: simplify thread state decoder
* 70c9d519 macho: thread state visualization
* a4789dc1 macho: timestamp mapper
* 2ccb8087 macho: update test cases v to dv
* 74abe990 macho: update tests
* 12eb7cc5 macho: use FieldUTF8NullFixedLen
* 5f4ad410 macho: use FieldUTF8NullFixedLen for segname
* f8690e6c mp4: Add colr box support
* b157751a mp4: Reformat and use dv in test
* 0a043f90 repl,interp: Refactor repl and slurp
* ca8cdadb repl: Add comments and query from/to  helper
* 9cb4205b repl: Correct error arrow position in color mode
* e238f292 repl: Speedup multi input to sub-repl
* 56ae4a0c test: Make expect cli test more robust

# 0.0.5

Improved binary slicing and bit reading, `avro_ocf` decoder (thanks @xentripetal), `asn1_ber` decoder, renamed `display` aliases, new `grep_by` and `paste` function.

## Changes
- Big internal bit reader refactor. Now much more consistent code and fixes some issues reading and decoding of binary arrays and binary slices. #123
  - Bit reading and IO have been moved to a `bitio` package.
  - Non-simple bit reading have been move out of `bitio` to `decode` package.
  - `[0,1,1,0,0,1,1,0,0,1,1,1,0,0,0,1 | tobits] | tobytes | tostring` returns `"fq"`.
  - `[.frames[0], .frames[-1]] | mp3` decode mp3 based on first and last frame of other mp3.
- Add `grep_by` that recursively selects using a filter condition and ignores errors. #102
  - `grep_by(.type == "trak")` finds all objects where `.type` is "trak" (all mp4 track boxes).
  - `grep_by(tonumber | . >= 40 and . <= 100)` find all numbers between 40 and 100.
  - `grep_by(format == "jpeg")` find all jpegs.
- Add `paste` function to read string from stdin util ^D (EOF). Can be used to paste in REPl etc. #143
  - `paste | frompem | asn1_ber | repl` wait for PEM encoded text on stdin (paste and press ^D), decode it with `asn1_ber` and start a sub-REPL with the result.
  - `paste | fromjson` decode pasted JSON.
  - `eval(paste)` eval pasted jq expression.
- Cleanup display aliases. Remove `v` and `f`, add `da`, `dd`, `dv` and `ddv`. #112
  - `d`/`d($opts)` display value and truncate long arrays and buffers
  - `da`/`da($opts)` display value and don't truncate arrays
  - `dd`/`dd($opts)` display value and don't truncate arrays or buffers
  - `dv`/`dv($opts)` verbosely display value and don't truncate arrays but truncate buffers
  - `ddv`/`ddv($opts)` verbosely display value and don't truncate arrays or buffers
- Refactor `radix` into `toradix($base)`/`fromradix($base)`. #139
- Remove `number_to_bytes`. Can be done with `tobytes`. #139
- Change `tobytes` to zero pad most significant bits to byte alignment if needed. #133
- Add `tobytes`/`tobits` variant that takes an argument to add extra padding. #133
  - `0xf | tobytes` 8 bit binary with last 4 bits set
  - `0xf | tobytes(4)` 32 bit binary with last 4 bits set
  - `0xf | tobits(12)` 12 bit binary with last 4 bits set
- Rename fq type buffer to binary as it makes more sense. #133
- Add `topem`/`frompem` to work with PEM encoding. #92
- Add Windows scoop install. #137 Thanks @thushan
- Add `display`, decode value, binary and binary array documentation. #118 #136 #133
- Add decode API documentation. #149
- Improved REPL completion for keys. #144
- Add `-o force=<bool>` option that sets force decode option. Same as `mp4({force: true})`. #126

## Decoder changes
- `avro_ocf` Add decoder. #38 Thanks @xentripetal
  - Full avro OCF support. Handles all primitive, complex, and logical types besides decimals.
  - Able to handle deflate, snappy, and null codecs for blocks.
- `asn1_ber` Add decoder. #92
  - Also decodes CER and DER (X.690) but with no additional validation at the moment.
  - Support all types but real type is currently limited to range for 64 bit integer/float.
  - Has `torepr` support.
  - No schema support.
- `aac_frame` Only decode object types we know about. #130
- `mp3` Shorter sync find heuristics. #104
- `mp4` Add `stz2` support
- `mp4` Add `pnot` (preview container) and `jP  ` (JPEG 2000) signature. #125

Also thanks to @Doctor-love for keeping things tidy.

## Changelog

* 6fc1efd9 Add test case with all data types
* ae4a6243 Adds Windows Scoop instructions for fq.
* 4b809a73 Change avro codec to funcs
* 66ca1f10 Change tests to use new verbose syntax
* 7345b8c7 Cleanup
* 07ddf36f Cleanup for linting
* 4508241b Cleanup snappy
* 0909fb6d Comment on snappy decompression
* 21cfc70c Dates need to specify UTC too
* 75b84961 Fix lint
* 7a8e3ca2 Hook into registry, add codecs
* 251053ef Initial pass on logical types
* 2605bce4 Lint and add basic doc
* ee184075 Parse header using avro decoders. Still not certain this is the best idea. Will get opinions before finalizing.
* ab50088d Polish of problem template and clarifying questions
* 27789f2d Regenerate docs
* 5a1d35e7 Remove redundant question and fix typo
* 31c4c0d3 Support snappy and deflate codecs
* 0300c955 Take heading off doc to match make doc format
* 6f57cdbf Timestamps should be UTC
* 06085a26 Undo change to doc/file.mp4. I have no idea how this got changed in the first place? Maybe some macos shenanigans.
* d137a72a Update docker-golang from 1.17.6 to 1.17.7
* 267e30ec Update github-go-version from 1.17.6, 1.17.6 to 1.17.7
* 1e859cda Update github-golangci-lint from 1.43.0 to 1.44.0
* 16849c8f Update linting
* d02d8968 Update make-golangci-lint from 1.43.0 to 1.44.0
* 68e85a2d Use existing scalar description helper
* 3bab3d65 aac_frame: Only try decode object types we know about
* 0829c167 asn1_ber: Add decoder
* 0312c92c asn1_ber: Add more doc and multiple outputs for frompem
* 06245d12 binary,decode,doc: Rename buffer to binary and add some documentation
* 7c521534 bitio,decode: Refactor bitio usage and make buffer slicing more correct
* 0d74e879 bitio,doc: Even more cleanup
* d854ed57 bitio: Cleanup documentation a bit
* 82aeb355 bitio: More doc cleanup
* 01ecde64 bump: Add snappy config
* de64a99e cleanup some docs, change enum to mapper, error zigzag on more than 8 bytes
* 6cd1c38f decode,scalar: Add scalar.Str{Uint/Int/F}ToSym to parse numbers
* 4ab6381d decode: Add scalars args to FieldRootBitBuf
* be71eb01 decode: Rename LenFn() to LimitedFn, add FramedFn and document
* 7bc25219 doc,interp: Add some example usages to cli help
* 8e47fb1a doc,matroska: Fix filename in example
* c15f5283 doc: Add format links to format table
* b86da7ae doc: Add initial decoder API documentation
* 49c90f89 doc: Add macOS security notes and move supported format up a bit
* 06b67e4b doc: Add more license details
* a8664ed5 doc: Add per format documentation
* 09552628 doc: Add snappy license
* 36307857 doc: Cleanup and add more decode value and binary documentation
* 710c29b2 doc: Color edges in diagram based on dest
* f0ce7179 doc: Document display and some more jq hints
* b3504680 doc: More decode API details and polish
* 6b51b067 doc: More display alias leftover fixes
* dd3e40fb doc: Unbreak formats_digaram.jq since radix change
* c52a1a23 doc: Use f($a; $b) instead of jq f/2 notation
* 233d86a3 fq: Add arch and os to --version
* b8efd8e5 fuzz: Fuzz all formats
* e1bdfdf8 fuzz: List seed numbers and make it build again
* 6090b65e fuzz: Make it compile again and run one format per fuzz
* aea48847 github: Add basic issue template
* b55ca2cd gojq: Rebase fq branch
* 47c978e4 goreleaser: Use zip for macos
* 85371173 id3v2: Should assert not validate magic
* d6ca4818 initial work for avro OCF files
* ca68e6a1 interp: Add Platform() method to OS interface
* e792598c interp: Add grep_by/1 to recursively match using a filter
* 0a1a5610 interp: Add missing default opts for tovalue
* 48a19cb8 interp: Add paste function to allow pasting text into REPL etc
* fc0aacb6 interp: Cleanup display aliases, now: d, da, dd, dv, ddv
* bf7a483f interp: Fix handling of group decode error from stdin
* 26d9650b interp: Refactor radix* into toradix($base)/fromradix($base)
* 366f6b18 interp: Support force decode as -o force=true
* 77ab667b interp: Use absolute path in errors
* c31ec2a3 interp: Use correct sym color
* 898dfec1 lint: Fix typeassert and case exhaustive warnings
* d5401166 make doc
* bf170be8 make: Cleanup some not very used targets
* 8d2d88f4 mp3: Decrease max sync seek length between frames to 4k
* d555c324 mp4,fuzz: Fatal error on infinite sgpd box entries
* 45b00aab mp4: Add stz2 support
* 092609be mp4: Add video preview (pnot) and JPEG 2000 (jP) signatures
* febce5a5 mpeg_spu: Fatal error on infinite loop
* c58ba28d mpeg_spu: Fatal error on unknown cmd
* d1943dad pcapng,fuzz: Fix infinite loop by fatal error on block length <= 0
* 2ab395a0 protobuf: Add note about sub message decoding
* af053811 repl,interp: Make stdio work during completion
* bd9be2c5 repl: Fix completion of non-underscore extkeys
* 69c745d3 simplify scalar usage
* 778a1a41 zip: Assert signature not validate

(Some commits have been removed from list for clarity)

# 0.0.4

## Changes
- Add a `torepr/0` function that converts decoded value into what it represents.
Initially works for:`bencode`, `cbor`, `bson` and `msgpack` #74
Example usage:
`fq torepr file.cbor`
`fq -i torepr file.cbor`
`fq torepr.field file.cbor`
`fq 'torepr | .field | ...' file.cbor`
- Add `stderr/0` function for jq compatibility #58
- Bitwise operators `band`, `bor` etc are now normal functions instead of operators. Was done to be syntax compatible with jq #61
Uses the same convention as jq math functions: Unary uses input `123 | bnot`, more than one argument all as arguments `band(123; 456)`
- Decode API now supports null values #81
- Decode API now supports arbitrary large integers #79
- TCP reassembly now supports streams with missing SYN/ACK #57
- Update readline package to version with less dependencies #83
- Make REPL prompt more jqish #70

## Decoder changes
- `bencode` Add decoder #64
- `cbor` Add decoder #75
- `msgpack` Add decoder #69
- `mp4` Much improved sample decode #82 #78
- `png` Decode PLTE and tRNS chunks #59
- `tar` Don't assume there is a end marker and support more than 2 blocks #86 #87

Also thanks to @Doctor-love for keeping things tidy.

## Changelog

* af8e7ef bencode: Add decoder
* 0b0f28e cbor: Add decoder
* 1383b41 decode,interp: Add arbitrary large integer support (BigInt)
* 548a065 decode,interp: Finish up nil value support
* ff5c0b8 decode: Error on negative number of bits when reading numbers
* cf8a50c decode: Use stable sort for values to not change order or values with same range start
* b4694b6 doc,dev: Add some more decoder implementation help
* 0c1716b doc: Add alpine and go run
* 809210b doc: Add more dev tips
* 59b8803 doc: Document dev dependencies and related PRs/issues etc
* 6ca4767 doc: Improve formats graph a bit
* 8e9700d doc: Improve readme a bit and add torepr example
* 0cf486d elf: fix all-platforms naming typo
* 263f1ae flac: Don't allow zero subframe sample size
* 729a6ca formats: Sort and make lists less likely to cause collision
* 78c0775 fq: Embed version in source
* aa7adb6 fq: Update version to 0.0.4
* 7461264 fuzz: Skip other tests when fuzzing
* be0ef80 interp,fq: Make bit operators normal functions
* a3cfcd0 interp: Add stderr again for jq compat
* 149cb3f interp: Add torepr/0 that converts decode value into what it reptresents
* b3a0980 interp: Document bit opts funcs and add some error tests
* 8d10423 make: Fix quote issue in release script
* 4a1e859 mp4: Improved stsz handling
* 61bf2ce mp4: Refactor sample decode into something more sane
* a6bf62c msgpack: Add decoder
* edad481 num,mathextra: Rename num package to mathextra
* bfc977b png: Decode PLTE and tRNS chunks and cleanup syms a bit
* 36d2891 readline: Update to version with less deps
* 9770b00 repl: Make prompt for array and iter more jqish
* ba1edef tar: Allow more than 2 zero end blocks at end
* 5921d76 tar: Don't assume there is a end marker
* edd0ae1 tcp,flow: By default allow missing syn/ack for now

(Some commits have been removed from list for clarity)

# 0.0.3

## Changes
- Now works on Windows #52
- Now builds and runs on 32-bit CPUs #30 @danfe
- `print/0`, `println/0` function now properly convert input to string if needed. #54
- `match` functions now don't try to be smart with buffers, use. `tobytes` etc instead. Less confusing and also unbreak `grep`:ing decode values. #53
Now this works: `fq 'grep("^strtab$") | parent | {name, string}' /bin/ls`
- Add Arch package #19 @orhun @ulrichSchreiner @dundee
- Add Nix package #22 @siraben @jtojnar @portothree
- Add FreeBSD port @danfe

## Decoder changes
- `bson` Add Binary JSON deccoder
- `ar` Add  Unix archive decoder
- `bsd_loopback_frame` Add BSD lookback frame decoder (used in pcap files)
- `elf` Now does a two-pass decode to correctly handle string table references
- `elf` Decode more sections: symbol tables and hashes
- `matroska` Assert sane tag size only for strings
- `pcap` Don't fail if incl_len > spanlen

Also thanks to @Doctor-love @mathieu-aubin for keeping things tidy.

## Changelog

* 628f0f4 bson: Add decoder
* 46b59d0 crc: Unbreak build on 32-bit arch
* 681dbc2 elf,ar: Add ar decoder, improved elf decoder
* e5c620d github,ci: Add windows, macos and 32-bit linux
* 52dddbb goreleaser: Use draft release to allow release note changes
* e365f22 interp: Cleanup stdio usage and functions
* 55b1d5c interp: Move _registry to decode
* b6515c8 interp: Remove buffer smartness for regexp match functions
* b867113 matroska: Assert sane tag size only for strings
* b9aef39 pcap,pcapng,bsd_loopback_frame: Add decoder, refactor link frame into a group
* af23eb8 pcap: Don't fail if incl_len > spanlen
* a41f0d4 windows: Correct @builtin include path join
* bf9e13c windows: Unbreak tests

(Some commits have been removed from list for clarity)

# 0.0.2

## Changelog

* 00f34c2 Update docker-golang from 1.17.3 to 1.17.4
* 05b179c Update docker-golang from 1.17.4 to 1.17.5
* c721ac8 Update github-go-version from 1.17.3, 1.17.3 to 1.17.4
* befe783 Update github-go-version from 1.17.4, 1.17.4 to 1.17.5
* 75aa475 decode: Generate Try?Scalar* methods for readers too
* 9b683cd decode: Cleanup some io panic(err)
* e1e8a23 doc: Add color/unicode section and move config section
* ee023d7 doc: Add some more related and similar projects
* d02c7c4 doc: Add some more usage examples
* be46d5f doc: Cleanup todo and add some dev notes
* 47deb4d doc: Fix interpretation typo
* fe68b51 doc: Improve readme text a bit
* 916cb30 doc: Improve usage examples as bit
* 5cb3496 docker: Fix broken build, copy fq.go
* c2131bb flac: Cleanup scalar usage and fix incorrect sample rates
* 1500fd9 gojq: Update to rebased fq fork
* 3601fe3 gzip: fuzz: Don't uncompress on unknown compress method (nil create reader fn)
* f4f6383 interp: Add ._index for values in arrays
* 4558192 interp: Cleanup buffer code and implement ExtType()
* e823475 interp: Fix help a bit
* fa350c6 interp: Move display to jq
* f7c7801 interp: Rework buffer regex support
* 6ed2e2e interp: dump: Indicate arrays using jq-syntax
* 9aec91a interp: match: Fix issue with regexp meta characters when matching using a buffer
* e5e81e7 make: Fix prof build issue
* e91b22b matroska,ebml: Use scalar and require sane tag size
* eb9698f mp4,ctts: Seem more usable to treat sample count/offset as signed
* c149732 mp4,trun,fuzz: Limit number of constant sample entries
* 406263b mp4: Add comment about hdlr.component_name prefix byte
* 28a3b71 mp4: Cleanup sample decode code
* 6278529 mp4: fuzz: Make sure stsz has sane number of entries on constant sample size
* 9f08af3 mpeg,aac: Factor out escape value decoding
* dc1aea3 opus: Cleanup endian usage and fix incorrect preskip decode
* 2b2320d pcap,flows: fuzz: Handle broken packets more nicely
* 1d7ace3 pcap,pcapng,tcp: Use capture length not original length
* b525d0b pcap: fuzz: Skip ssl2 packet if too short
* 713ffe4 scalar: Add Require* and Require/Validate/AssertRange*
* f348002 sll2: fuzz: Limit address length to max 8 bytes
* eb4718f tar: Cleanup api usage
* 45026eb tar: Cleanup constant usage a bit more
* 91cc6d8 tar: Fix size decode regression after cleanup
* 07a2ebe tiff,fuzz: Fatal error on infinite ifd loops
* 91217e8 tiff: Fix endian typo and cleanup todos
* 3850968 udp: Use proper udp payload format var name
* 6a8d77b vorbis_comment: Cleanup endian usage and naming a bit
* 57e9f41 vorbis_packet: Cleanup endian usage
* f5a4d26 vscode: Use tabsize 2 for jq files
* dd883b3 wav: Cleanup endian usage
* e260830 webp: Cleanup endian usage



# 0.0.1

## Changelog

* 4242bf6 *_annexb: Refactor into avc/hevc_annexb
* e86b45b Add *grep/1/2 and find/1/2
* 36fd74a Add comment how raw byte regexp matching works
* d1be167 Add decode struct each order test
* 7f36f70 Add to/bytes/bits[range]
* 571bf29 Change project title
* 95ec5e1 Cleanup and rename s/BufferView/BufferRange
* 9797cdc Cleanup dev/snippets.jq
* 565f18d Correct avc_au format variable name
* 6a1fa04 Decode hvc1 as hevc samplesa also
* 74bad2d Fix broken value.fqtest
* 514739a Give proper error on missing short flag
* 57f0ec1 Improve cli help a bit
* 9704659 Init
* f33b310 Refactor decode.Value gojq bindings
* 3d90b6d Remove fixed comment
* e4e269b Rename and move cli test to pkg/cli
* 1bd34bf Same args error behavior as jq
* 3693667 Start of configurable json bit buffer formats
* 590ee52 Update ci-golang from 1.17.0 to 1.17.1
* 72a3f69 Update ci-golangci-lint from 1.42.0 to 1.42.1
* 078cf29 Update docker-golang from 1.17.1 to 1.17.2
* f5a8484 Update docker-golang from 1.17.2 to 1.17.3
* 1371bc7 Update docker-golangci-lint from 1.42.1 to 1.43.0
* f1ad700 Update github-go-version from 1.17.1, 1.17.1 to 1.17.2
* 314bd17 Update github-go-version from 1.17.2, 1.17.2 to 1.17.3
* 53a8d91 Update github-golangci-lint from 1.42.1 to 1.43.0
* 898cd26 Update golang from 1.17.0 to 1.17.1
* 31cd26e Update golangci-lint from 1.42.0 to 1.42.1
* d4b2d58 Update gomod-mapstructure from 1.4.2 to 1.4.3
* 1decf85 Update make-golangci-lint from 1.42.0 to 1.42.1
* 233aaa1 Update make-golangci-lint from 1.42.1 to 1.43.0
* 15e9f6f ansi: Correct background reset code
* 1ccab2d apev2: Add test
* 27e4770 apev2: Fatal if > 1000 tags
* 3bf1a57 avc: Cleanup and add color names etc
* 777191f avc: Correct sign expgolomb decode
* 20021f4 bitio: Handle < 0 nbits
* 7c4b0b3 bitio: Simplify by embedding reader
* f600f2e build: Require go 1.17
* 77f97aa builtin: Add chunk_by, count_by and debug
* 27ba359 bump: Add action and cleanup names
* 71e87e6 bump: Cleanup config, add config for release.yml
* 8f2f524 bump: Move bump config to where it's used
* 184df0a cli: A bit clearner array and iter prompt
* d350971 cli: Add --decode-file VAR PATH support
* 17104f0 cli: Add --options to make help a bit nicer
* 6356a84 cli: Add --raw-string
* 93fd097 cli: Add -M -C support and default to color if tty
* 8dc0f06 cli: Add completion tests
* 2010cac cli: Add error test
* 706b2f2 cli: Add exit 2 (like jq) for no args
* 569b631 cli: Add output join tests
* 49f541c cli: Add proper repl iterator support
* 3304f29 cli: Add string_input options test
* 0dd848d cli: Better filenames in errors
* 22eb53d cli: Cleanup MaybeLogFile
* c14c29a cli: Cleanup and more comments
* 10d7ed7 cli: Fix error filename on script error
* 2b8d11d cli: Fix indent
* 2d4eb9c cli: Fix non-string variables and var(; f) variant to delete etc
* f5ffd32 cli: Implement --arg, --argjson and --rawfile
* ae5566a cli: Include paths and some refactor
* ec98fd3 cli: Make --argjson and --decode-file error similar
* ac8cfca cli: Make --raw-string work with input/0 and inputs/0
* b33f2cd cli: Make profile build optional and move it to cli
* e2ff2a2 cli: Move help/0 to inter.jq, better help for -n
* d0bb9a5 cli: Nicer grammar for --null-input help
* 78eb737 cli: Nicer usage and indent input iteration
* 394e2b3 cli: Only compelete at end or whitespace
* 8d1fafe cli: Only show fq info for --help
* 2684ed2 cli: Prepare completion for better variables support
* e666380 cli: REPL and multiple files
* 75cf46f cli: Refactor options code
* f8ab00e cli: Remove unused eval debug arg
* 5c8fb5f cli: Revert accidental history path change
* 55cd45f cli: Simplify code
* 2874bc7 cli: Unbreak colors in windows
* 48517c7 cli: Unbreak part of completion
* 518b725 cli: Use format/0 to check if value is a format
* 21bef18 cli: User defined global vars
* 538f4ff cli: add -nul-output compat
* 85d1719 cli: jq compat, multiple -L
* c8f0264 cli: more jq compat
* f893295 cli: rename fq.jq to interp.jq
* 1436fdc completion: Better  and _internal handling
* c7416e6 decode, interp: More buffer reuse
* 6ee7977 decode,format: Allow root array
* 58ba84f decode,interp: Add RecoverableErrorer interface instead of enumerate
* b66ed32 decode,interp: Make fuzzing work again and cleanup fatal/error code
* 5052bae decode,interp: Refactor to allow decode/fillgap a range
* d4142b8 decode,png: Add FieldFormatReaderLen, refactor out zlib to format
* 826c509 decode: Add Generated header for scalar
* 3db11d3 decode: Add UTF8Fn functions and trim some null terminated strings
* c083a9e decode: Fix MapRawToScalar regression
* c17483d decode: Fix accidental rename
* 6fba1a8 decode: Fix bitbuf root handling a bit
* d1e1cd9 decode: Fix walk root depth issue causing dump to indent incorrectly
* 1b32b42 decode: Major decode API refactor
* 7af191d decode: Move io helper into *D
* 26d615b decode: Move name/description into DecodeOptions
* 986d5ec decode: Move registry package to decode/registry and add a format group type
* ede2e77 decode: Nicer scalar template and add doc
* 6207fcc decode: Pass context to be able to cancel properly
* 5d98a69 decode: Refactor Error/Fatal into printf functions
* 9f55b6e decode: Refactor and add symbol struct and mapping
* 2fc0a71 decode: Refactor scalar usage
* 8eaba88 decode: Refactor walk code a bit, add WalkRoot* to stay inside one root
* f40320b decode: Remove D.Scalar* and add d.(Try)FieldScala*Fn instead
* c155c89 decode: Rename format *Decode to *Format
* f801cc0 decode: Rename s/FieldTryFormat/TryFieldFormat for consistency
* 776a6b3 decode: Reuse read buffer per decode to speed things up
* 9d116df decode: Rework use of TryFieldReaderRangeFormat
* d48ebc1 decode: Simplify Compound.Children
* 473b224 decode: Simplify and move format arg into DecodeOptions
* 0480a2f decode: Some format decode and sub buffer work
* a49e924 decode: Use golang.org/x/text for text decoding
* af3e6b1 dev: Add format_summary.jq
* 4c6de82 dev: Add summary how start to dump tree works
* a926c8f dev: Some debug notes
* b17a715 difftest: Remove accidental log
* 1e1ad14 difftest: Sync code
* 3cea849 dns: Cleanup a bit
* d469edf doc,make: Correctly strip out graphviz version from svg
* 15d85e1 doc: Add find/1 find/2
* 26ea6d8 doc: Cleanup
* a131210 doc: Cleanup and add note about . argument
* 0a97f86 doc: Cleanup and note about repl limit
* 8440e8a doc: Cleanup todo and dev a bit
* 15b6d64 doc: Cleanup up a bit
* 1047d90 doc: Document io packages a bit
* 1a0089e doc: Fix typo and some improvements
* 97c7403 doc: Fix typo in README
* 0814206 doc: Fix typos and old examples
* 0e8c82a doc: Improve readme goals a bit
* 2f9d93d doc: Improved readme and cleanup todo
* 242525f doc: Move formats to own file
* 583bc38 doc: Regenerate demo.svg
* a050adc doc: Regenerate demo.svg
* 18e3e20 doc: Regenerate demo.svg with newer ansisvg
* 1f61704 doc: Regenerate svg after ansisvg monospace update
* a7459b3 doc: Some basic usage and cleanup
* 6fa5ae8 doc: Some fixes
* 00b7c18 doc: Some more doc work
* 07c7daa doc: Some rewording
* 5a82224 doc: Update
* 21a74fa doc: Update
* d6d3265 doc: Update README a bit
* df5bd19 doc: Update TODO
* fb13fe5 doc: Update todo
* eee3c4e doc: Use unicode pipe in demo
* b9b0326 doc: fq - jq for files
* 790198c docker: Make image build again and install expect
* 2387ec8 docker: Use golang 1.17
* f8e5944 dump,json: Properly figure if compound or not
* ee972f4 dump: Add ascii header
* 69c6d15 editlit.jq: Update after mp4 field renames
* 27095b5 es: Decode MPEG ASC if audio stream
* 761c411 exif: Add note about JPEGInterchangeFormat
* 344f628 flac: Calculate correct md5 when total samples count is zero
* ce044ba flac: Cleanup
* 2af08da flac: Fix block_size regression
* 4f0bf92 flac: Make md5_calculated be a buffer
* 797bd4d flac: Refactor flac_metadatablock into flac_metadatablocks
* 1f26d4f flac_frame: Correctly read escaped samples and also a bit less allocations
* bc9951c flac_frame: Fail if trying to decode outside block size
* 723542a flac_frame: Make utf8Uint reader more correct and robust
* b643e22 flac_frame: Support non-8 bit align sample size
* 2daa738 flac_frame: Use d.Invalid for possible errors
* c1d9b4d flac_metadatablocks: typ > 127 can't happen, add app id
* 984ba1a flac_metadatablocks: type >= 127 is invalid
* 2b35d28 flac_picture,mp4: avif images and data fallback is image format fails
* 509b8f8 flac_picture: Add picture_type names
* 4f8d037 format,decode: Some crc and endian refactor
* aa38ccf format,interp: Use MustGroup and add probe order test
* 577c0f5 format: Add panic if register after resolve
* 606c0b6 format: Add vorbis-comment-picture test, add .gitignore and cleanup
* 798141a format: Cleanup comments
* c0eebcc format: Remove unused ProtoBufType
* 4b48828 format: Rename source file to match format name
* d1b514e format: Some claeanup
* ec97eca format: Split default.go into format.go and shared.go
* 25f5ad7 fq,cli: Rename chunk  to streaks, cleanup
* 905c0ab fq: Add chunk/1
* 46d37ef fq: Add cli sanity test
* b849895 fq: Add truncate array support to dump/display
* 8cb380e fq: Generate decode alises code
* dfcefc1 fq: Make format/0 native for performance
* ba273be fq: Make relative include work with @builtin etc
* d23edaa fq: Rename bits/0 bytes/0 to tobits tobytes, remove string/0
* a7a58c8 fq: Rename main.go to fq.go
* 834f4a5 fq: use jq functions for all display alises
* 31d7611 fqtest: Add env support and isterminal/width/height support
* 285356d fqtest: Cleanup and dont assert when WRITE_ACTUAL
* 86b34a3 fqtest: Fix section regexp
* dee10db fqtest: No need to escape empty stdout
* 6916880 fqtest: Refactor our script part to own package
* 01d8a90 funcs: Add delta/0 delta_by/0
* 96f7a75 funcs: Fix typo add count/0
* 962d84d funcs: Make intdiv truncate to int
* 08ec4f0 funcs: Remove unused string function
* d5c084c funcs: chunk_by comment
* 80eaa46 funcs: format helper
* c770a75 funcs: make in_bytes/bits_range more generic
* 6ff5aca gif: Support GIF87a
* 0d87018 github: Install expect to cli test
* 8e3d9d8 github: bump: Checkout with bump token so it's used when push
* 0c7fa09 gojq: Initial update support
* 1888bb2 gojq: Remove div operator
* d7dbe7c gojq: Update fq branch to fix mod (-1 % 256) difference
* bf5c222 gojq: Update fq fork
* 093ee71 gojq: Update fq fork
* 1d15c1d gojq: Update fq fork
* 3044fef gojq: Update fq fork
* 5a5f01e gojq: Update fq fork
* bfec366 gojq: Update fq fork
* 4104a18 gojq: Update fq fork
* 31a5047 gojq: Update fq fork
* 0b90558 gojq: Update fq fork
* 8277b79 gojq: Update fq fork, support JQValue alt //
* a34784d gojq: Update fq gojq fork
* 7ad3d25 gojq: Update fq gojq fork
* f5164ee gojq: Update fq gojq fork
* ac7568d gojq: Update gojq fq fork
* f828ae1 gojq: Update gojq fq fork
* 845bc6b gojq: Update gojq fq fork
* d365ab7 gojq: Update gojq fq fork
* 6420928 gojq: Update gojq fq fork
* c8776ab gojq: Update gojq fq fork
* f4cd7bf gojq: Update jq fork
* b75da30 gojq: Update jq fork
* 03af2b5 gojq: Update rebased fq branch
* c3b7b5c gojq: Update to rebased fq branch
* 2d6573d gojq: Update to rebased fq branch
* fdb811e gojqextra,decode: Add generic lazy JQValue
* aab32cf gojqextra,interp: Add lazy string to speed usage of decode value buffer where string is not used
* 9035278 gojqextra: Move errors to own file
* 86092e6 golangci: unused: Assume go 1.17
* ac106d9 goreleaser: Cleanup
* ee611a4 gzip,bzip2: Calculate CRC
* 5344c7e icc: Add mluc support
* 50d00d8 icc_profile: Add element alignment bytes field
* 59e7fd4 icc_profile: Clamp alignment to end instead of check last tag (might not be last in buffer)
* a5b802b icc_profile: Correct alignment byte count
* 57a1207 id3v2: Support GEOB tags
* f55b1af inet: Add tcp and ipv4 reassembly
* 26c594f input: make -R mimic jq on io error
* d9b45ba internal: _global_var returns new value instead of _global_state
* f9f8660 interp,decode: Add force option to ignore asserts
* 6a15625 interp,decode: Refactor out Scalar from Value and merge Array/Struct into Compound
* cc5f405 interp,format: Update tests after decode refactor/tosym/toactual
* 9f2dddc interp,gojqextra: Make buffers values even more lazy and error early on non-scalar calls
* bf19588 interp: Add ansi helper
* 7298a4c interp: Add buffer match support to find and grep
* 0d693aa interp: Add line between usage and args help
* c997536 interp: Add root, buffer_root, format_root, parent and parents
* 69e4eea interp: Better error if format/group is not found
* 7423f45 interp: Cleanup output types
* 8d442b8 interp: Cleanup stdin reading and add more option tests
* b641c77 interp: Cleanup unfinished/broken preview
* 67898cb interp: Cleanup, use BufferRange for _open, progress for all decode
* 0660ff0 interp: Clear up confusing --rawfile (add a jq alias)
* 3fafee8 interp: Clearer help for -d
* 879bb56 interp: Correctly check if _decode_value, add more has/1 tests
* 2e964fa interp: Disable progress after decode is done
* 80b5d66 interp: Document inputs and repl/cli details
* 13fae09 interp: Don't print context cancel
* 96cc128 interp: Eval options in jq instead of calling jq from go
* 178032e interp: Fix $opts shadowing in decode
* 4eccb1e interp: Fix broken aes_ctr, should return buffer instead of []byte
* 110c86b interp: Fix broken dynamic width/height
* 80a6997 interp: Fix file leak and always cache for now
* cfdd922 interp: Fix prompt issue with format
* 93322bc interp: Implement format/0 in jq
* 8e5442f interp: Limit how often decode progress fn is called
* f1fcbe5 interp: Make has/1 work for _ext keys
* 826c8bd interp: Make include abs path work again
* 6034ad7 interp: More sure stdOS stops the signal forward gorutine
* 6cacc9b interp: Move *CtxWriter to ioextra, some comments cleanup
* afb1050 interp: Move _decode_value to value.jq
* 3e7e133 interp: Move formats func def to jq
* 5cd5633 interp: Move jq function impls closer to where they belong
* 996be0f interp: Move more options code into options.jq
* 976e992 interp: Move progress logic to jq
* 528e6b9 interp: Refactor and use mapstructure
* 472c1ef interp: Refactor out string_input to own function
* 54e121c interp: Refactor repl inputs a bit
* 618c1ea interp: Refactor/Rename BufferView
* d6d9484 interp: Remove --options, probably just confusing
* b024316 interp: Remove accidental extra space after path in prompt
* ff143d0 interp: Remove redundant decode arg
* abcecb8 interp: Remove unused []byte type
* 01a407a interp: Rename s/bufferrange/buffer and cleanup
* 567bc4b interp: Rename to* funcs to match
* 36e5562 interp: Rename value.* to decode.* as it makes more sense
* 0cce5ec interp: Reorganize, move out repl and options, more functions to funcs.jq
* eedfd16 interp: Replace find with overloaded match that support buffer
* 58bf069 interp: Return []byte value as a buffer for now
* 1c3c65b interp: Rework buffer, still confusing
* 97f7317 interp: Rework repl prompt code and fix some whitespace issues again
* 4af5739 interp: Rework string/buffer for decode values
* ff2077b interp: Simplify Function, aliases done in jq now
* 1fe5d95 interp: Some better naming and typos
* 1b7b2f9 interp: Use gojqextra.NonUpdatableTypeError
* 6ce4ba9 interp: Use snake_case for all options
* 9cba69e interp: Use todescription in tests
* cf26b1f interp: _readline: Use _repeat_break, add test
* 07b4210 interp: add topath/0 and make todescription return null if there is none
* 16d1f45 interp: find buffer should always use ByteRuneReader
* 3ff0c9b lint: Enable errcheck adnd revive
* 13e98d4 lint: Fix unused bufferRange and toBufferView
* 13d5cbd lint: Remove usused nolint, should somehow tell about decode.Copy
* ad850d6 lint: want explicit types in gojqextra.go
* 90c19c6 make,test: Move testjq to own script and reuse fq binary
* 052b9c0 make: Build with -s -w same as goreleaser
* 1557e14 make: CGO_ENABLED=0 for static build
* 8019212 make: Enable -race for tests
* 263a77f make: Mark actual and cover as phony
* b0694f5 make: Move build flags to vars
* 75b59db make: Move doc generate to helper script
* fcbfc29 make: Move testjq.sh into pkg/interp
* 5bf4bc7 make: Rename testwrite to actual
* 40e26e8 make: Rename to testjq
* ab8080f make: Reorganize test targets
* 98a8bae make: Silence git if no repo version found
* 95b9c32 make: doc/formats.svg: Ignore graphviz version to get less diff
* 30ad643 make: go build args should not be quoted
* f1507f7 mod: Use proper path and dont use replace
* c4a3120 mp3,README: some cleanup
* da386ea mp3: Be more relaxed with zero padding, just warn
* 1325e5c mp3: Cleanup comments
* e104748 mp3: Continue try find frames on error
* 046f2fd mp3: Don't allow more than 64k between frames
* 527f917 mp3: Error if > 5 unique header configs
* 4344b62 mp3: Only look for supported mp3 sync headers
* 8a4f66f mp3: Probe order after formwats with raw samples and similar sync headers
* 60df04b mp3_frame: Only supports layer 3, fixes some probe issues
* ed21f36 mp3_frame: Rename samples_per_frame to sample_count
* db586eb mp4,matroska: Add *_path/0 variant that uses format_root
* 9ac17bd mp4: Add comment about future truncate to size option
* 2e71fa1 mp4: Add smhd box
* 278e909 mp4: Add tapt, prof, enof and clap boxes
* 161dcaf mp4: Better fragmented mp4 support
* e47888e mp4: Fix 64bit size regression
* 0801882 mp4: Fix field name typo for sample_composition_time_offsets_present
* f322e78 mp4: More _time decoding
* 6b8d26c mp4: Properly decode tfra
* f2c1327 mp4: Use ISOBMFF naming and some more tests
* 589207d mp4: Use descriptor field for all descriptor boxes
* c7d45ff mp4: add pssh_playready format
* e6cb708 mpeg: Nicer sym and description
* f4b11b4 mpeg_annexb: Add format
* 6a8ba31 number_to_bytes: Force int to make it work with bigint
* 88eade9 ogg: Add flac support
* 7f76986 ogg: Cleanup bitio in format out, maybe later
* 7b7faaf pcap: Add pcap, pcapng, ether8023, ipv4, udp, udp
* fc76907 png: Add proper color type
* 5c733ad readline: Update fq fork
* 4cee498 readme: Nicer demo
* ffb5adf registry: Move to pkg/registry, feels better
* dcceaa4 repl: Fix help and error message a bit
* 44d8b66 repl: Give error if repl is used in non-repl mode
* 1d0ebb5 repl: Handle directives, add tests
* c9777aa repl: completion in jq
* ccf6cab repl: use map in _query_slurp_wrap
* 74b5750 shquote: Remove unnecessary sb.Reset()
* 44251ca snippet: Add mp4_matrix_structure_rotation
* 79a1aea snippets: Add urldecode
* 684193a snippets: add changes/1
* fda1dda snippets: urlencode: only 0-9a-f
* 5ad048d tar: Fix 0 trim regression
* 49d2e61 tar: Unbreak num parsing and add test
* c8fad57 tiff: Fix reading of mluc tags with multiple records
* b55f24a todo: Add ignore range check idea
* bc1b3bf todo: Add note about symbols and iprint improvements
* 1d83554 todo: Add note about test and capture with buffer
* f839317 todo: Add some known issues
* 4d94c9a todo: Clenaup a  bit
* 17a708f todo: Remove fixed repl item
* 0af4c2b todo: Update about readline
* d03a1c9 todo: add echo '{} {} {}' | jq difference
* f9622c2 vorbis_comment: Fix field name typo
* be0fdbe vp9: Add profile and fix reserved_zero field collision
* 646f902 vpx_ccr: Add color  names
* b0ad3f2 w
* 45afbe6 wip
* d838d2f zip: Add format decoder
* 9029143 zip: Fix nested decode for none compress

