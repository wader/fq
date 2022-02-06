`protobuf` decoder can be used to decode sub messages:

```
fq -d protobuf '.fields[6].wire_value | protobuf | d'
```
