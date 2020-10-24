package decode

type Dep struct {
	Names   []string
	Formats *[]*Format
}

type Format struct {
	Name      string
	Groups    []string
	MIMEs     []string
	DecodeFn  func(d *Common) interface{}
	SkipProbe bool
	Deps      []Dep
}

func FormatFn(d func(c *Common) interface{}) []*Format {
	return []*Format{{
		DecodeFn: d,
	}}
}

type DecoderFn struct {
	Common
	decode func(c *Common)
}

func (d *DecoderFn) Decode() {
	d.decode(&d.Common)
}
