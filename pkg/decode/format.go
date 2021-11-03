package decode

import "io/fs"

type Dependency struct {
	Names   []string
	Formats *[]*Format // TODO: rename to outFormats to make it clear it's used to assign?
}

type Format struct {
	Name         string
	ProbeOrder   int // probe order is from low to hi value then by name
	Description  string
	Groups       []string
	DecodeFn     func(d *D, in interface{}) interface{}
	RootArray    bool
	RootName     string
	Dependencies []Dependency
	Files        fs.ReadDirFS
}

func FormatFn(d func(d *D, in interface{}) interface{}) []*Format {
	return []*Format{{
		DecodeFn: d,
	}}
}
