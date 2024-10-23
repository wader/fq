package json

import (
	"bytes"
	"embed"
	stdjson "encoding/json"
	"errors"
	"io"

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
	interp.RegisterFormat(
		format.JSON,
		&decode.Format{
			Description: "JavaScript Object Notation",
			ProbeOrder:  format.ProbeOrderTextJSON,
			Groups:      []*decode.Group{format.Probe},
			DecodeFn:    decodeJSON,
			Functions:   []string{"_todisplay"},
		})
	interp.RegisterFS(jsonFS)
	interp.RegisterFunc1("_to_json", toJSON)
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
				d.Fatalf("%s", err.Error())
			}
			break
		}

		vs = append(vs, v)
	}

	if !lines && (len(vs) != 1 || !foundEOF) {
		d.Fatalf("trialing data after top-level value")
	}

	var s scalar.Any
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

func decodeJSON(d *decode.D) any {
	return decodeJSONEx(d, false)
}

type ToJSONOpts struct {
	Indent int
}

// TODO: share with interp code
func makeEncoder(opts ToJSONOpts) *colorjson.Encoder {
	return colorjson.NewEncoder(colorjson.Options{
		Color:  false,
		Tab:    false,
		Indent: opts.Indent,
		ValueFn: func(v any) (any, error) {
			switch v := v.(type) {
			case gojq.JQValue:
				return v.JQValueToGoJQ(), nil
			default:
				return v, nil
			}
		},
		Colors: colorjson.Colors{},
	})
}

func toJSON(_ *interp.Interp, c any, opts ToJSONOpts) any {
	cj := makeEncoder(opts)
	bb := &bytes.Buffer{}
	if err := cj.Marshal(c, bb); err != nil {
		return err
	}
	return bb.String()
}
