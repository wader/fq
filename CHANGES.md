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
* c0202483 hevc_vpc,hevc_sps: Use same nameing for profile as in spec
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
- Speedup interpeter by skipping redundant includes. #172
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
  - Support all types but real type is currently limited to range for 64 bit interger/float.
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
* 5a1d35e7 Remove redudant question and fix typo
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
* 8e47fb1a doc,matroska: Fix filesname in example
* c15f5283 doc: Add format links to format table
* b86da7ae doc: Add inital decoder API documentation
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
* 36d2891 readline: Update to verison with less deps
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
* 9b683cd deocde: Cleanup some io panic(err)
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
* c14c29a cli: Cleanup and more commens
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
* 1a0089e doc: Fix typo and some improvments
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
* 2daa738 flac_frame: Use d.Invalid for posssible errors
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
* 08ec4f0 funcs: Remove unsued string function
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
* 95b9c32 make: doc/formats.svg: Ignore graphviz verison to get less diff
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
* bc1b3bf todo: Add note about symbols and iprint improvments
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

