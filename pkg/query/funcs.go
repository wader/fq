package query

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"fq/internal/ansi"
	"fq/internal/asciiwriter"
	"fq/internal/hexdump"
	"fq/internal/hexpairwriter"
	"fq/internal/num"
	"fq/internal/progressreadseeker"
	"fq/pkg/bitio"
	"fq/pkg/decode"
	"fq/pkg/format"
	"io"
	"io/ioutil"
	"math/big"
	"net/url"
	"strings"

	"github.com/chzyer/readline"
)

var fqModuleSrc = `
def _options_default_color: tty(1).is_terminal and env.CLICOLOR!=null;
def _options_default_unicode: tty(1).is_terminal and env.CLIUNICODE!=null;
def _options_default_linebytes: if tty(1).is_terminal then [((tty(1).size[0] div 10) div 2) * 2, 4] | max else 16 end;
def _options_default_displaybytes: _options_default_linebytes;


# convert number to array of bytes
def number_to_bytes($bits):
	def _number_to_bytes($d):
		if . > 0 then
			. % $d, (. div $d | _number_to_bytes($d))
		else
			empty
		end;
	if . == 0 then [0]
	else [_number_to_bytes(1 bsl $bits)] | reverse end;
def number_to_bytes:
	number_to_bytes(8);


def from_base($base;$table):
	split("")
	| reverse
	| map($table[.])
	| if . == null then error("invalid char \(.)") else . end
	| reduce .[] as $c
		# state: [power, ans]
		([1,0]; (.[0] * $base) as $b | [$b, .[1] + (.[0] * $c)])
	| .[1];

def to_base($base;$table):
	def stream:
		recurse(if . > 0 then . div $base else empty end) | . % $base;
	if . == 0 then
		"0"
	else
		[stream] |
		reverse  |
		.[1:] |
		if $base <= ($table | length) then
			map($table[.]) | join("")
		else
			error("base too large")
		end
	end;

def base2:
	if (. | type) == "number" then to_base(2;"01")
	else from_base(2;{"0": 0, "1": 1}) end;

def base8:
	if (. | type) == "number" then to_base(8;"01234567")
	else from_base(8;{"0": 0, "1": 1, "2": 2, "3": 3,"4": 4, "5": 5, "6": 6, "7": 7}) end;

def base16:
	if (. | type) == "number" then to_base(16;"0123456789abcdef")
	else from_base(16;{
		"0": 0, "1": 1, "2": 2, "3": 3,"4": 4, "5": 5, "6": 6, "7": 7, "8": 8, "9": 9,
		"a": 10, "b": 11, "c": 12, "d": 13, "e": 14, "f": 15
	})
	end;

def base62:
	if (. | type) == "number" then to_base(62;"0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")
	else from_base(62;{
		"0": 0, "1": 1, "2": 2, "3": 3,"4": 4, "5": 5, "6": 6, "7": 7, "8": 8, "9": 9,
		"A": 10, "B": 11, "C": 12, "D": 13, "E": 14, "F": 15, "G": 16,
		"H": 17, "I": 18, "J": 19, "K": 20, "L": 21, "M": 22, "N": 23,
		"O": 24, "P": 25, "Q": 26, "R": 27, "S": 28, "T": 29, "U": 30,
		"V": 31, "W": 32, "X": 33, "Y": 34, "Z": 35,
		"a": 36, "b": 37, "c": 38, "d": 39, "e": 40, "f": 41, "g": 42,
		"h": 43, "i": 44, "j": 45, "k": 46, "l": 47, "m": 48, "n": 49,
		"o": 50, "p": 51, "q": 52, "r": 53, "s": 54, "t": 55, "u": 56,
		"v": 57, "w": 58, "x": 59, "y": 60, "z": 61
	})
	end;

def base62sp:
	if (. | type) == "number" then to_base(62;"0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	else from_base(62;{
		"0": 0, "1": 1, "2": 2, "3": 3,"4": 4, "5": 5, "6": 6, "7": 7, "8": 8, "9": 9,
		"a": 10, "b": 11, "c": 12, "d": 13, "e": 14, "f": 15, "g": 16,
		"h": 17, "i": 18, "j": 19, "k": 20, "l": 21, "m": 22, "n": 23,
		"o": 24, "p": 25, "q": 26, "r": 27, "s": 28, "t": 29, "u": 30,
		"v": 31, "w": 32, "x": 33, "y": 34, "z": 35,
		"A": 36, "B": 37, "C": 38, "D": 39, "E": 40, "F": 41, "G": 42,
		"H": 43, "I": 44, "J": 45, "K": 46, "L": 47, "M": 48, "N": 49,
		"O": 50, "P": 51, "Q": 52, "R": 53, "S": 54, "T": 55, "U": 56,
		"V": 57, "W": 58, "X": 59, "Y": 60, "Z": 61
	})
	end;

def base62:
	if (. | type) == "number" then to_base(62;"ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")
	else from_base(62;{
		"A": 0, "B": 1, "C": 2, "D": 3, "E": 4, "F": 5, "G": 6,
		"H": 7, "I": 8, "J": 9, "K": 10, "L": 11, "M": 12, "N": 13,
		"O": 14, "P": 15, "Q": 16, "R": 17, "S": 18, "T": 19, "U": 20,
		"V": 21, "W": 22, "X": 23, "Y": 24, "Z": 25,
		"a": 26, "b": 27, "c": 28, "d": 29, "e": 30, "f": 31, "g": 32,
		"h": 33, "i": 34, "j": 35, "k": 36, "l": 37, "m": 38, "n": 39,
		"o": 40, "p": 41, "q": 42, "r": 43, "s": 44, "t": 45, "u": 46,
		"v": 47, "w": 48, "x": 49, "y": 50, "z": 51,
		"0": 52, "1": 53, "2": 54, "3": 55, "4": 56, "5": 57, "6": 58, "7": 59, "8": 60, "9": 61
	})
	end;

def base64:
	if (. | type) == "number" then to_base(64;"ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/")
	else from_base(64;{
		"A": 0, "B": 1, "C": 2, "D": 3, "E": 4, "F": 5, "G": 6,
		"H": 7, "I": 8, "J": 9, "K": 10, "L": 11, "M": 12, "N": 13,
		"O": 14, "P": 15, "Q": 16, "R": 17, "S": 18, "T": 19, "U": 20,
		"V": 21, "W": 22, "X": 23, "Y": 24, "Z": 25,
		"a": 26, "b": 27, "c": 28, "d": 29, "e": 30, "f": 31, "g": 32,
		"h": 33, "i": 34, "j": 35, "k": 36, "l": 37, "m": 38, "n": 39,
		"o": 40, "p": 41, "q": 42, "r": 43, "s": 44, "t": 45, "u": 46,
		"v": 47, "w": 48, "x": 49, "y": 50, "z": 51,
		"0": 52, "1": 53, "2": 54, "3": 55, "4": 56, "5": 57, "6": 58, "7": 59, "8": 60, "9": 61,
		"+": 62, "/": 63
	})
	end;

# from https://rosettacode.org/wiki/Non-decimal_radices/Convert#jq
# unknown author
# Convert the input integer to a string in the specified base (2 to 36 inclusive)
def _convert(base):
	def stream:
		recurse(if . > 0 then . div base else empty end) | . % base;
	if . == 0 then
		"0"
	else
		[stream] |
		reverse  |
		.[1:] |
		if base <  10 then
			map(tostring) | join("")
		elif base <= 36 then
			map(if . < 10 then 48 + . else . + 87 end) | implode
		else
			error("base too large")
		end
	end;

# input string is converted from "base" to an integer, within limits
# of the underlying arithmetic operations, and without error-checking:
def _to_i(base):
	explode
	| reverse
	| map(if . > 96  then . - 87 else . - 48 end)  # "a" ~ 97 => 10 ~ 87
	| reduce .[] as $c
		# state: [power, ans]
		([1,0]; (.[0] * base) as $b | [$b, .[1] + (.[0] * $c)])
	| .[1];

# like iprint
def i:
	{
		bin: "0b\(base2)",
		oct: "0o\(base8)",
		dec: "\(.)",
		hex: "0x\(base16)",
		str: ([.] | implode),
	};

def _formats_dot:
	"# ... | dot -Tsvg -o formats.svg",
	"digraph formats {",
	"  node [shape=\"box\",style=\"rounded,filled\"]",
	"  edge [arrowsize=\"0.7\"]",
	(.[] | "  \(.name) -> {\(.dependencies | flatten? | join(" "))}"),
	(.[] | .name as $name | .groups[]? | "  \(.) -> \($name)"),
	(keys[] | "  \(.) [color=\"paleturquoise\"]"),
	([.[].groups[]?] | unique[] | "  \(.) [color=\"palegreen\"]"),
	"}";

def field_inrange($p): ._type == "field" and ._range.start <= $p and $p < ._range.stop;

`

