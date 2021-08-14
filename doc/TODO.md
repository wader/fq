### TODO and ideas

#### CLI and REPL

- ctxstack index cancel wrong order skip?
- Pager for long output. Configurable? `$PAGER`? only explicit with some kind of syntax? `.. | less` but how?
- Nicer context cancel message
- `dump` cancel output of large root value, ex: `.frames`. Problem is dump is done by parent repl. 
- Error position "^" pointer?
- Configurable history file/name?
- Reset color at prompt? context cancel
- Auto complete $variables
- Auto complete keys that need escaping, now just filtered out
- Auto complete add "." just one and is object

#### Language

- Variables somehow? global `$VAR`?
- Cleanup/Make binary buffers make sense.
- gojq uses golang `int` for slice indexes, might be issue for non-64bit cpus

#### Functions

- buffer truncate, left/right pad?
- `toimage`? can be done in CLI with "\x1b]1337" but maybe something for a UI?
- `toplot`?
- `dump` should handle binary, make column code more generic? share with `hexdump`? (bindump also?)
- `dump` colorize/notify row range discontinuity
- `dump` truncate long arrays in output unless verbose? `dd`?
- `hexdump` etc should handle binary non byte aligned data
- `tojvalue` handle binary somehow, base64 string, truncate? md5 digest etc? configurable? `{binarystring: "md5"}`
- Function to search in binary data, regexp?
- `grep` like function somehow?
- Cleanup rework cipher functions, `ctr(aes("key"), "iv")` or `cipher(ctr("iv"), aes("key))`?
- `open` when to close file?
- Safe mode interpreter?
- Allow/deny `open` in autocomplete
- `open` leak, file and ctxreadseeker

### Tests

- WRITE_ACTUAL does not preserve comment order
- empty file test
- CLI tests, raw write, colors?
- Interactive tests

#### Documentation

- Nicer README, repl example (generate somehow)
- `help("topic")`?
- Generate from source
- `-n`, `inputs/0` and `input/0` behavior. Same as jq.
- Mention `empty.something`?
- `repl` behavior
  - `1,2,3 | repl` will start 3 repl:s
  - Behave same as CLi expres: iterate over all inputs `[1,2,3] | repl({}; .[])`
- Known issus and confusing things
  - Symbolic number has to use `._symbol` for now. For example matroska ids are number ID:s that have symbolic string names.
- Use https://github.com/fadado/JBOL/blob/master/doc/JQ-Distilled.md notation
- Decoder write guide
  - Endian inherited per buffer, reset on new buffer
  - Invalid on zero length input, assert one valid frame etc
  - Try validate input to make it not ambiguous with other decoders
  - Try to not seek and read at end while validating or early, will break progress indicator if not
  - Split bit flags etc into a field with subfields for each bit
  - Try keep code as declarative as possible
  - Split into multiple sub formats if possible
  - See the decoded tree as user interface but still has to represent the actual bit structure

#### Decode

- Store original filename somewhere? root description?
- Nicer DSL
  - More optional things? optional args or return value to modify?
- nicer "synthetic" values? now zero length
- Value should have raw, tranlated, symbolic and description? struct of map functions?
  - If symbolic should resolve to string? would make `.field == "abc"` work instead of now `.field._symbolic`
- Array as root value, adts, avc_au etc
- Decode framed/limited? framed adds unknown fields?
- Add unknown for arrays?
- Cleanup and rethink nested buffers (zip, muxed like ogg)
  - Root of nested buffer, what range?
- `dump` has some bug not showing buffer "nesting" level in the address column
- Endian bitfield helper (elf etc)
- Cleanup checksums, should just be fields and add warning if mismatch?
- Decoder in jq
  - Use jq array/object syntax and pass around decode context, collect fields and build tree
- Somehow control/limit nested decoding, depth/exclude/include? `probe({depth:1})` etc? per format skip options?
- `json` decoder?
- Can't use range while decoding, not calculated yet
- Keep track of encoding for values, u16le, utf8, varint etc

#### Formats

- Pass argument to format
- Move format helpers like `mp4_path` to mp4 format code?
- Value decoder in jq `u(32)`, `u32`?
- Warnings and errors
  - `mp4` sample counts
  - `flac` truncated picture, mix sample rate, bitdepth etc?
- `protobuf` schema?
- `matroska` crc
- `mp4` styp segment test
- Document maturity/completeness
- Refactor *[]decode.Format into something more abstract, group?
- Add `dsf` format
- Make `json` format more normal? is a bit a of a special case now

#### Scripts

- Summary script, tree with format specific summaries like codec, sample count etc etc?
- Probe tool with common field names
- MIME codec encode/decoder "avc1.PPCCLL" etc https://tools.ietf.org/html/rfc6381#section-3.3
- Validate scripts for mp4, matroska

#### gojq

- `JQValue` tests
- Common errors with gojq? re-implemented now
- `JQValue` update/assign function in interface, proper or just one to return error for now
- `join` can be exponential, try add with strings faster, use add `["a","b","c"] | add`?
- `0b` -> `1.7976931348623157e+308` something fishy with bin/hex/... literals change
- Do something similar to `builtin.go` in gojq to speedup a bit
- remove `scopedump`?

#### Readline

- Use something else than `github.com/chzyer/readline`?
- Fix backspace as start clears whole line bug
- Undo ctrl+_ would be nice

#### Big things

- UI, web interface? multiple repl windows? nicer way of showing overlapping fiends in hex etc?
- jupyter notebook integration
- FUSE interface
- Lazy decode, should work on known sizes? could also save memory be re-decode?
