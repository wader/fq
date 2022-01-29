Supports `torepr`, ex:

```
fq -d cbor torepr file.cbor
fq -d cbor 'torepr.field' file.cbor
fq -d cbor 'torepr | .field' file.cbor
fq -d cbor 'torepr | grep("abc")' file.cbor
```