func buildDumpOptions(ms ...map[string]interface{}) decode.DumpOptions {
	var opts decode.DumpOptions
	for _, m := range ms {
		if m != nil {
			mapSetDumpOptions(&opts, m)
		}
	}
	opts.Decorator = decoratorFromDumpOptions(opts)

	return opts
}

func mapSetDumpOptions(d *decode.DumpOptions, m map[string]interface{}) {
	if v, ok := m["maxdepth"]; ok {
		d.MaxDepth = num.MaxInt(0, toIntZ(v))
	}
	if v, ok := m["verbose"]; ok {
		d.Verbose = toBoolZ(v)
	}
	if v, ok := m["color"]; ok {
		d.Color = toBoolZ(v)
	}
	if v, ok := m["unicode"]; ok {
		d.Unicode = toBoolZ(v)
	}
	if v, ok := m["linebytes"]; ok {
		d.LineBytes = num.MaxInt(0, toIntZ(v))
	}
	if v, ok := m["displaybytes"]; ok {
		d.DisplayBytes = num.MaxInt64(0, toInt64Z(v))
	}
	if v, ok := m["addrbase"]; ok {
		d.AddrBase = num.ClampInt(2, 36, toIntZ(v))
	}
	if v, ok := m["sizebase"]; ok {
		d.SizeBase = num.ClampInt(2, 36, toIntZ(v))
	}
}

