### List indexes in a mosaic

```sh
$ fq '.mosaic.index[].idx_description' file.mosaic
```

### List only small objects

```sh
$ fq '.mosaic.tile[].object[] | select(.size < 1000)' file.mosaic
```

### Check if the reported size is correct

```sh
$ fq 'add(.mosaic.tile[].object[].size) == .mosaic.container_meta_data.objects_total_size.value' basic.mosaic
```

Command above only works if the file does not use compression:

### See if a file uses compression

```sh
$ fq '.mosaic.container_meta_data.compression_method' basic.mosaic
```

### Authors

- [@martinkirch](https://github.com/martinkirch/)


### References

- https://docs.softwareheritage.org/devel/swh-mosaic/
- https://gitlab.softwareheritage.org/swh/devel/swh-mosaic
