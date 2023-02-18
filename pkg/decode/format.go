package decode

type Group []Format

type Dependency struct {
	Names []string
	Group *Group
}

type Format struct {
	Name               string
	ProbeOrder         int // probe order is from low to hi value then by name
	Description        string
	Groups             []string
	DecodeFn           func(d *D) any
	DefaultInArg       any
	RootArray          bool
	RootName           string
	Dependencies       []Dependency
	Help               FormatHelp
	Functions          []string
	SkipDecodeFunction bool
}

type HelpExample struct {
	Comment string
	Code    string
}

type HelpFunction struct {
	Name     string
	Examples []HelpExample
}

type HelpReference struct {
	Title string
	URL   string
}

type FormatHelp struct {
	Notes      string
	Functions  []HelpFunction
	References []HelpReference
}

func FormatFn(d func(d *D) any) Group {
	return Group{{
		DecodeFn: d,
	}}
}
