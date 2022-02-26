# MACHO Format

Decodes vanilla and FAT Macho binaries

## Examples

```
fq . -d macho fq
```

Can be used to decode nested parts:

```
./fq '.load_commands[] | select(.cmd=="segment_64")' -d macho fq
```

## References:
- https://github.com/aidansteele/osx-abi-macho-file-format-reference
