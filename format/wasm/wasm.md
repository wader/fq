### Count opcode usage
```sh
$ fq '.sections[] | select(.id == "code_section") | [.. | .opcode? // empty] | count | map({key: .[0], value: .[1]}) | from_entries' file.wasm
```

### List exports and imports
```sh
$ fq '.sections | {import: map(select(.id == "import_section").content.im.x[].nm.b), export: map(select(.id == "export_section").content.ex.x[].nm.b)}' file.wasm
```

### Authors
- Takashi Oguma
[@bitbears-dev](https://github.com/bitbears-dev)
[@0xb17bea125](https://twitter.com/0xb17bea125)

### References
- https://webassembly.github.io/spec/core/
