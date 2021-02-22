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
	"log"
	"os"
	"strings"
)

type StandardOS struct{}

func (StandardOS) Stdin() io.Reader                        { return os.Stdin }
func (StandardOS) Stdout() io.Writer                       { return os.Stdout }
func (StandardOS) Stderr() io.Writer                       { return os.Stderr }
func (StandardOS) Environ() []string                       { return os.Environ() }
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

type OptionValueFlag map[string]string

func (o OptionValueFlag) String() string {
	return "options"
}

func (o OptionValueFlag) Set(v string) error {
	parts := strings.SplitN(v, "=", 2)
	if len(parts) < 2 {
		return fmt.Errorf("not key=value")
	}
	(map[string]string)(o)[parts[0]] = parts[1]

	return nil
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
	//allFormats := m.Registry.MustAll()

	/*

		fs := flag.NewFlagSet("", flag.ContinueOnError)
		fs.SetOutput(m.OS.Stderr())
		versionFlag := fs.Bool("version", false, fmt.Sprintf("Show version(%s)", fq.Version))
		// TODO: confusing that "decode" is default?
		decodeFormatFlag := fs.String("d", "probe", "Decode format")
		noInputFlag := fs.Bool("n", false, "No input")
		fileFlag := fs.String("f", "", "Read script from file")
		replFlag := fs.Bool("i", false, "REPL")
		// TODO: refactor our to jq helper function?
		opts := map[string]string{
			"maxdepth":     "0",
			"verbose":      "false",
			"color":        `_options_default_color`,
			"unicode":      `_options_default_unicode`,
			"raw":          `_options_default_raw`,
			"linebytes":    `_options_default_linebytes`,
			"displaybytes": `_options_default_displaybytes`,
			"addrbase":     "16",
			"sizebase":     "10",
		}
		optsFlag := OptionValueFlag(opts)
		// TODO: show options? do in jq show values?
		fs.Var(optsFlag, "o", "key=value option, eg: color=true")
		fs.Usage = func() {
			maxNameLen := 0
			maxDescriptionLen := 0
			for _, f := range allFormats {
				maxNameLen = num.MaxInt(maxNameLen, len(f.Name))
				maxDescriptionLen = num.MaxInt(maxDescriptionLen, len(f.Description))
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
	*/
	//filename := fs.Arg(0)
	filename := "asd"

	var args []interface{}
	for _, a := range m.OS.Args() {
		args = append(args, a)
	}

	q := query.NewQuery(query.QueryOptions{
		Variables: map[string]interface{}{
			"$FILENAME": filename,
			"$VERSION":  fq.Version,
			"$ARGS":     args,
		},
		Registry: m.Registry,
		//Options:  opts,
		OS: m.OS,
	})

	runMode := query.ScriptMode

	/*
		if *replFlag {
			runMode = query.REPLMode
		}

		src := ""
		if *fileFlag != "" {
			r, err := m.OS.Open(*fileFlag)
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
					*decodeFormatFlag)
			}
			if e := fs.Arg(1); e != "" {
				srcs = append(srcs, e)
			}
			if *replFlag {
				srcs = append(srcs, `push`)
			}
			src = strings.Join(srcs, " | ")
		}

		_ = src
	*/

	i, err := q.Eval(context.Background(), runMode, nil, "main($ARGS)", query.WriterOutput{Ctx: context.Background(), W: m.OS.Stdout()})
	if err != nil {
		return err
	}
	for {
		v, ok := i.Next()
		log.Printf("v: %#+v\n", v)
		if !ok {
			break
		} else if err, ok := v.(error); ok {
			fmt.Fprintln(m.OS.Stderr(), err)
			break
		} else if d, ok := v.([2]interface{}); ok {
			fmt.Fprintf(m.OS.Stdout(), "%s: %v\n", d[0], d[1])
		}
	}

	// if *replFlag {
	// 	if err := q.REPL(context.Background()); err != nil {
	// 		return err
	// 	}
	// }

	return nil
}
