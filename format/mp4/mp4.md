Supports `mp4_path`:

```
$ fq 'mp4_path(".moov.trak[1]")' file.mp4
     │00 01 02 03 04 05 06 07 08 09│0123456789│.boxes[3].boxes[1]{}:
0x4f6│                     00 00 02│       ...│  size: 573
0x500│3d                           │=         │
0x500│   74 72 61 6b               │ trak     │  type: "trak" (Container for an individual track or stream)
0x500│               00 00 00 5c 74│     ...\t│  boxes[0:3]:
0x50a│6b 68 64 00 00 00 03 00 00 00│khd.......│
0x514│00 00 00 00 00 00 00 00 01 00│..........│
*    │until 0x739.7 (565)          │          │
```

```
$ fq 'first(grep_by(.type == "trak")) | mp4_path' file.mp4
".moov.trak"
```
