### TODO and ideas

#### CLI and REPL

- Pager when long output. Configurable? `$PAGER`? only explicit with some kind of syntax? `.. | less` but how?
- Nicer context cancel message
- `dump` cancel output of large root value, ex: `.frames`. Problem is dump is done by parent repl. 
- Error position "^" pointer?

#### Functions

- buffer truncate, left/right pad?
- `toimage`? can be done i CLI with "\x1b]1337" but maybe something for a UI?
- `toplot`?
- `dump` should handle binary, make column code more generic? share with `hexdump`? (bindump also?)
- `hexdump` etc should handle binary non byte aligned data
- `tojvalue` handle binary somehow, base64 string, truncate? md5 digest etc?
- Function to search in binary data, regexp?

### Documentation

- TODO

#### Decode

- Nicer DSL
  - More optional things? optional args or return value to modify?
- Value should have raw, tranlated, symbolic and description? struct of map functions?
- Array as root value, adts, avc_au etc
- Decode framed/limited? framed adds unknown fields?
- Add unknown for arrays?
- Cleanup and rethink nested buffers (zip, muxed like ogg)
- `dump` has some bug not showing buffer "nesting" level in the address column

#### Formats

- Pass argument to format
- Move format helpers like `mp4_path` to mp4 format code?
- Decoder in jq
- Value decoder in jq `u(32)`, `u32`?
- Warnings and errors

#### Documentation/FAQ

- Mention `empty.something`?

#### Scripts

- Summary script, tree with format specific summaries like codec, sample count etc etc?

## gojq

- `JQValue` tests

### UI

- Web?
