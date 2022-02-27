Decodes vanilla and FAT Macho binaries

## Examples

To decode the MacOS build of `fq`

```
fq . -d macho fq
```

---
**NOTE**
`-d macho` is usually not needed.

---

Can be used to decode nested parts:

```
./fq '.load_commands[] | select(.cmd=="segment_64")' -d macho fq
```

## References:
- https://github.com/aidansteele/osx-abi-macho-file-format-reference
