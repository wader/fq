package cli

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"fq/pkg/bitio"
	"fq/pkg/decode"
	"fq/pkg/format"
	"io"
	"io/ioutil"
	"sort"
	"strings"

	"github.com/itchyny/gojq"
)

type Main struct {
	OS       OS
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
	probeFormats := m.Registry.MustGroup(format.PROBE)

	fs := flag.NewFlagSet("", flag.ContinueOnError)
	fs.SetOutput(m.OS.Stderr())
	dotFlag := fs.Bool("dot", false, "Output dot format graph (... | dot -Tsvg -o formats.svg)")
	formatNameFlag := fs.String("f", "probe", "Format name")
	maxDisplayBytes := fs.Int64("d", 16, "Max display bytes")
	// verboseFlag := fs.Bool("v", false, "Verbose output")
	fs.Usage = func() {
		maxNameLen := 0
		for _, f := range allFormats {
			if len(f.Name) > maxNameLen {
				maxNameLen = len(f.Name)
			}
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
		fmt.Fprintf(fs.Output(), "Name:%s    MIME:\n", pad(maxNameLen, "Name:"))
		for _, f := range formatsSorted {
			fmt.Fprintf(fs.Output(), "%s%s    %s\n", f.Name, pad(maxNameLen, f.Name), strings.Join(f.MIMEs, ", "))
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

	var rs io.ReadSeeker
	if filename != "" && filename != "-" {
		f, err := m.OS.Open(filename)
		if err != nil {
			return err
		}
		defer f.Close()
		rs = f
	} else {
		filename = "stdin"
		buf, err := ioutil.ReadAll(m.OS.Stdin())
		if err != nil {
			return err
		}
		rs = bytes.NewReader(buf)
	}

	if *formatNameFlag != "" {
		var err error
		probeFormats, err = m.Registry.Group(*formatNameFlag)
		if err != nil {
			return fmt.Errorf("%s: %s", *formatNameFlag, err)
		}
	}
	bb, err := bitio.NewBufferFromReadSeeker(rs)
	if err != nil {
		return err
	}

	dumpDefaultOpts := decode.DumpOptions{
		LineBytes:       16,
		MaxDisplayBytes: *maxDisplayBytes,
		AddrBase:        16,
		SizeBase:        10,
	}

	fqFuncs := map[string]gojq.Function{
		"bits": {
			Argcount: 1,
			Callback: func(c interface{}, a []interface{}) interface{} {
				if v, ok := c.(*decode.Value); ok {
					bb, err := v.RootBitBuf.BitBufRange(v.Range.Start, v.Range.Len)
					if err != nil {
						return err
					}
					return bb
				}

				// TODO: passthru c? move raw function?
				return nil
			},
		},
		"string": {
			Argcount: 1,
			Callback: func(c interface{}, a []interface{}) interface{} {
				var bb *bitio.Buffer
				switch cc := c.(type) {
				case *decode.Value:
					bb, err = cc.RootBitBuf.BitBufRange(cc.Range.Start, cc.Range.Len)
					if err != nil {
						return err
					}
				case *bitio.Buffer:
					bb = cc
				default:
					return fmt.Errorf("value is not a decode value or bit buffer")
				}

				sb := &strings.Builder{}
				if _, err := io.Copy(sb, bb); err != nil {
					return err
				}

				return string(sb.String())
			},
		},
		"probe": {
			Argcount: 1<<2 | 1<<1 | 1<<0,
			Callback: func(c interface{}, a []interface{}) interface{} {
				var bb *bitio.Buffer
				switch cc := c.(type) {
				case *decode.Value:
					bb, err = cc.RootBitBuf.BitBufRange(cc.Range.Start, cc.Range.Len)
					if err != nil {
						return err
					}
				case *bitio.Buffer:
					bb = cc
				default:
					return fmt.Errorf("value is not a decode value or bit buffer")
				}

				formats := probeFormats
				if len(a) == 1 {
					groupName, ok := a[0].(string)
					if !ok {
						return fmt.Errorf("format name is not a string")
					}

					formats, err = m.Registry.Group(groupName)
					if err != nil {
						return fmt.Errorf("%s: %s", groupName, err)
					}
				}

				// TODO: hmm
				name := "unname"
				if len(a) == 2 {
					var ok bool
					name, ok = a[1].(string)
					if !ok {
						return fmt.Errorf("name is not a string")
					}
				}

				dv, _, errs := decode.Probe(name, bb, formats)
				if dv == nil {
					return errs
				}

				return dv
			},
		},
		"hexdump": {
			Argcount: 1 << 0,
			Callback: func(c interface{}, a []interface{}) interface{} {
				var bb *bitio.Buffer
				switch cc := c.(type) {
				case *decode.Value:
					bb, err = cc.RootBitBuf.BitBufRange(cc.Range.Start, cc.Range.Len)
					if err != nil {
						return err
					}
				case *bitio.Buffer:
					bb = cc
				default:
					return fmt.Errorf("value is not decode value or a bit buffer")
				}

				hw := hex.Dumper(m.OS.Stdout())
				defer hw.Close()
				if _, err := io.Copy(hw, bb); err != nil {
					return err
				}

				return c
			},
		},
		"dump": {
			Argcount: 1<<1 | 1<<0,
			Callback: func(c interface{}, a []interface{}) interface{} {
				var v *decode.Value
				switch cc := c.(type) {
				case *decode.Value:
					v = cc
				case *decode.D:
					v = cc.Value
				default:
					return fmt.Errorf("%v: value is not a decode value", c)
				}

				maxDepth := 0
				if len(a) == 1 {
					var ok bool
					maxDepth, ok = a[0].(int)
					if !ok {
						return fmt.Errorf("max depth is not a int")
					}
					if maxDepth < 0 {
						return fmt.Errorf("max depth can't be negative")
					}
				}

				opts := dumpDefaultOpts
				opts.MaxDepth = maxDepth

				if err := v.Dump(m.OS.Stdout(), opts); err != nil {
					return err
				}

				return c
			},
		},
	}

	argQ := fs.Arg(1)
	if fs.Arg(1) == "" {
		argQ = "."
	}
	q := fmt.Sprintf(`probe($FQ_FORMAT; $FQ_FILENAME) | %s`, argQ)

	query, err := gojq.Parse(q)
	if err != nil {
		panic(err)
	}
	code, err := gojq.Compile(
		query,
		gojq.WithVariables([]string{
			"$FQ_FORMAT",
			"$FQ_FILENAME",
		}),
		gojq.WithExtraFunctions(fqFuncs),
	)
	if err != nil {
		panic(err)
	}

	iter := code.Run(bb, *formatNameFlag, filename)

	for {
		v, ok := iter.Next()
		if !ok {
			break
		}
		if err, ok := v.(error); ok {
			fmt.Fprintf(m.OS.Stderr(), "%s\n", err)
			break
		}

		switch vv := v.(type) {
		case *decode.Value:
			if err := vv.Dump(m.OS.Stdout(), dumpDefaultOpts); err != nil {
				return err
			}
		case *decode.D:
			if err := vv.Value.Dump(m.OS.Stdout(), dumpDefaultOpts); err != nil {
				return err
			}
		case *bitio.Buffer:
			io.Copy(m.OS.Stdout(), vv)
		case string:
			fmt.Fprintln(m.OS.Stdout(), vv)
		default:
			json.NewEncoder(m.OS.Stdout()).Encode(v)
		}
	}

	return nil
}
