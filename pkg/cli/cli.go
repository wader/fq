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

	gojqOptions := []gojq.CompilerOption{
		gojq.WithFunction("bits", 0, 0, func(c interface{}, a []interface{}) interface{} {
			if v, ok := c.(*decode.Value); ok {
				bb, err := v.RootBitBuf.BitBufRange(v.Range.Start, v.Range.Len)
				if err != nil {
					return err
				}
				return bb
			}

			// TODO: passthru c? move raw function?
			return nil
		}),
		gojq.WithFunction("string", 0, 0, func(c interface{}, a []interface{}) interface{} {
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
		}),
		gojq.WithFunction("probe", 0, 2, func(c interface{}, a []interface{}) interface{} {
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

			opts := map[string]interface{}{}
			formats := probeFormats

			if len(a) >= 1 {
				formatName, ok := a[0].(string)
				if !ok {
					return fmt.Errorf("format name is not a string")
				}

				if strings.HasSuffix(formatName, ".jq") {
					formats, err = m.Registry.Group("jq")

					script, err := ioutil.ReadFile(formatName)
					if err != nil {
						return err
					}
					opts["script"] = string(script)
				} else {
					formats, err = m.Registry.Group(formatName)
					if err != nil {
						return fmt.Errorf("%s: %s", formatName, err)
					}
				}
			}

			// TODO: hmm
			name := "unname"
			if len(a) >= 2 {
				var ok bool
				name, ok = a[1].(string)
				if !ok {
					return fmt.Errorf("name is not a string")
				}
			}

			dv, _, errs := decode.Probe(name, bb, formats, decode.ProbeOptions{FormatOptions: opts})
			if dv == nil {
				return errs
			}

			return dv
		}),

		gojq.WithFunction("hexdump", 0, 0, func(c interface{}, a []interface{}) interface{} {
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
		}),
		gojq.WithFunction("dump", 0, 1, func(c interface{}, a []interface{}) interface{} {
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
		}),
		gojq.WithVariables([]string{
			"$FQ_FORMAT",
			"$FQ_FILENAME",
		}),
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
	code, err := gojq.Compile(query, gojqOptions...)
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
