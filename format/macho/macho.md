Supports decoding vanilla and FAT Mach-O binaries.

#### Examples

To decode the macOS build of `fq`:

```
fq . /path/to/fq
```

```
fq '.load_commands[] | select(.cmd=="segment_64")' /path/to/fq
```

Note you can use `-d macho` to decode a broken Mach-O binary.

#### References:
- https://github.com/aidansteele/osx-abi-macho-file-format-reference
