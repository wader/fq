package ebml_matroska

// https://raw.githubusercontent.com/cellar-wg/matroska-specification/aa2144a58b661baf54b99bab41113d66b0f5ff62/ebml_matroska.xml
//go:generate sh -c "go run ../ebml/gen/main.go ebml_matroska.xml ebml_matroska fq/format/matroska/ebml '' Segment | gofmt > ebml_matroska_gen.go"
