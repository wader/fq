package decode

type Dep struct {
	Names   []string
	Formats *[]*Format
}

type Format struct {
	Name      string
	Groups    []string
	MIMEs     []string
	New       func() Decoder
	SkipProbe bool
	Deps      []Dep
}

func FormatFn(d func(c *Common)) []*Format {
	return []*Format{{
		New: func() Decoder { return &DecoderFn{decode: d} },
	}}
}

type DecoderFn struct {
	Common
	decode func(c *Common)
}

func (d *DecoderFn) Decode() {
	d.decode(&d.Common)
}
