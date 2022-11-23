Decode to a slice and indexable binary of bits.

### Slice and decode bit range

```sh
$ echo 'some {"a":1} json' | fq -d bits '.[40:-48] | fromjson'
{
  "a": 1
}
```

## Index bits

```sh
✗ echo 'hello' | fq -d bits '.[4]'
1
$ echo 'hello' | fq -c -d bits '[.[range(8)]]'
[0,1,1,0,1,0,0,0]
```
