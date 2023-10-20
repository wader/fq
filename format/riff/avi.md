### Samples

AVI has many redundant ways to index samples so currently `.streams[].samples` will only include samples the most "modern" way used in the file. That is in order of stream super index, movi ix index then idx1 index.

### Extract samples for stream 1

```sh
$ fq '.streams[1].samples[] | tobytes' file.avi > stream01.mp3
```

### Show stream summary
```sh
$ fq -o decode_samples=false '[.chunks[0] | grep_by(.id=="LIST" and .type=="strl") | grep_by(.id=="strh") as {$type} | grep_by(.id=="strf") as {$format_tag, $compression} | {$type,$format_tag,$compression}]' *.avi
```

### Speed up decoding by disabling sample and extended chunks decoding

If your not interested in sample details or extended chunks you can speed up decoding by using:
```sh
$ fq -o decode_samples=false -o decode_extended_chunks=false d file.avi
```

### References

- [AVI RIFF File Reference](https://learn.microsoft.com/en-us/windows/win32/directshow/avi-riff-file-reference)
- [OpenDML AVI File Format Extensions](http://www.jmcgowan.com/odmlff2.pdf)
