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
	var v interface{}
	if err := jd.Decode(&v); err != nil {
		d.Invalid(err.Error())
	}

	d.Value.V = decode.JSON{V: v}

	return nil
}
