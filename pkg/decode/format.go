package decode

type Dependency struct {
	Names   []string
	Formats *[]*Format
}

type Format struct {
	Name         string
	Description  string
	Groups       []string
	MIMEs        []string
	DecodeFn     func(d *D, in interface{}) interface{}
	Dependencies []Dependency
}

func FormatFn(d func(d *D, in interface{}) interface{}) []*Format {
	return []*Format{{
		DecodeFn: d,
	}}
}
