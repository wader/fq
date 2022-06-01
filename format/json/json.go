package json

import (
	"bytes"
	"embed"
	stdjson "encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"

	"github.com/wader/fq/format"
	"github.com/wader/fq/internal/colorjson"
	"github.com/wader/fq/pkg/bitio"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
	"github.com/wader/fq/pkg/scalar"
	"github.com/wader/gojq"
)

//go:embed json.jq
var jsonFS embed.FS

func init() {
	interp.RegisterFormat(decode.Format{
		Name:        format.JSON,
		Description: "JavaScript Object Notation",
		ProbeOrder:  format.ProbeOrderText,
		Groups:      []string{format.PROBE},
		DecodeFn:    decodeJSON,
		Functions:   []string{"_todisplay"},
		Files:       jsonFS,
	})
	interp.RegisterFunc1("_tojson", toJSON)
}

func decodeJSON(d *decode.D, _ any) any {
	br := d.RawLen(d.Len())

	// keep in sync with gojq fromJSON
	jd := stdjson.NewDecoder(bitio.NewIOReader(br))
	jd.UseNumber()
	var s scalar.S
	if err := jd.Decode(&s.Actual); err != nil {
		d.Fatalf(err.Error())
	}
	if err := jd.Decode(new(any)); !errors.Is(err, io.EOF) {
		d.Fatalf("trialing data after top-level value")
	}

	s.Actual = gojq.NormalizeNumbers(s.Actual)

	// switch s.Actual.(type) {
	// case map[string]any,
	// 	[]any:
	// default:
	// 	d.Fatalf("top-level not object or array")
	// }

	d.Value.V = &s
	d.Value.Range.Len = d.Len()

	return nil
}

type ToJSONOpts struct {
	Indent int
}

func toJSON(_ *interp.Interp, c any, opts ToJSONOpts) any {
	// TODO: share
	cj := colorjson.NewEncoder(
		false,
		false,
		opts.Indent,
		func(v any) any {
			switch v := v.(type) {
			case gojq.JQValue:
				return v.JQValueToGoJQ()
			case nil, bool, float64, int, string, *big.Int, map[string]any, []any:
				return v
			default:
				panic(fmt.Sprintf("toValue not a JQValue value: %#v %T", v, v))
			}
		},
		colorjson.Colors{},
	)
	bb := &bytes.Buffer{}
	if err := cj.Marshal(c, bb); err != nil {
		return err
	}
	return bb.String()
}
