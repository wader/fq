package json

import (
	stdjson "encoding/json"

	"github.com/wader/fq/format"
	"github.com/wader/fq/format/registry"
	"github.com/wader/fq/pkg/decode"
)

func init() {
	registry.MustRegister(&decode.Format{
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
	var s decode.Scalar
	if err := jd.Decode(&s.Actual); err != nil {
		d.Invalid(err.Error())
	}
	switch s.Actual.(type) {
	case map[string]interface{},
		[]interface{}:
	default:
		d.Invalid("root not object or array")
	}
	// TODO: root not array/struct how to add unknown gaps?
	// TODO: ranges not end up correct
	d.Value.V = s
	d.Value.Range.Len = jd.InputOffset() * 8

	return nil
}
