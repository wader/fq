Supports ZIP64.

## Timestamp and time zones

The timestamp accessed via `.local_files[].last_modification` is encoded in ZIP files using [MS-DOS representation](https://learn.microsoft.com/en-us/windows/win32/api/oleauto/nf-oleauto-dosdatetimetovarianttime) which lacks a known time zone. Probably the local time/date was used at creation. The `unix_guess` field in `last_modification` is a guess assuming the local time zone was UTC at creation.

### References
- https://pkware.cachefly.net/webdocs/casestudies/APPNOTE.TXT
- https://opensource.apple.com/source/zip/zip-6/unzip/unzip/proginfo/extra.fld
- https://formats.kaitai.io/dos_datetime/
- https://learn.microsoft.com/en-us/windows/win32/api/oleauto/nf-oleauto-dosdatetimetovarianttime
