Supports `matroska_path`, ex:

```
$ fq 'matroska_path(".Segment.Tracks[0]") file.mkv
     │00 01 02 03 04 05 06 07 08 09│0123456789│.elements[1].elements[3]{}:
0x122│         16 54 ae 6b         │   .T.k   │  id: "Tracks" (0x1654ae6b) (A Top-Level Element of information with many tracks described.)
     │                             │          │  type: "master" (7)
0x122│                     4d bf   │       M. │  size: 3519
0x122│                           bf│         .│  elements[0:3]:
0x12c│84 cf 8b db a0 ae 01 00 00 00│..........│
0x136│00 00 00 78 d7 81 01 73 c5 88│...x...s..│
*    │until 0xee9.7 (3519)         │          │
```

```
$ fq 'first(grep_by(.id == "Tracks")) | matroska_path' test.mkv
".Segment.Tracks"
```
