The TAP- (and BLK-) format is nearly a direct copy of the data that is stored
in real tapes, as it is written by the ROM save routine of the ZX-Spectrum.
A TAP file is simply one data block or a group of 2 or more data blocks, one
followed after the other. The TAP file may be empty.

You will often find this format embedded inside the TZX tape format.

The default file extension is `.tap`.

### Processing JSON files

When needing to process a generated JSON file it's recommended to convert the
plain data bytes to an array by setting `bits_format=byte_array`:

```bash
fq -o bits_format=byte_array -d tap -V d /path/to/file.tap
```

### Authors

- Michael R. Cook work.mrc@pm.me, original author

### References

- https://worldofspectrum.net/zx-modules/fileformats/tapformat.html