func decoratorFromDumpOptions(opts decode.DumpOptions) decode.Decorator {
	colStr := "|"
	if opts.Unicode {
		colStr = "\xe2\x94\x82"
	}
	nameFn := func(s string) string { return s }
	valueFn := func(s string) string { return s }
	byteFn := func(b byte, s string) string { return s }
	column := colStr + "\n"
	if opts.Color {
		nameFn = func(s string) string { return ansi.FgBrightBlue + s + ansi.Reset }
		valueFn = func(s string) string { return ansi.FgBrightCyan + s + ansi.Reset }
		byteFn = func(b byte, s string) string {
			switch {
			case b == 0:
				return ansi.FgBrightBlack + s + ansi.Reset
			case b >= 32 && b <= 126, b == '\r', b == '\n', b == '\f', b == '\t', b == '\v':
				return ansi.FgWhite + s + ansi.Reset
			default:
				return ansi.FgBrightWhite + s + ansi.Reset
			}
		}
		column = ansi.FgWhite + colStr + ansi.Reset + "\n"
	}

	return decode.Decorator{
		Name:   nameFn,
		Value:  valueFn,
		Byte:   byteFn,
		Column: column,
	}
}

type Decorators struct {
	Name   func(s string) string
	Value  func(s string) string
	Byte   func(b byte, s string) string
	Column string
}

