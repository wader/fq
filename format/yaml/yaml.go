package yaml

// TODO: yaml type eval? walk eval?

import (
	"embed"
	"errors"
	"io"

	"github.com/wader/fq/format"
	"github.com/wader/fq/internal/gojqex"
	"github.com/wader/fq/pkg/bitio"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
	"github.com/wader/fq/pkg/scalar"
	"gopkg.in/yaml.v3"
)

//go:embed yaml.jq
var yamlFS embed.FS

func init() {
	interp.RegisterFormat(decode.Format{
		Name:        format.YAML,
		Description: "YAML Ain't Markup Language",
		ProbeOrder:  format.ProbeOrderTextFuzzy,
		Groups:      []string{format.PROBE},
		DecodeFn:    decodeYAML,
		Functions:   []string{"_todisplay"},
	})
	interp.RegisterFS(yamlFS)
	interp.RegisterFunc0("toyaml", toYAML)
}

func decodeYAML(d *decode.D, _ any) any {
	br := d.RawLen(d.Len())
	var r any

	yd := yaml.NewDecoder(bitio.NewIOReader(br))
	if err := yd.Decode(&r); err != nil {
		d.Fatalf("%s", err)
	}
	if err := yd.Decode(new(any)); !errors.Is(err, io.EOF) {
		d.Fatalf("trialing data after top-level value")
	}

	var s scalar.S
	s.Actual = r

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

func toYAML(_ *interp.Interp, c any) any {
	b, err := yaml.Marshal(gojqex.Normalize(c))
	if err != nil {
		return err
	}
	return string(b)
}
