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

### Authors
- Pavel Safonov
p.n.safonov@gmail.com
[@pnsafonov](https://github.com/pnsafonov)

### References
- https://www.postgresql.org/docs/current/storage-page-layout.html