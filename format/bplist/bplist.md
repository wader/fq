### Show full decoding
```sh
$ fq -d bplist dv Info.plist

     |00 01 02 03 04 05 06 07 08 09 0a 0b 0c 0d 0e 0f|0123456789abcdef|.{}: testdata/Info.plist (bplist) 0x0-0x983.7 (2436)
     |                                               |                |  header{}: 0x0-0x7.7 (8)
0x000|62 70 6c 69 73 74                              |bplist          |    magic: "bplist" (valid) 0x0-0x5.7 (6)
0x000|                  30 30                        |      00        |    version: "00" (valid) 0x6-0x7.7 (2)
     |                                               |                |  objects{}: 0x8-0x87f.7 (2168)
0x000|                        df                     |        .       |    type: "dict" (13) (Dictionary) 0x8-0x8.3 (0.4)
0x000|                        df                     |        .       |    size_bits: 15 0x8.4-0x8.7 (0.4)
0x000|                           10                  |         .      |    large_size_marker: 1 (valid) 0x9-0x9.3 (0.4)
0x000|                           10                  |         .      |    exponent: 0 0x9.4-0x9.7 (0.4)
0x000|                              16               |          .     |    size_bigint: 22 0xa-0xa.7 (1)
     |                                               |                |    size: 22 0xb-NA (0)
     |                                               |                |    entries[0:22]: 0xb-0x87f.7 (2165)
     |                                               |                |      [0]{}: entry 0xb-0x207.7 (509)
0x000|                                 01            |           .    |        key_index: 1 0xb-0xb.7 (1)
0x020|   17                                          | .              |        value_index: 23 0x21-0x21.7 (1)
     |                                               |                |        key{}: 0x37-0x43.7 (13)
0x030|                     5c                        |       \        |          type: "ascii_string" (5) (ASCII encoded string) 0x37-0x37.3 (0.4)
0x030|                     5c                        |       \        |          size_bits: 12 0x37.4-0x37.7 (0.4)
     |                                               |                |          size: 12 0x38-NA (0)
0x030|                        43 46 42 75 6e 64 6c 65|        CFBundle|          value: "CFBundleName" 0x38-0x43.7 (12)

...
<snip>
...
```

### Get JSON representation
```
$ fq '. | torepr' com.apple.UIAutomation.plist
{
  "UIAutomationEnabled": true
}

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
