PostgreSQL

### Decode content of pg_control file

```sh
$ fq -d pg_control -o flavour=postgres14 d pg_control
```

### Specific fields can be got by request

```sh
$ fq -d pg_control -o flavour=postgres14 ".state, .check_point_copy.redo, .wal_level" pg_control
```

### To see heap page's content
```sh
$ fq -d pg_heap -o flavour=postgres14 ".[0]" 16994
```

### To see page's header

```sh
$ fq -d pg_heap -o flavour=postgres14 ".[0].page_header" 16994
```

### First and last item pointers on first page

```sh
$ fq -d pg_heap -o flavour=postgres14 ".[0].pd_linp[0, -1]" 16994
```

### First and last tuple on first page

```sh
$ fq -d pg_heap -o flavour=postgres14 ".[0].tuples[0, -1]" 16994
```

### Btree index meta page

```sh
$ fq -d pg_btree -o flavour=postgres14 ".[0] | d" 16404
```

### Btree index page

```sh
$ fq -d pg_btree -o flavour=postgres14 ".[1]" 16404
```

### References
- https://github.com/postgres/postgres/blob/REL_14_2/src/include/catalog/pg_control.h
- https://www.postgresql.org/docs/current/storage-page-layout.html
