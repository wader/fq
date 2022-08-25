package toml

import (
	"bytes"
	"embed"

	"github.com/BurntSushi/toml"
	"github.com/wader/fq/format"
	"github.com/wader/fq/internal/gojqex"
	"github.com/wader/fq/pkg/bitio"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
	"github.com/wader/fq/pkg/scalar"
)

//go:embed toml.jq
var tomlFS embed.FS

func init() {
	interp.RegisterFormat(decode.Format{
		Name:        format.TOML,
		Description: "Tom's Obvious, Minimal Language",
		ProbeOrder:  format.ProbeOrderTextFuzzy,
		Groups:      []string{format.PROBE},
		DecodeFn:    decodeTOML,
		Functions:   []string{"_todisplay"},
	})
	interp.RegisterFS(tomlFS)
	interp.RegisterFunc0("totoml", toTOML)
}

func decodeTOML(d *decode.D, _ any) any {
	br := d.RawLen(d.Len())
	var r any

	if _, err := toml.NewDecoder(bitio.NewIOReader(br)).Decode(&r); err != nil {
		d.Fatalf("%s", err)
	}
	var s scalar.S
	s.Actual = gojqex.Normalize(r)

	// TODO: better way to handle that an empty file is valid toml and parsed as an object
	switch v := s.Actual.(type) {
	case map[string]any:
		if len(v) == 0 {
			d.Fatalf("root object has no values")
		}
	case []any:
	default:
		d.Fatalf("root not object or array")
	}

	d.Value.V = &s
	d.Value.Range.Len = d.Len()

	return nil
}

func toTOML(_ *interp.Interp, c any) any {
	if c == nil {
		return gojqex.FuncTypeError{Name: "totoml", V: c}
	}

	b := &bytes.Buffer{}
	if err := toml.NewEncoder(b).Encode(gojqex.Normalize(c)); err != nil {
		return err
	}
	return b.String()
}