// TODO: make it nicer somehow?
func (q *Query) makeFunctions(opts QueryOptions) []Function {
	fs := []Function{
		{[]string{"tty"}, 1, 1, q.tty},
		{[]string{"options"}, 0, 1, q.options},

		{[]string{"help"}, 0, 0, q.help},
		{[]string{"open"}, 0, 1, q.open},
		{[]string{"dump", "d"}, 0, 1, q.makeDumpFn(nil)},
		{[]string{"verbose", "v"}, 0, 1, q.makeDumpFn(map[string]interface{}{"verbose": true})},
		{[]string{"hexdump", "hd", "h"}, 0, 1, q.hexdump},
		{[]string{"bits"}, 0, 2, q.bits},
		{[]string{"string"}, 0, 0, q.string_},
		{[]string{"decode"}, 0, 1, q.makeDecodeFn(opts.Registry, opts.Registry.MustGroup(format.PROBE))},
		{[]string{"u"}, 0, 1, q.u},
		{[]string{"push"}, 0, 0, q.push},
		{[]string{"pop"}, 0, 0, q.pop},
		{[]string{"_value_keys"}, 0, 0, q._valueKeys},
		{[]string{"formats"}, 0, 0, q.formats},
		{[]string{"preview", "p"}, 0, 0, q.preview},
		{[]string{"md5"}, 0, 0, q.md5},
		{[]string{"base64"}, 0, 0, q.base64},
		{[]string{"unbase64"}, 0, 0, q.unbase64},
		{[]string{"hex"}, 0, 0, q.hex},
		{[]string{"unhex"}, 0, 0, q.unhex},
		{[]string{"query_escape"}, 0, 0, q.queryEscape},
		{[]string{"query_unescape"}, 0, 0, q.queryUnescape},
		{[]string{"path_escape"}, 0, 0, q.pathEscape},
		{[]string{"path_unescape"}, 0, 0, q.pathUnescape},
		{[]string{"aes_ctr"}, 1, 2, q.aesCtr},

		{[]string{"json"}, 0, 0, q._json},
	}
	for name, f := range q.opts.Registry.Groups {
		fs = append(fs, Function{[]string{name}, 0, 0, q.makeDecodeFn(opts.Registry, f)})
	}

	return fs
}

func (q *Query) tty(c interface{}, a []interface{}) interface{} {
	fd, ok := a[0].(int)
	if !ok {
		return fmt.Errorf("%v: value is not a number", a[0])
	}
	w, h, _ := readline.GetSize(fd)
	return map[string]interface{}{
		"is_terminal": readline.IsTerminal(fd),
		"size":        []interface{}{w, h},
	}
}

func (q *Query) options(c interface{}, a []interface{}) interface{} {
	if len(a) > 0 {
		opts, ok := a[0].(map[string]interface{})
		if !ok {
			return fmt.Errorf("%v: value is not object", a[0])
		}
		q.runContext.opts = opts
	}
	return q.runContext.opts
}

func (q *Query) _json(c interface{}, a []interface{}) interface{} {
	bb, _, _, err := toBitBuf(c)
	if err != nil {
		return err
	}

	buf := &bytes.Buffer{}
	if _, err := io.Copy(buf, bb); err != nil {
		return err
	}

	var vv interface{}
	if err := json.Unmarshal(buf.Bytes(), &vv); err != nil {
		return err
	}

	return vv

}

func (q *Query) hexdump(c interface{}, a []interface{}) interface{} {
	bb, r, _, err := toBitBuf(c)
	if err != nil {
		return err
	}

	return func(stdout io.Writer) error {
		bitsByteAlign := r.Start % 8
		bb, err := bb.BitBufRange(r.Start-bitsByteAlign, r.Len+bitsByteAlign)
		if err != nil {
			return err
		}

		var opts decode.DumpOptions
		if len(a) >= 1 {
			opts = buildDumpOptions(q.runContext.opts, a[0].(map[string]interface{}))
		} else {
			opts = buildDumpOptions(q.runContext.opts)
		}

		d := opts.Decorator
		hw := hexdump.New(
			stdout,
			(r.Start-bitsByteAlign)/8,
			num.DigitsInBase(bitio.BitsByteCount(r.Stop()+bitsByteAlign), true, opts.AddrBase),
			opts.AddrBase,
			opts.LineBytes,
			func(b byte) string { return d.Byte(b, hexpairwriter.Pair(b)) },
			func(b byte) string { return d.Byte(b, asciiwriter.SafeASCII(b)) },
			d.Column,
		)
		if _, err := io.Copy(hw, bb); err != nil {
			return err
		}
		hw.Close()
		return nil
	}
}

