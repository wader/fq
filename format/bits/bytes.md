Decode to a slice and indexable binary of bytes.

### Slice out byte ranges

```sh
$ echo -n 'hello' | fq -d bytes '.[-3:]' > last_3_bytes
$ echo -n 'hello' | fq -d bytes '[.[-2:], .[0:2]] | tobytes' > first_last_2_bytes_swapped
```

### Slice and decode byte range

```sh
$ echo 'some {"a":1} json' | fq -d bytes '.[5:-6] | fromjson'
{
  "a": 1
}
```

## Index bytes

```sh
$ echo 'hello' | fq -d bytes '.[1]'
101
```
