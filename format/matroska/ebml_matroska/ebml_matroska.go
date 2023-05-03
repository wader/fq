package ebml_matroska

// https://raw.githubusercontent.com/ietf-wg-cellar/matroska-specification/master/ebml_matroska.xml
//go:generate sh -c "go run ../ebml/gen/main.go ebml_matroska.xml ebml_matroska github.com/wader/fq/format/matroska/ebml Segment | gofmt -s > ebml_matroska_gen.go"
