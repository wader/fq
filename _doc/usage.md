### Useful tricks

#### `.. | select(...)` fails with `expected an ... but got: ...`

Try add `select(...)?` the select expression assumes it will get and object etc.

#### Manual decode

Sometimes fq fails to decode or you know there is valid data buried inside some binary or maybe
you know the format of some unknown value. Then you can decode manually.

<pre>
# try decode a `mp3_frame` that failed to decode
$ fq file.mp3 .unknown0 mp3_frame
# skip first 10 bytes then decode as `mp3_frame`
$ fq file.mp3 .unknown0._bytes[10:] mp3_frame
</pre>

#### Run pipelines using CLI arguments
<pre sh>
$ fq file.mp3 .frames[0].header.bitrate radix2 
"1101101011000000"
</pre>
instead of:
<pre sh>
$ fq file.mp3 '.frames[0].header.bitrate | radix2' 
"1101101011000000"
</pre>
this can also be used with interactive mode
```sh
$ fq -i file.flac .metadatablocks[0] 
.metadatablocks[0] flac_metadatablock> 
```

#### appending to array is slow

Try to use `map` or `foreach` instead.

#### Use `print` and `println` to produce more friendly compact output

```
> [[0,"a"],[1,"b"]]
[
  [
    0,
    "a"
  ],
  [
    1,
    "b"
  ]
]
> [[0,"a"],[1,"b"]] | .[] | "\(.[0]): \(.[1])" | println
0: a
1: b
```

#### Run interactive mode with no input
```sh
fq -i
null>
```