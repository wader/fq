package json

import (
	stdjson "encoding/json"

	"github.com/wader/fq/format"
	"github.com/wader/fq/format/registry"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/scalar"
)

// TODO: should read multiple json values or just one?
// TODO: root not array/struct how to add unknown gaps?
// TODO: ranges not end up correct
// TODO: use jd.InputOffset() * 8?

func init() {
	registry.MustRegister(decode.Format{
		Name:        format.JSON,
		Description: "JSON",
		ProbeOrder:  100, // last
		Groups:      []string{format.PROBE},
		DecodeFn:    decodeJSON,
	})
}

func decodeJSON(d *decode.D, in interface{}) interface{} {
	bb := d.RawLen(d.Len())
	jd := stdjson.NewDecoder(bb)
	var s scalar.S
	if err := jd.Decode(&s.Actual); err != nil {
		d.Fatalf(err.Error())
	}
	switch s.Actual.(type) {
	case map[string]interface{},
		[]interface{}:
	default:
		d.Fatalf("root not object or array")
	}

	d.Value.V = &s
	d.Value.Range.Len = d.Len()

	return nil
}
