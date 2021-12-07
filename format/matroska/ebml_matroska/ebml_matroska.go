//nolint:revive
package ebml_matroska

// https://raw.githubusercontent.com/cellar-wg/matroska-specification/aa2144a58b661baf54b99bab41113d66b0f5ff62/ebml_matroska.xml
//go:generate sh -c "go run ../ebml/gen/main.go ebml_matroska.xml ebml_matroska github.com/wader/fq/format/matroska/ebml github.com/wader/fq/pkg/scalar Segment | gofmt -s > ebml_matroska_gen.go"
