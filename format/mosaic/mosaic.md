
    # Decode file as mosaic
    $ fq -d mosaic . file

    # List indexes in a mosaic
    $ fq '.mosaic.index[].idx_description' file.mosaic

    # List only small objects
    $ fq '.mosaic.tile[].object[] | select(.size < 1000)' file.mosaic

    # Check if the reported size is correct
    $ fq 'add(.mosaic.tile[].object[].size) == .mosaic.container_meta_data.objects_total_size.value' basic.mosaic

    # Command above only works if the file does not use compression - check it with:
    $ fq '.mosaic.container_meta_data.compression_method | d' basic.mosaic

### Authors

- [@martinkirch](https://github.com/martinkirch/)


### References

- https://docs.softwareheritage.org/devel/swh-mosaic/
- https://gitlab.softwareheritage.org/swh/devel/swh-mosaic
