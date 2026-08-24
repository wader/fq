## EBML helpers

This module does not implements a decoder, but helper functions to decode EBML-based
formats (in particular, Matroska).

### Elements generator

`format/ebml/gen` is a program that can generate FQ elements from an EBML schema written
as XML:

    go run format/ebml/gen/main.go \
        format/matroska/ebml_matroska/ebml_matroska.xml \
        ebml_matroska \
        github.com/wader/fq/format/ebml Segment \
        | gofmt -s > format/matroska/ebml_matroska/ebml_matroska_gen.go

### References

- https://www.rfc-editor.org/info/rfc8794
