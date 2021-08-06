package json

import (
	stdjson "encoding/json"
	"fq/format"
	"fq/format/registry"
	"fq/pkg/decode"
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
	bb := d.BitBufLen(d.Len())
	jd := stdjson.NewDecoder(bb)
	if err := jd.Decode(&d.Value.V); err != nil {
		d.Invalid(err.Error())
	}
	switch d.Value.V.(type) {
	case map[string]interface{},
		[]interface{}:
	default:
		d.Invalid("root not object or array")
	}
	// TODO: root not array/struct how to add unknown gaps?
	// TODO: ranges not end up correct
	d.Value.Range.Len = jd.InputOffset() * 8

	return nil
}
