### Decode content of pg_control file

```sh
$ fq -d pg_control -o flavour=postgres14 d pg_control
```

### Specific fields can be got by request

```sh
$ fq -d pg_control -o flavour=postgres14 ".state, .check_point_copy.redo, .wal_level" pg_control
```

### Authors
- Pavel Safonov
p.n.safonov@gmail.com
[@pnsafonov](https://github.com/pnsafonov)

### References
- https://github.com/postgres/postgres/blob/REL_14_2/src/include/catalog/pg_control.h