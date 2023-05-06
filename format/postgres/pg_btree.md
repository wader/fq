### Btree index meta page

```sh
$ fq -d pg_btree -o flavour=postgres14 ".[0] | d" 16404
```

### Btree index page

```sh
$ fq -d pg_btree -o flavour=postgres14 ".[1]" 16404
```

### Authors
- Pavel Safonov
p.n.safonov@gmail.com
[@pnsafonov](https://github.com/pnsafonov)

### References
- https://www.postgresql.org/docs/current/storage-page-layout.html