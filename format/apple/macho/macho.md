Supports decoding vanilla and FAT Mach-O binaries.

### Select 64bit load segments

```sh
$ fq '.load_commands[] | select(.cmd=="segment_64")' file
```

### References
- https://github.com/aidansteele/osx-abi-macho-file-format-reference

### Authors
- Sıddık AÇIL
acils@itu.edu.tr
[@Akaame](https://github.com/Akaame)
