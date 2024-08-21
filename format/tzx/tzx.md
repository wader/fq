`TZX` is a file format designed to preserve cassette tapes compatible with the
ZX Spectrum computers, although some specialized versions of the format have
been defined for other machines such as the Amstrad CPC and C64.

The format was originally created by Tomaz Kac, who was maintainer until
`revision 1.13`, before passing it to Martijn v.d. Heide. For a brief period
the company Ramsoft became the maintainers, and created revision `v1.20`.

The default file extension is `.tzx`.

### Processing JSON files

When needing to process a generated JSON file it's recommended to convert the
plain data bytes to an array by setting `bits_format=byte_array`:

```bash
fq -o bits_format=byte_array -d tzx -V d /path/to/file.tzx
```

### Authors

- Michael R. Cook work.mrc@pm.me, original author

### References

- https://worldofspectrum.net/TZXformat.html
