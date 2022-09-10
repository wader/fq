#### Convert represented value to JSON

```
$ fq -d bson torepr file.bson
```

### Filter represented value

```
$ fq -d bson 'torepr | select(.name=="bob")' file.bson
```

#### References
- https://bsonspec.org/spec.html
