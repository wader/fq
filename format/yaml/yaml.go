package yaml

// TODO: yaml type eval? walk eval?

import (
	"bytes"
	"embed"
	"errors"
	"io"

	"github.com/wader/fq/format"
	"github.com/wader/fq/internal/gojqx"
	"github.com/wader/fq/pkg/bitio"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
	"github.com/wader/fq/pkg/scalar"
	"gopkg.in/yaml.v3"
)

//go:embed yaml.jq
var yamlFS embed.FS

func init() {
	interp.RegisterFormat(
		format.YAML,
		&decode.Format{
			Description: "YAML Ain't Markup Language",
			ProbeOrder:  format.ProbeOrderTextFuzzy,
			Groups:      []*decode.Group{format.Probe},
			DecodeFn:    decodeYAML,
			DefaultInArg: format.YAML_In{
				MultiDocument: false,
			},
			Functions: []string{"_todisplay"},
		})
	interp.RegisterFS(yamlFS)
	interp.RegisterFunc1("_to_yaml", toYAML)
}

func decodeYAML(d *decode.D) any {
	var yi format.YAML_In
	d.ArgAs(&yi)

	br := d.RawLen(d.Len())

	var vs []any

	yd := yaml.NewDecoder(bitio.NewIOReader(br))
	for {
		var v any
		err := yd.Decode(&v)
		if err != nil {
			if len(vs) == 0 {
				d.Fatalf("%s", err)
			} else if errors.Is(err, io.EOF) {
				break
			} else {
				d.Fatalf("trialing data after document")
			}
		}

		vs = append(vs, v)
	}

	var s scalar.Any
	if !yi.MultiDocument && len(vs) == 1 {
		s.Actual = gojqx.Normalize(vs[0])
	} else {
		s.Actual = gojqx.Normalize(vs)
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

type ToYAMLOpts struct {
	Indent        int  `default:"4"` // 4 is default for gopkg.in/yaml.v3
	MultiDocument bool `default:"false"`
}

func toYAML(_ *interp.Interp, c any, opts ToYAMLOpts) any {
	c = gojqx.Normalize(c)

	cs, isArray := c.([]any)
	if opts.MultiDocument {
		if !isArray {
			return gojqx.FuncTypeError{Name: "to_yaml", V: c}
		}
	} else {
		cs = []any{c}
	}

	b := &bytes.Buffer{}
	e := yaml.NewEncoder(b)
	for _, c := range cs {
		// yaml.SetIndent panics if < 0
		if opts.Indent >= 0 {
			e.SetIndent(opts.Indent)
		}
		if err := e.Encode(gojqx.Normalize(c)); err != nil {
			return err
		}
	}

	return b.String()
}
