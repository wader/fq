### Enumerate contents
```sh
$ fq '.objects[] | select(.id == "code_section") | [.. | .opcode? // empty] | count | map({key: .[0], value: .[1]}) | from_entries' file.wasm
```

### Inspect trailer
```sh
$ fq '.sections | {import: map(select(.id == "import_section").content.im.x[].nm.b), export: map(select(.id == "export_section").content.ex.x[].nm.b)}' file.wasm
```

### Authors
- David McDonald
[@dgmcdona](https://github.com/dgmcdona)

### References
- http://fileformats.archiveteam.org/wiki/Property_List/Binary
- https://medium.com/@karaiskc/understanding-apples-binary-property-list-format-281e6da00dbd
- https://opensource.apple.com/source/CF/CF-550/CFBinaryPList.c
