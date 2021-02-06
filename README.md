## FQ

jq for binary files

```
# duration of mp3 file
$ fq file.mp3 '[.frames[] | .samples_per_frame / .sample_rate] | add'
7504.169795907116

# width/height of embedded id3v2 jpeg picture
$ fq file.mp3 '.header.frames[] | select(.id == "APIC").picture.segments[] | select(.code._symbol == "SOF0")'
   |                                               |                |.header.frames[8].picture.segments[4]:
3b0|ff                                             |.               |  prefix: ff
3b0|   c0                                          | .              |  code: SOF0 (192)
3b0|      00 11                                    |  ..            |  Lf: 17
3b0|            08                                 |    .           |  P: 8
3b0|               01 40                           |     .@         |  Y: 320
3b0|                     01 40                     |       .@       |  X: 320
3b0|                           03                  |         .      |  Nf: 3
3b0|                              01 22 00 02 11 01|          ."....|  frame_components[3]:
3c0|03 11 01                                       |...             |
$  fq file.mp3 '.header.frames[] | select(.id == "APIC").picture.segments[] | select(.code._symbol == "SOF0") | {X,Y}'
{
  "X": 320,
  "Y": 320
}

# bitrate of first two and last frames
$ fq file.mp3 '[(.frames[0:2] + .frames[-3:-1])[].bitrate]'
[
  128000,
  128000,
  128000,
  128000
]

# mp4 sidx
$ fq test.mp4 '.. | select(.type == "sidx") | [.index_table[] | {size, duration}]'
[
  {
    "duration": 94208,
    "size": 104035
  },
  {
    "duration": 33792,
    "size": 38475
  }
]
```

## TODOs and ideas


### TODOs

- Nested BitBufs, how to show? what about ranges? for example compressed data, demuxed ogg
- Clean up panics, errors, better partial decode
- bitio.MultiBitReader to save memory
- Cleanup decoder API
- Save encoding for values, LE, BE, varint etc
- Cleanup decoders
- Document decode maturity/completeness
- Embed jq code using go 1.16 embed

### Ideas

- Some kind of UI, web and cli? maybe would be nice to show hex dump etc with overlapping fields?
- Would it be possible to save memory by just record range/decoder at first decode and
then decode as needed later?
- Move more things to jq code, dumper, CLI, help, REPL?

## Thanks

Would not be possible without [itchyny](https://github.com/itchyny)'s
jq implementation [gojq](https://github.com/itchyny/gojq).