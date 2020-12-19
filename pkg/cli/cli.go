package cli

// TODO: REPL and reading from stdin?

import (
	"flag"
	"fmt"
	"fq/pkg/decode"
	"fq/pkg/osenv"
	"fq/pkg/query"
	"os"
	"sort"
	"strings"
)

type Main struct {
	OS       osenv.OS
	Registry *decode.Registry
}

func StandardOSMain(r *decode.Registry) {
	if err := (Main{
		OS:       osenv.StandardOS{},
		Registry: r,
	}).Run(); err != nil {
		os.Exit(1)
	}
}

// Run cli main
func (m Main) Run() error {
	err := m.run()
	if err != nil && err != flag.ErrHelp {
		fmt.Fprintf(m.OS.Stderr(), "%s\n", err)
	}
	return err
}

func (m Main) run() error {
	allFormats := m.Registry.MustAll()

	fs := flag.NewFlagSet("", flag.ContinueOnError)
	fs.SetOutput(m.OS.Stderr())
	dotFlag := fs.Bool("dot", false, "Output dot format graph (... | dot -Tsvg -o formats.svg)")
	formatNameFlag := fs.String("f", "probe", "Format name")
	maxDisplayBytes := fs.Int64("d", 16, "Max display bytes")
	replFlag := fs.Bool("i", false, "REPL")
	// verboseFlag := fs.Bool("v", false, "Verbose output")
	fs.Usage = func() {
		maxNameLen := 0
		maxDescriptionLen := 0
		for _, f := range allFormats {
			m := func(a, b int) int {
				if a > b {
					return a
				}
				return b
			}
			maxNameLen = m(maxNameLen, len(f.Name))
			maxDescriptionLen = m(maxDescriptionLen, len(f.Description))
		}

		formatsSorted := make([]*decode.Format, len(allFormats))
		copy(formatsSorted, allFormats)
		sort.Slice(formatsSorted, func(i, j int) bool {
			return formatsSorted[i].Name < formatsSorted[j].Name
		})

		pad := func(n int, s string) string { return strings.Repeat(" ", n-len(s)) }
		fmt.Fprintf(fs.Output(), "Usage: %s [FLAGS] FILE [EXP]\n", m.OS.Args()[0])
		fs.PrintDefaults()
		fmt.Fprintf(fs.Output(), "\n")
		fmt.Fprintf(fs.Output(), "Name:%s  Description:%s  MIME:\n", pad(maxNameLen, "Name:"), pad(maxNameLen, "Description:"))
		for _, f := range formatsSorted {
			fmt.Fprintf(fs.Output(), "%s%s  %s%s  %s\n", f.Name, pad(maxNameLen, f.Name), f.Description, pad(maxDescriptionLen, f.Description), strings.Join(f.MIMEs, ", "))
		}
	}
	if err := fs.Parse(m.OS.Args()[1:]); err != nil {
		return err
	}

	if *dotFlag {
		m.Registry.Dot(m.OS.Stdout())
		return nil
	}

	filename := fs.Arg(0)

	q := query.NewQuery(query.QueryOptions{
		Filename: filename,
		Registry: m.Registry,
		DumpOptions: decode.DumpOptions{
			LineBytes:       16,
			MaxDisplayBytes: *maxDisplayBytes,
			AddrBase:        16,
			SizeBase:        10,
		},
		OS: m.OS,
	})

	src := fs.Arg(1)
	if fs.Arg(1) == "" {
		src = "."
	}
	src = fmt.Sprintf(`open($FILENAME) | %s | %s | dot`, *formatNameFlag, src)

	if _, err := q.Run(src); err != nil {
		return err
	}

	if *replFlag {
		if err := q.REPL(); err != nil {
			return err
		}
	}

	return nil
}