func (q *Query) formats(c interface{}, a []interface{}) interface{} {

	allFormats := map[string]*decode.Format{}

	for _, fs := range q.opts.Registry.Groups {
		for _, f := range fs {
			if _, ok := allFormats[f.Name]; ok {
				continue
			}
			allFormats[f.Name] = f
		}
	}

	vs := map[string]interface{}{}
	for _, f := range allFormats {
		vf := map[string]interface{}{
			"name":        f.Name,
			"description": f.Description,
		}

		var dependenciesVs []interface{}
		for _, d := range f.Dependencies {
			var dNamesVs []interface{}
			for _, n := range d.Names {
				dNamesVs = append(dNamesVs, n)
			}
			dependenciesVs = append(dependenciesVs, dNamesVs)
		}
		if len(dependenciesVs) > 0 {
			vf["dependencies"] = dependenciesVs
		}
		var groupsVs []interface{}
		for _, n := range f.Groups {
			groupsVs = append(groupsVs, n)
		}
		if len(groupsVs) > 0 {
			vf["groups"] = groupsVs
		}

		vs[f.Name] = vf
	}

	return vs
}

func (q *Query) preview(c interface{}, a []interface{}) interface{} {
	v, err := toValue(c)
	if err != nil {
		return fmt.Errorf("%v: value is not a decode value", c)
	}
	return func(stdout io.Writer) error {
		if err := v.Preview(stdout); err != nil {
			return err
		}
		return nil
	}
}

func (q *Query) help(c interface{}, a []interface{}) interface{} {
	return queryErrorFn(func(stdout io.Writer) error {
		for _, f := range q.functions {
			var names []string
			for _, n := range f.Names {
				for j := f.MinArity; j <= f.MaxArity; j++ {
					names = append(names, fmt.Sprintf("%s/%d", n, j))
				}
			}
			fmt.Fprintf(stdout, "%s\n", strings.Join(names, ", "))
		}
		return nil
	})
}

func (q *Query) open(c interface{}, a []interface{}) interface{} {
	var rs io.ReadSeeker

	var filename string
	if len(a) == 1 {
		var err error
		filename, err = toString(a[0])
		if err != nil {
			return fmt.Errorf("%s: %w", filename, err)
		}
	}

	if filename == "" || filename == "-" {
		filename = "stdin"
		buf, err := ioutil.ReadAll(q.opts.OS.Stdin())
		if err != nil {
			return err
		}
		rs = bytes.NewReader(buf)
	} else {
		f, err := q.opts.OS.Open(filename)
		if err != nil {
			return err
		}

		// TODO: cleanup? bitbuf have optional close method etc?
		// if c, ok := f.(io.Closer); ok {
		// 	c.Close()
		// }

		rs = f
	}

	//TODO: how to know when decode is done?
	// TODO: refactor
	bPos, err := rs.Seek(0, io.SeekCurrent)
	if err != nil {
		return err
	}
	bEnd, err := rs.Seek(0, io.SeekEnd)
	if err != nil {
		return err
	}
	if _, err := rs.Seek(bPos, io.SeekStart); err != nil {
		return err
	}

	// TODO: make nicer
	// we don't want to print any progress things after decode is done
	var decodeDoneFn func()
	if q.runContext.mode == REPLMode {
		decodeDone := false
		decodeDoneFn = func() {
			fmt.Fprint(q.runContext.stdout, "\r")
			decodeDone = true
		}

		rs = progressreadseeker.New(rs, bEnd, func(readBytes int64, length int64) {
			if decodeDone {
				return
			}
			fmt.Fprintf(q.runContext.stdout, "\r%.1f%%", (float64(readBytes)/float64(length))*100)
		})
	}

	bb, err := bitio.NewBufferFromReadSeeker(rs)
	if err != nil {
		return err
	}

	return &bitBufFile{
		bb:           bb,
		filename:     filename,
		decodeDoneFn: decodeDoneFn,
	}
}

