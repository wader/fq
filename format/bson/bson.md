### Limitations

- The decimal128 type is not supported for decoding, will just be treated as binary

### Convert represented value to JSON

```
$ fq -d bson torepr file.bson
```

### Filter represented value

```
$ fq -d bson 'torepr | select(.name=="bob")' file.bson
```

### Authors
- Mattias Wadman mattias.wadman@gmail.com, original author
- Matt Dale [@matthewdale](https://github.com/matthewdale), additional types and bug fixes

### References
- https://bsonspec.org/spec.html
