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
	"os"
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
	forceFormatNameFlag := fs.String("f", "", "Force format")
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
		os.Exit(0)
	}

	var rs io.ReadSeeker
	if fs.Arg(0) != "" && fs.Arg(0) != "-" {
		f, err := m.OS.Open(fs.Arg(0))
		if err != nil {
			return err
		}
		defer f.Close()
		rs = f
	} else {
		buf, err := ioutil.ReadAll(m.OS.Stdin())
		if err != nil {
			return err
		}
		rs = bytes.NewReader(buf)
	}

	if *forceFormatNameFlag != "" {
		var err error
		probeFormats, err = m.Registry.Group(*forceFormatNameFlag)
		if err != nil {
			return fmt.Errorf("%s: %s", *forceFormatNameFlag, err)
		}
	}
	bb, err := bitio.NewBufferFromReadSeeker(rs)
	if err != nil {
		return err
	}

	// TODO: how to skip probe at all in some cases?
	q := "probe"
	if fs.Arg(1) != "" {
		q = fmt.Sprintf("probe | (%s)", fs.Arg(1))
	}

	query, err := gojq.Parse(q)
	if err != nil {
		panic(err)
	}

	code, err := gojq.Compile(query, gojq.WithExtraFunctions(map[string]gojq.Function{
		"bits": {
			Argcount: 1,
			Callback: func(c interface{}, a []interface{}) interface{} {
				if v, ok := c.(*decode.Value); ok {
					bb, err := v.BitBuf.BitBufRange(v.Range.Start, v.Range.Len)
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
					bb, err = cc.BitBuf.BitBufRange(cc.Range.Start, cc.Range.Len)
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
			Argcount: 1<<1 | 1<<0,
			Callback: func(c interface{}, a []interface{}) interface{} {
				var bb *bitio.Buffer
				switch cc := c.(type) {
				case *decode.Value:
					bb, err = cc.BitBuf.BitBufRange(cc.Range.Start, cc.Range.Len)
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

				// TODO: name? inject filename as variable for root? other just "." or arg?
				dv, _, errs := decode.Probe("bla", bb, formats)
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
					bb, err = cc.BitBuf.BitBufRange(cc.Range.Start, cc.Range.Len)
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
				default:
					return fmt.Errorf("value is not a decode value")
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

				if err := v.Dump(maxDepth, m.OS.Stdout()); err != nil {
					return err
				}

				return c
			},
		},
	}))
	if err != nil {
		panic(err)
	}

	iter := code.Run(bb)

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
			fmt.Fprintf(m.OS.Stdout(), "%s:\n", vv.Path())
			if err := vv.Dump(0, m.OS.Stdout()); err != nil {
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
