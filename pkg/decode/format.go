package decode

type Dep struct {
	Names   []string
	Formats *[]*Format
}

type Format struct {
	Name      string
	Groups    []string
	MIMEs     []string
	DecodeFn  func(d *D) interface{}
	SkipProbe bool
	Deps      []Dep
}

func FormatFn(d func(d *D) interface{}) []*Format {
	return []*Format{{
		DecodeFn: d,
	}}
}