func (q *Query) makeDumpFn(fnOpts map[string]interface{}) func(c interface{}, a []interface{}) interface{} {
	return func(c interface{}, a []interface{}) interface{} {
		v, err := toValue(c)
		if err != nil {
			return fmt.Errorf("%v: value is not a decode value", c)
		}

		var opts decode.DumpOptions
		if len(a) >= 1 {
			opts = buildDumpOptions(q.runContext.opts, fnOpts, a[0].(map[string]interface{}))
		} else {
			opts = buildDumpOptions(q.runContext.opts, fnOpts)
		}

		return func(stdout io.Writer) error {
			if err := v.Dump(stdout, opts); err != nil {
				return err
			}
			return nil
		}
	}
}

func (q *Query) makeDecodeFn(registry *decode.Registry, decodeFormats []*decode.Format) func(c interface{}, a []interface{}) interface{} {
	return func(c interface{}, a []interface{}) interface{} {
		// TODO: progress hack
		// would be nice to move progress code into decode but it might be
		// tricky to keep track of absolute positions in the underlaying readers
		// when it uses BitBuf slices, maybe only in Pos()?
		if bbf, ok := c.(*bitBufFile); ok {
			if bbf.decodeDoneFn != nil {
				defer bbf.decodeDoneFn()
			}
		}

		bb, r, filename, err := toBitBuf(c)
		if err != nil {
			return err
		}
		bb, err = bb.BitBufRange(r.Start, r.Len)
		if err != nil {
			return err
		}

		opts := map[string]interface{}{}

		name := "unnamed"
		if filename != "" {
			name = filename
		}

		if len(a) >= 1 {
			formatName, err := toString(a[0])
			if err != nil {
				return fmt.Errorf("%s: %w", formatName, err)
			}
			decodeFormats, err = registry.Group(formatName)
			if err != nil {
				return fmt.Errorf("%s: %w", formatName, err)
			}
		}

		dv, _, errs := decode.Decode(name, bb, decodeFormats, decode.DecodeOptions{FormatOptions: opts})
		if dv == nil {
			return errs
		}

		return dv
	}
}

func (q *Query) _valueKeys(c interface{}, a []interface{}) interface{} {
	if v, ok := c.(*decode.Value); ok {
		var vs []interface{}
		for _, s := range v.SpecialPropNames() {
			vs = append(vs, s)
		}
		return vs
	}
	return nil
}

func (q *Query) bits(c interface{}, a []interface{}) interface{} {
	bb, r, _, err := toBitBuf(c)
	if err != nil {
		return err
	}
	bb, err = bb.BitBufRange(r.Start, r.Len)
	if err != nil {
		return err
	}

	startArg := int64(0)
	endArg := int64(-1)
	toAbs := func(v int64, l int64) int64 {
		if v < 0 {
			return l + v + 1
		}
		return v
	}

	if len(a) >= 1 {
		startArg, err = toInt64(a[0])
		if err != nil {
			return err
		}
	}
	if len(a) >= 2 {
		endArg, err = toInt64(a[1])
		if err != nil {
			return err
		}
	}

	startArg = toAbs(startArg, bb.Len())
	endArg = toAbs(endArg, bb.Len())

	bb, err = bb.BitBufRange(startArg, endArg-startArg)
	if err != nil {
		return err
	}

	return bb
}

