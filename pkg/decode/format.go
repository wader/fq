package decode

type Group struct {
	Name         string
	Formats      []*Format
	DefaultInArg any
}

type Dependency struct {
	Groups []*Group
	Out    *Group
}

type Format struct {
	Name               string
	ProbeOrder         int // probe order is from low to hi value then by name
	Description        string
	Groups             []*Group
	DecodeFn           func(d *D) any
	DefaultInArg       any
	RootArray          bool
	RootName           string
	Dependencies       []Dependency
	Functions          []string
	SkipDecodeFunction bool
}

func FormatFn(fn func(d *D) any) *Group {
	return &Group{
		Formats: []*Format{
			{DecodeFn: fn},
		},
	}
}
