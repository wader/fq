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
		ProbeOrder:  format.ProbeOrderTextJSON,
		Groups:      []string{format.PROBE},
		DecodeFn:    decodeJSON,
		Functions:   []string{"_todisplay"},
	})
	interp.RegisterFS(jsonFS)
	interp.RegisterFunc1("_tojson", toJSON)
}

func decodeJSONEx(d *decode.D, lines bool) any {
	var vs []any

	// keep in sync with gojq fromJSON
	jd := stdjson.NewDecoder(bitio.NewIOReader(d.RawLen(d.Len())))
	jd.UseNumber()

	foundEOF := false

	for {
		var v any
		if err := jd.Decode(&v); err != nil {
			if errors.Is(err, io.EOF) {
				foundEOF = true
				if lines {
					break
				} else if len(vs) == 1 {
					break
				}
			} else if lines {
				d.Fatalf(err.Error())
			}
			break
		}

		vs = append(vs, v)
	}

	if !lines && (len(vs) != 1 || !foundEOF) {
		d.Fatalf("trialing data after top-level value")
	}

	var s scalar.S
	if lines {
		if len(vs) == 0 {
			d.Fatalf("not lines found")
		}
		s.Actual = gojq.NormalizeNumbers(vs)
	} else {
		s.Actual = gojq.NormalizeNumbers(vs[0])
	}
	d.Value.V = &s
	d.Value.Range.Len = d.Len()

	return nil
}

func decodeJSON(d *decode.D, _ any) any {
	return decodeJSONEx(d, false)
}

type ToJSONOpts struct {
	Indent int
}

// TODO: share with interp code
func makeEncoder(opts ToJSONOpts) *colorjson.Encoder {
	return colorjson.NewEncoder(
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
}

func toJSON(_ *interp.Interp, c any, opts ToJSONOpts) any {
	cj := makeEncoder(opts)
	bb := &bytes.Buffer{}
	if err := cj.Marshal(c, bb); err != nil {
		return err
	}
	return bb.String()
}