func (q *Query) string_(c interface{}, a []interface{}) interface{} {
	var bb *bitio.Buffer
	switch cc := c.(type) {
	case *decode.Value:
		var err error
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
}

func (q *Query) u(c interface{}, a []interface{}) interface{} {
	bb, r, _, err := toBitBuf(c)
	if err != nil {
		return err
	}

	nBits := r.Len
	if len(a) == 1 {
		n, err := toInt64(a[0])
		if err != nil {
			return err
		}
		nBits = n
	}

	bb, err = bb.BitBufRange(r.Start, nBits)
	if err != nil {
		return err
	}

	// TODO: smart and maybe use int if bits can fit?
	bi := new(big.Int)
	for i := bb.Len() - 1; i >= 0; i-- {
		v, err := bb.Bool()
		if err != nil {
			return err
		}
		if v {
			bi.SetBit(bi, int(i), 1)
		}
	}

	return bi
}

func (q *Query) push(c interface{}, a []interface{}) interface{} {
	if _, ok := c.(error); !ok {
		q.runContext.pushVs = append(q.runContext.pushVs, c)
	}
	return func(stdout io.Writer) error {
		// nop
		return nil
	}

}

func (q *Query) pop(c interface{}, a []interface{}) interface{} {
	q.runContext.pops++
	return func(stdout io.Writer) error {
		// nop
		return nil
	}
}

func (q *Query) md5(c interface{}, a []interface{}) interface{} {
	bb, _, _, err := toBitBuf(c)
	if err != nil {
		return err
	}

	md5 := md5.New()
	if _, err := io.Copy(md5, bb); err != nil {
		return err
	}

	return md5.Sum(nil)
}

func (q *Query) base64(c interface{}, a []interface{}) interface{} {
	bb, _, _, err := toBitBuf(c)
	if err != nil {
		return err
	}

	b64Buf := &bytes.Buffer{}
	b64 := base64.NewEncoder(base64.StdEncoding, b64Buf)
	if _, err := io.Copy(b64Buf, bb); err != nil {
		return err
	}
	b64.Close()

	return b64Buf.Bytes()
}

func (q *Query) unbase64(c interface{}, a []interface{}) interface{} {
	bb, _, _, err := toBitBuf(c)
	if err != nil {
		return err
	}

	buf := &bytes.Buffer{}
	if _, err := io.Copy(buf, base64.NewDecoder(base64.StdEncoding, bb)); err != nil {
		return err
	}

	return buf.Bytes()
}

func (q *Query) hex(c interface{}, a []interface{}) interface{} {
	bb, _, _, err := toBitBuf(c)
	if err != nil {
		return err
	}

	buf := &bytes.Buffer{}
	if _, err := io.Copy(hex.NewEncoder(buf), bb); err != nil {
		return err
	}

	return buf.String()
}

func (q *Query) unhex(c interface{}, a []interface{}) interface{} {
	bb, _, _, err := toBitBuf(c)
	if err != nil {
		return err
	}

	b64Buf := &bytes.Buffer{}
	if _, err := io.Copy(b64Buf, hex.NewDecoder(bb)); err != nil {
		return err
	}

	return b64Buf.Bytes()
}

func (q *Query) queryEscape(c interface{}, a []interface{}) interface{} {
	s, err := toString(c)
	if err != nil {
		return err
	}
	return url.QueryEscape(s)
}

func (q *Query) queryUnescape(c interface{}, a []interface{}) interface{} {
	s, err := toString(c)
	if err != nil {
		return err
	}
	u, err := url.QueryUnescape(s)
	if err != nil {
		return err
	}
	return u
}
func (q *Query) pathEscape(c interface{}, a []interface{}) interface{} {
	s, err := toString(c)
	if err != nil {
		return err
	}
	return url.PathEscape(s)
}

func (q *Query) pathUnescape(c interface{}, a []interface{}) interface{} {
	s, err := toString(c)
	if err != nil {
		return err
	}
	u, err := url.PathUnescape(s)
	if err != nil {
		return err
	}
	return u
}

func (q *Query) aesCtr(c interface{}, a []interface{}) interface{} {
	keyBytes, err := toBytes(a[0])
	if err != nil {
		return err
	}

	switch len(keyBytes) {
	case 16, 24, 32:
	default:
		return fmt.Errorf("key length should be 16, 24 or 32 bytes, is %d bytes", len(keyBytes))
	}

	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return err
	}

	var ivBytes []byte
	if len(a) >= 2 {
		var err error
		ivBytes, err = toBytes(a[1])
		if err != nil {
			return err
		}
		if len(ivBytes) != block.BlockSize() {
			return fmt.Errorf("iv length should be %d bytes, is %d bytes", block.BlockSize(), len(ivBytes))
		}
	} else {
		ivBytes = make([]byte, block.BlockSize())
	}

	bb, _, _, err := toBitBuf(c)
	if err != nil {
		return err
	}

	buf := &bytes.Buffer{}
	reader := &cipher.StreamReader{S: cipher.NewCTR(block, ivBytes), R: bb}
	if _, err := io.Copy(buf, reader); err != nil {
		return err
	}

	return buf.Bytes()
}
