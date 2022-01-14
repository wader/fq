package decode

import "io/fs"

type Group []Format

type Dependency struct {
	Names []string
	Group *Group
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
	ToRepr       string
}

func FormatFn(d func(d *D, in interface{}) interface{}) Group {
	return Group{{
		DecodeFn: d,
	}}
}
