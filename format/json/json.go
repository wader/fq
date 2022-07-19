package json

import (
	stdjson "encoding/json"

	"github.com/wader/fq/format"
	"github.com/wader/fq/pkg/bitio"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
	"github.com/wader/fq/pkg/scalar"
)

// TODO: should read multiple json values or just one?
// TODO: root not array/struct how to add unknown gaps?
// TODO: ranges not end up correct
// TODO: use jd.InputOffset() * 8?

func init() {
	interp.RegisterFormat(decode.Format{
		Name:        format.JSON,
		Description: "JSON",
		ProbeOrder:  100, // last
		Groups:      []string{format.PROBE},
		DecodeFn:    decodeJSON,
	})
}

func decodeJSON(d *decode.D, _ any) any {
	br := d.RawLen(d.Len())
	jd := stdjson.NewDecoder(bitio.NewIOReader(br))
	var s scalar.S
	if err := jd.Decode(&s.Actual); err != nil {
		d.Fatalf(err.Error())
	}
	switch s.Actual.(type) {
	case map[string]any,
		[]any:
	default:
		d.Fatalf("root not object or array")
	}

	d.Value.V = &s
	d.Value.Range.Len = d.Len()

	return nil
}
