package cli

import (
	"context"
	"flag"
	"fmt"
	"fq"
	"fq/pkg/decode"
	"fq/pkg/osenv"
	"fq/pkg/query"
	"io"
	"io/ioutil"
	"os"
	"sort"
	"strings"
)

type StandardOS struct{}

func (StandardOS) Stdin() io.Reader                        { return os.Stdin }
func (StandardOS) Stdout() io.Writer                       { return os.Stdout }
func (StandardOS) Stderr() io.Writer                       { return os.Stderr }
func (StandardOS) Args() []string                          { return os.Args }
func (StandardOS) Open(name string) (io.ReadSeeker, error) { return os.Open(name) }

func StandardOSMain(r *decode.Registry) {
	if err := (Main{
		OS:       StandardOS{},
		Registry: r,
	}).Run(); err != nil {
		os.Exit(1)
	}
}

type Main struct {
	OS       osenv.OS
	Registry *decode.Registry
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
	versionFlag := fs.Bool("version", false, fmt.Sprintf("Show version(%s)", fq.Version))
	formatNameFlag := fs.String("f", "probe", "Format name")
	noInputFlag := fs.Bool("n", false, "No input")
	maxDisplayBytes := fs.Int64("d", 16, "Max display bytes")
	scriptFlag := fs.String("s", "", "Script path")
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
		fmt.Fprintf(fs.Output(), "Usage: %s [FLAGS] [FILE] [EXP]\n", m.OS.Args()[0])
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
	if *versionFlag {
		fmt.Fprintln(m.OS.Stdout(), fq.Version)
		return nil
	}
	filename := fs.Arg(0)

	q := query.NewQuery(query.QueryOptions{
		Variables: map[string]interface{}{
			"FILENAME": filename,
		},
		Registry: m.Registry,
		DumpOptions: decode.DumpOptions{
			LineBytes:       16,
			MaxDisplayBytes: *maxDisplayBytes,
			AddrBase:        16,
			SizeBase:        10,
		},
		OS: m.OS,
	})

	src := ""
	if *scriptFlag != "" {
		r, err := m.OS.Open(*scriptFlag)
		if err != nil {
			return err
		}
		scriptBytes, err := ioutil.ReadAll(r)
		if err != nil {
			return err
		}
		src = string(scriptBytes)
	} else {
		var srcs []string
		if !*noInputFlag {
			srcs = append(srcs,
				`open($FILENAME)`,
				*formatNameFlag)
		}
		if e := fs.Arg(1); e != "" {
			srcs = append(srcs, e)
		}
		if *replFlag {
			srcs = append(srcs, `push`)
		}
		src = strings.Join(srcs, " | ")
	}

	if _, err := q.Run(context.Background(), src, m.OS.Stdout()); err != nil {
		return err
	}

	if *replFlag {
		if err := q.REPL(context.Background()); err != nil {
			return err
		}
	}

	return nil
}
