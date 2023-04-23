### Known bugs to fix

- `fq -n '"aabbccdd" | hex | tobytes[1:] | decode("bytes") | tobytes'` create binary `aabbcc` should be `bbccdd`. I think decode (raw in this case) is confused by root value buffer.
- Buffers/string duality is confusing, most string functions should be wrapped to understand binary.
- REPL cancel seems to sometimes exit a sub-REPl without properly cleanup options.
- Value errors, can only be accessed with `._error`.
- Framed (add unknown gaps) decode should be on struct level not format?
- `tovalue({bits_format: "base64"})` only affect root value.
- Auto complete of non-global variables is broken. `scope` is broken for variables.
- `echo '{} {} {}' | jq` vs `echo '{} {} {}' | fq` works differently. fq currently decodes one root format and might add unknown gap fields etc. Maybe should work differently for `json` format?
- `format/0` overlap with jq builtin `format/1`. What to rename it to? `decode_format`?
- repl expression returning a value that produced lots of output can't be interrupted. This is because ctrl-c currently only interrupts the eval interpreter, outputted value is printed (`display`) by parent interpreter.
- Rework cli/repl user interrupt (context cancel via ctrl-c), see comment in Interp.Main
- Optimize `Interp.Options` calls, now called per display. Cache per eval? needs to handle nested evals.
- `<array decode value>[{start: ...: end: ...}]` syntax a bit broken.
- REPL completion might have side effcts. Make interp.Function type know and wrap somehow? input, inputs, open, ...
- Rework group arguments so that `{is_probe:true}` is not needed. Look up group name and see if it has an argument somehow?

### TODO and ideas

#### CLI

- `--args` support
- Reset color at prompt? context cancel

#### CLI and REPL

- ctxstack index cancel wrong order, should just skip?
- Pager for long output. Configurable? `$PAGER`? only explicit with some kind of syntax? `.. | less` but how?
- `dump` cancel output of large root value, ex: `.frames`. Problem is dump is done by parent repl.
- Error position "^" pointer?
- Configurable history file/name?
- Auto complete $variables
- Auto complete keys that need escaping, now just filtered out
- Auto complete add "." just one and is object
- Use JQ_COLORS but extended to allow name= also?

#### Language

- Cleanup/Make binary buffers make sense.
- gojq uses golang `int` for slice indexes, might be issue for non-64bit cpus

#### Functions

- buffer truncate, left/right pad?
- `toimage`? can be done in CLI with "\x1b]1337" but maybe something for a UI?
- `toplot`?
- `dump` should handle binary, make column code more generic? share with `hexdump`? (bindump also?)
- `dump` colorize/notify row range discontinuity
- `hexdump` etc should handle binary non byte aligned data
- Cleanup rework cipher functions, `ctr(aes("key"), "iv")` or `cipher(ctr("iv"), aes("key))`?
- `open` when to close file?
- Safe mode interpreter?
- Allow/deny `open` in autocomplete
- `open` leak, file and ctxreadseeker
- Summary tree with format specific summaries for each format, sample count etc etc?
- List all unique paths in some compact form?

### Tests

- update tests does not preserve comment order for readlines
- empty file test
- CLI tests, raw write, colors?
- Interactive tests

#### Documentation

- `help("topic")`?
- Generate from source
- `-n`, `inputs/0` and `input/0` behavior. Same as jq.
- Mention `empty.something`?
- Use https://github.com/fadado/JBOL/blob/master/doc/JQ-Distilled.md notation

#### Decode

- Use interfaces to save memory, "Value V" interface so can have U, Str, etc implementations?
- Array of "decorations" sym, display format?
- Store original filename somewhere? description for now
- Nicer "synthetic" values? now zero length
- Cleanup and rethink nested buffers (zip, muxed like ogg)
- Endian bitfield helper (elf etc)
- Cleanup checksums, should just be fields and add warning if mismatch?
- Decoder in jq
  - Use jq array/object syntax and pass around decode context, collect fields and build tree
- Somehow control/limit nested decoding, depth/exclude/include? `probe({depth:1})` etc? per format skip options?
- Can't use range while decoding, not calculated yet
- Keep track of encoding for values, u16le, utf8, varint etc
- Option to ignore range checks, decode until read error instead. Ex: mp4 with truncated mdat.

#### Formats

- `asn1_ber` `asn1_der`, `asn1_cer` decoder
- `flatbuffer` decoder
- `capnproto` decoder
- Pass argument to format
- Value decoder in jq `u(32)`, `u32`?
- Warnings and errors
  - `mp4` sample counts
  - `flac` truncated picture, mix sample rate, bitdepth etc?
- `protobuf` schema?
- `matroska` crc
- `mp4` styp segment test
- Document maturity/completeness
- Add `dsf` format
- Make `json` format more normal? is a bit a of a special case now
- exif in mp4 (heif/heic):<br>
`. as $r | grep("iloc") | parent.items[] | select(.id == (first($r | grep("exif";"i")) | parent.id)).extends[0] as $e | $r | tobytes[$e.offset+10:$e.offset+$e.length] | exif`

#### Scripts

- Probe tool with common field names
- MIME codec encode/decoder "avc1.PPCCLL" etc https://tools.ietf.org/html/rfc6381#section-3.3
- Validate scripts for mp4, matroska

#### gojq

- Common errors with gojq? re-implemented now
- `0b` -> `1.7976931348623157e+308` something fishy with bin/hex/... literals change
- Do something similar to `builtin.go` in gojq to speedup a bit
- remove `scopedump`?

#### Big things

- fq play website?
- UI, web interface? tree interface, multiple repl windows? nicer way of showing overlapping fiends in hex etc?
- jupyter notebook integration
- FUSE interface
- Lazy decode, should work on known sizes? could also save memory be re-decode?
