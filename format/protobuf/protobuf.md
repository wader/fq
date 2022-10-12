### Can decode sub messages

```sh
$ fq -d protobuf '.fields[6].wire_value | protobuf | d' file
```

### References
- https://developers.google.com/protocol-buffers/docs/encoding
