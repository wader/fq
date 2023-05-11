package interp

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"math/big"
	"reflect"
	"strings"
	"time"

	"github.com/mitchellh/copystructure"
	"github.com/wader/fq/internal/bitioex"
	"github.com/wader/fq/internal/gojqex"
	"github.com/wader/fq/internal/ioex"
	"github.com/wader/fq/internal/mapstruct"
	"github.com/wader/fq/pkg/bitio"
	"github.com/wader/fq/pkg/decode"

	"github.com/wader/gojq"
)

func init() {
	RegisterFunc0("_registry", (*Interp)._registry)
	RegisterFunc1("_tovalue", (*Interp)._toValue)
	RegisterFunc2("_decode", (*Interp)._decode)
}

type expectedExtkeyError struct {
	Key string
}

func (err expectedExtkeyError) Error() string {
	return "expected a extkey but got: " + err.Key
}

// TODO: redo/rename
// used by _isDecodeValue
type DecodeValue interface {
	Value
	ToBinary

	DecodeValue() *decode.Value
}

func (i *Interp) _registry(c any) any {
	uniqueFormats := map[string]*decode.Format{}

	groups := map[string]any{}
	formats := map[string]any{}

	for _, g := range i.Registry.Groups() {
		var group []any

		for _, f := range g.Formats {
			group = append(group, f.Name)
			if _, ok := uniqueFormats[f.Name]; ok {
				continue
			}
			uniqueFormats[f.Name] = f
		}

		groups[g.Name] = group
	}

	for _, f := range uniqueFormats {
		vf := map[string]any{
			"name":                 f.Name,
			"description":          f.Description,
			"probe_order":          f.ProbeOrder,
			"root_name":            f.RootName,
			"root_array":           f.RootArray,
			"skip_decode_function": f.SkipDecodeFunction,
		}

		var dependenciesVs []any
		for _, d := range f.Dependencies {
			var dNamesVs []any
			for _, g := range d.Groups {
				dNamesVs = append(dNamesVs, g.Name)
			}
			dependenciesVs = append(dependenciesVs, dNamesVs)
		}
		if len(dependenciesVs) > 0 {
			vf["dependencies"] = dependenciesVs
		}
		var groupsVs []any
		for _, g := range f.Groups {
			groupsVs = append(groupsVs, g.Name)
		}
		if len(groupsVs) > 0 {
			vf["groups"] = groupsVs
		}
		if f.DefaultInArg != nil {
			doc := map[string]any{}
			st := reflect.TypeOf(f.DefaultInArg)
			for i := 0; i < st.NumField(); i++ {
				f := st.Field(i)
				if v, ok := f.Tag.Lookup("doc"); ok {
					doc[mapstruct.CamelToSnake(f.Name)] = v
				}
			}
			vf["decode_in_arg_doc"] = doc

			args, err := mapstruct.ToMap(f.DefaultInArg)
			if err != nil {
				return err
			}

			// filter out internal field without documentation
			for k := range args {
				if _, ok := doc[k]; !ok {
					delete(args, k)
				}
			}
			vf["decode_in_arg"] = gojqex.Normalize(args)
		}

		if f.Functions != nil {
			var ss []any
			for _, f := range f.Functions {
				ss = append(ss, f)
			}
			vf["functions"] = ss
		}

		formats[f.Name] = vf
	}

	var files []any
	for _, fs := range i.Registry.FSs {
		ventries := []any{}

		entries, err := fs.ReadDir(".")
		if err != nil {
			return err
		}

		for _, e := range entries {
			f, err := fs.Open(e.Name())
			if err != nil {
				return err
			}
			b, err := io.ReadAll(f)
			if err != nil {
				return err
			}

			ventries = append(ventries, map[string]any{
				"name": e.Name(),
				"data": string(b),
			})
		}

		files = append(files, ventries)
	}

	return map[string]any{
		"groups":  groups,
		"formats": formats,
		"files":   files,
	}
}

func (i *Interp) _toValue(c any, om map[string]any) any {
	return toValue(
		func() *Options {
			opts := OptionsFromValue(om)
			return &opts
		},
		c,
	)
}

type decodeOpts struct {
	Force    bool
	Progress string
	Remain   map[string]any `mapstruct:",remain"`
}

func (i *Interp) _decode(c any, format string, opts decodeOpts) any {
	var filename string

	// TODO: progress hack
	// would be nice to move all progress code into decode but it might be
	// tricky to keep track of absolute positions in the underlaying readers
	// when it uses BitBuf slices, maybe only in Pos()?
	if bbf, ok := c.(*openFile); ok {
		filename = bbf.filename

		if opts.Progress != "" {
			evalProgress := func(c any) {
				// {approx_read_bytes: 123, total_size: 123} | opts.Progress
				_, _ = i.EvalFuncValues(
					i.EvalInstance.Ctx,
					c,
					opts.Progress,
					nil,
					EvalOpts{output: ioex.DiscardCtxWriter{Ctx: i.EvalInstance.Ctx}},
				)
			}
			lastProgress := time.Now()
			bbf.progressFn = func(approxReadBytes, totalSize int64) {
				// make sure to not call too often as it's quite expensive
				n := time.Now()
				if n.Sub(lastProgress) < 200*time.Millisecond {
					return
				}
				lastProgress = n
				evalProgress(
					map[string]any{
						"approx_read_bytes": approxReadBytes,
						"total_size":        totalSize,
					},
				)
			}
			// when done decoding, tell progress function were done and disable it
			defer func() {
				bbf.progressFn = nil
				evalProgress(nil)
			}()
		}
	}

	bv, err := toBinary(c)
	if err != nil {
		return err
	}

	formatName, err := toString(format)
	if err != nil {
		return err
	}
	decodeGroup, err := i.Registry.Group(formatName)
	if err != nil {
		return err
	}

	dv, formatOut, err := decode.Decode(i.EvalInstance.Ctx, bv.br, decodeGroup,
		decode.Options{
			IsRoot:      true,
			FillGaps:    true,
			Force:       opts.Force,
			Range:       bv.r,
			Description: filename,
			ParseOptsFn: func(init any) any {
				v, err := copystructure.Copy(init)
				if err != nil {
					return nil
				}

				if len(opts.Remain) > 0 {
					if err := mapstruct.ToStruct(opts.Remain, &v); err != nil {
						// TODO: currently ignores failed struct mappings
						return nil
					}
				}
				// nil if same as init
				if reflect.DeepEqual(init, v) {
					return nil
				}

				return v
			},
		},
	)
	if dv == nil {
		var decodeFormatsErr decode.FormatsError
		if errors.As(err, &decodeFormatsErr) {
			var vs []any
			for _, fe := range decodeFormatsErr.Errs {
				vs = append(vs, fe.Value())
			}

			return valueError{vs}
		}

		return valueError{err}
	}

	var formatOutMap any

	if formatOut != nil {
		formatOutMap, err = mapstruct.ToMap(formatOut)
		if err != nil {
			return err
		}
	}

	return makeDecodeValueOut(dv, decodeValueValue, formatOutMap)
}

func valueKey(name string, a, b func(name string) any) any {
	if strings.HasPrefix(name, "_") {
		return a(name)
	}
	return b(name)
}
func valueHas(key any, a func(name string) any, b func(key any) any) any {
	stringKey, ok := key.(string)
	if ok && strings.HasPrefix(stringKey, "_") {
		if err, ok := a(stringKey).(error); ok {
			return err
		}
		return true
	}
	return b(key)
}

// TODO: make more efficient somehow? shallow values but might be hard
// when things like tovalue.key should behave like a jq value and not a decode value etc
func toValue(optsFn func() *Options, v any) any {
	nv, _ := gojqex.ToGoJQValueFn(v, func(v any) (any, bool) {
		switch v := v.(type) {
		case JQValueEx:
			if optsFn == nil {
				return v.JQValueToGoJQ(), true
			}
			return v.JQValueToGoJQEx(optsFn), true
		case gojq.JQValue:
			return v.JQValueToGoJQ(), true
		default:
			return v, true
		}
	})
	return nv
}

type decodeValueKind int

const (
	decodeValueValue decodeValueKind = iota
	decodeValueActual
	decodeValueSym
)

func makeDecodeValue(dv *decode.Value, kind decodeValueKind) any {
	return makeDecodeValueOut(dv, kind, nil)
}

func makeDecodeValueOut(dv *decode.Value, kind decodeValueKind, out any) any {
	switch vv := dv.V.(type) {
	case *decode.Compound:
		if vv.IsArray {
			return NewArrayDecodeValue(dv, out, vv)
		}
		return NewStructDecodeValue(dv, out, vv)

	case Scalarable:
		// TODO: rethink value/actual/sym handling
		var vvv any
		switch kind {
		case decodeValueValue:
			vvv = vv.ScalarValue()
		case decodeValueActual:
			vvv = vv.ScalarActual()
		case decodeValueSym:
			vvv = vv.ScalarSym()
		}

		switch vvv := vvv.(type) {
		case bitio.ReaderAtSeeker:
			// is lazy so that in situations where the decode value is only used to
			// create another binary we don't have to read and create a string, ex:
			// .unknown0 | tobytes[1:] | ...
			return decodeValue{
				JQValue: &gojqex.Lazy{
					Type:     "string",
					IsScalar: true,
					Fn: func() (gojq.JQValue, error) {
						buf := &bytes.Buffer{}
						vvvC, err := bitio.CloneReader(vvv)
						if err != nil {
							return nil, err
						}
						if _, err := bitioex.CopyBits(buf, vvvC); err != nil {
							return nil, err
						}
						return gojqex.String([]rune(buf.String())), nil
					},
				},
				decodeValueBase: decodeValueBase{dv: dv},
				bitsFormat:      true,
			}
		case bool:
			return decodeValue{
				JQValue:         gojqex.Boolean(vvv),
				decodeValueBase: decodeValueBase{dv: dv},
			}
		case int:
			return decodeValue{
				JQValue:         gojqex.Number{V: vvv},
				decodeValueBase: decodeValueBase{dv: dv},
			}
		case int64:
			return decodeValue{
				JQValue:         gojqex.Number{V: big.NewInt(vvv)},
				decodeValueBase: decodeValueBase{dv: dv},
			}
		case uint64:
			return decodeValue{
				JQValue:         gojqex.Number{V: new(big.Int).SetUint64(vvv)},
				decodeValueBase: decodeValueBase{dv: dv},
			}
		case float64:
			return decodeValue{
				JQValue:         gojqex.Number{V: vvv},
				decodeValueBase: decodeValueBase{dv: dv},
			}
		case string:
			return decodeValue{
				JQValue:         gojqex.String(vvv),
				decodeValueBase: decodeValueBase{dv: dv},
			}
		case []any:
			return decodeValue{
				JQValue:         gojqex.Array(vvv),
				decodeValueBase: decodeValueBase{dv: dv},
			}
		case map[string]any:
			return decodeValue{
				JQValue:         gojqex.Object(vvv),
				decodeValueBase: decodeValueBase{dv: dv},
			}
		case nil:
			return decodeValue{
				JQValue:         gojqex.Null{},
				decodeValueBase: decodeValueBase{dv: dv},
			}
		case *big.Int:
			return decodeValue{
				JQValue:         gojqex.Number{V: vvv},
				decodeValueBase: decodeValueBase{dv: dv},
			}
		case Binary:
			return vvv

		default:
			panic(fmt.Sprintf("unreachable vv %#+v", vvv))
		}
	default:
		panic(fmt.Sprintf("unreachable dv %#+v", dv))
	}
}

type decodeValueBase struct {
	dv  *decode.Value
	out any
}

func (dvb decodeValueBase) DecodeValue() *decode.Value {
	return dvb.dv
}

func (dvb decodeValueBase) Display(w io.Writer, opts Options) error { return dump(dvb.dv, w, opts) }
func (dvb decodeValueBase) ToBinary() (Binary, error) {
	return Binary{br: dvb.dv.RootReader, r: dvb.dv.InnerRange(), unit: 8}, nil
}
func (decodeValueBase) ExtType() string { return "decode_value" }
func (dvb decodeValueBase) ExtKeys() []string {
	kv := []string{
		"_actual",
		"_bits",
		"_buffer_root",
		"_bytes",
		"_description",
		"_format_root",
		"_gap",
		"_len",
		"_name",
		"_parent",
		"_path",
		"_root",
		"_start",
		"_stop",
		"_sym",
	}

	if _, ok := dvb.dv.V.(*decode.Compound); ok {
		kv = append(kv,
			"_error",
			"_format",
			"_out",
		)

		if dvb.dv.Index != -1 {
			kv = append(kv, "_index")
		}
	}

	return kv
}

func (dvb decodeValueBase) JQValueKey(name string) any {
	dv := dvb.dv

	switch name {
	case "_actual":
		switch dv.V.(type) {
		case Scalarable:
			return makeDecodeValue(dv, decodeValueActual)
		default:
			return nil
		}
	case "_bits":
		return Binary{
			br:   dv.RootReader,
			r:    dv.Range,
			unit: 1,
		}
	case "_buffer_root":
		// TODO: rename?
		return makeDecodeValue(dv.BufferRoot(), decodeValueValue)
	case "_bytes":
		return Binary{
			br:   dv.RootReader,
			r:    dv.Range,
			unit: 8,
		}
	case "_description":
		switch vv := dv.V.(type) {
		case *decode.Compound:
			if vv.Description == "" {
				return nil
			}
			return vv.Description
		case Scalarable:
			desc := vv.ScalarDescription()
			if desc == "" {
				return nil
			}
			return desc
		default:
			return nil
		}
	case "_format_root":
		// TODO: rename?
		return makeDecodeValue(dv.FormatRoot(), decodeValueValue)
	case "_gap":
		switch vv := dv.V.(type) {
		case Scalarable:
			return vv.ScalarIsGap()
		default:
			return false
		}
	case "_len":
		return big.NewInt(dv.Range.Len)
	case "_name":
		return dv.Name
	case "_parent":
		if dv.Parent == nil {
			return nil
		}
		return makeDecodeValue(dv.Parent, decodeValueValue)
	case "_path":
		return valuePath(dv)
	case "_root":
		return makeDecodeValue(dv.Root(), decodeValueValue)
	case "_start":
		return big.NewInt(dv.Range.Start)
	case "_stop":
		return big.NewInt(dv.Range.Stop())
	case "_sym":
		switch dv.V.(type) {
		case Scalarable:
			return makeDecodeValue(dv, decodeValueSym)
		default:
			return nil
		}

	case "_error":
		var formatErr decode.FormatError
		if errors.As(dv.Err, &formatErr) {
			return formatErr.Value()
		}
		return nil
	case "_format":
		if dv.Format != nil {
			return dv.Format.Name
		}
		return nil
	case "_out":
		return dvb.out

	case "_index":
		if dv.Index != -1 {
			return dv.Index
		}
	}

	return expectedExtkeyError{Key: name}
}

var _ DecodeValue = decodeValue{}

type decodeValue struct {
	gojq.JQValue
	decodeValueBase
	bitsFormat bool
}

func (v decodeValue) JQValueKey(name string) any {
	return valueKey(name, v.decodeValueBase.JQValueKey, v.JQValue.JQValueKey)
}
func (v decodeValue) JQValueHas(key any) any {
	return valueHas(key, v.decodeValueBase.JQValueKey, v.JQValue.JQValueHas)
}
func (v decodeValue) JQValueToGoJQEx(optsFn func() *Options) any {
	if !v.bitsFormat {
		return v.JQValueToGoJQ()
	}

	bv, err := v.decodeValueBase.ToBinary()
	if err != nil {
		return err
	}
	br, err := bv.toReader()
	if err != nil {
		return err
	}

	brC, err := bitio.CloneReaderAtSeeker(br)
	if err != nil {
		return err
	}

	s, err := optsFn().BitsFormatFn(brC)
	if err != nil {
		return err
	}
	return s
}

// decode value array

var _ DecodeValue = ArrayDecodeValue{}

type ArrayDecodeValue struct {
	gojqex.Base
	decodeValueBase
	*decode.Compound
}

func NewArrayDecodeValue(dv *decode.Value, out any, c *decode.Compound) ArrayDecodeValue {
	return ArrayDecodeValue{
		decodeValueBase: decodeValueBase{dv: dv, out: out},
		Base:            gojqex.Base{Typ: gojq.JQTypeArray},
		Compound:        c,
	}
}

func (v ArrayDecodeValue) JQValueKey(name string) any {
	return valueKey(name, v.decodeValueBase.JQValueKey, v.Base.JQValueKey)
}
func (v ArrayDecodeValue) JQValueSliceLen() any { return len(v.Compound.Children) }
func (v ArrayDecodeValue) JQValueLength() any   { return len(v.Compound.Children) }
func (v ArrayDecodeValue) JQValueIndex(index int) any {
	// -1 outside after string, -2 outside before string
	if index < 0 {
		return nil
	}
	return makeDecodeValue((v.Compound.Children)[index], decodeValueValue)
}
func (v ArrayDecodeValue) JQValueSlice(start int, end int) any {
	vs := make([]any, end-start)
	for i, e := range (v.Compound.Children)[start:end] {
		vs[i] = makeDecodeValue(e, decodeValueValue)
	}
	return vs
}
func (v ArrayDecodeValue) JQValueEach() any {
	props := make([]gojq.PathValue, len(v.Compound.Children))
	for i, f := range v.Compound.Children {
		props[i] = gojq.PathValue{Path: i, Value: makeDecodeValue(f, decodeValueValue)}
	}
	return props
}
func (v ArrayDecodeValue) JQValueKeys() any {
	vs := make([]any, len(v.Compound.Children))
	for i := range v.Compound.Children {
		vs[i] = i
	}
	return vs
}
func (v ArrayDecodeValue) JQValueHas(key any) any {
	return valueHas(
		key,
		v.decodeValueBase.JQValueKey,
		func(key any) any {
			intKey, ok := key.(int)
			if !ok {
				return gojqex.HasKeyTypeError{L: gojq.JQTypeArray, R: fmt.Sprintf("%v", key)}
			}
			return intKey >= 0 && intKey < len(v.Compound.Children)
		})
}
func (v ArrayDecodeValue) JQValueToGoJQEx(optsFn func() *Options) any {
	opts := optsFn()

	vs := make([]any, 0, len(v.Compound.Children))
	for _, f := range v.Compound.Children {
		switch s := f.V.(type) {
		case Scalarable:
			if s.ScalarIsGap() && opts.SkipGaps {
				// skip, note for arrays this will affect indexes
				continue
			}
		}

		vs = append(vs, makeDecodeValue(f, decodeValueValue))
	}
	return vs
}
func (v ArrayDecodeValue) JQValueToGoJQ() any {
	return v.JQValueToGoJQEx(func() *Options { return &Options{} })
}

// decode value struct

var _ DecodeValue = StructDecodeValue{}

type StructDecodeValue struct {
	gojqex.Base
	decodeValueBase
	*decode.Compound
}

func NewStructDecodeValue(dv *decode.Value, out any, c *decode.Compound) StructDecodeValue {
	return StructDecodeValue{
		decodeValueBase: decodeValueBase{dv: dv, out: out},
		Base:            gojqex.Base{Typ: gojq.JQTypeObject},
		Compound:        c,
	}
}

func (v StructDecodeValue) JQValueLength() any   { return len(v.Compound.Children) }
func (v StructDecodeValue) JQValueSliceLen() any { return len(v.Compound.Children) }
func (v StructDecodeValue) JQValueKey(name string) any {
	if strings.HasPrefix(name, "_") {
		return v.decodeValueBase.JQValueKey(name)
	}
	if v.Compound.ByName != nil {
		if f, ok := v.Compound.ByName[name]; ok {
			return makeDecodeValue(f, decodeValueValue)
		}
	}

	return nil
}
func (v StructDecodeValue) JQValueEach() any {
	props := make([]gojq.PathValue, len(v.Compound.Children))
	for i, f := range v.Compound.Children {
		props[i] = gojq.PathValue{Path: f.Name, Value: makeDecodeValue(f, decodeValueValue)}
	}
	return props
}
func (v StructDecodeValue) JQValueKeys() any {
	vs := make([]any, len(v.Compound.Children))
	for i, f := range v.Compound.Children {
		vs[i] = f.Name
	}
	return vs
}
func (v StructDecodeValue) JQValueHas(key any) any {
	return valueHas(
		key,
		v.decodeValueBase.JQValueKey,
		func(key any) any {
			stringKey, ok := key.(string)
			if !ok {
				return gojqex.HasKeyTypeError{L: gojq.JQTypeObject, R: fmt.Sprintf("%v", key)}
			}
			for _, f := range v.Compound.Children {
				if f.Name == stringKey {
					return true
				}
			}
			return false
		},
	)
}
func (v StructDecodeValue) JQValueToGoJQEx(optsFn func() *Options) any {
	opts := optsFn()

	vm := make(map[string]any, len(v.Compound.Children))
	for _, f := range v.Compound.Children {
		switch s := f.V.(type) {
		case Scalarable:
			if s.ScalarIsGap() && opts.SkipGaps {
				continue
			}
		}

		vm[f.Name] = makeDecodeValue(f, decodeValueValue)
	}
	return vm
}
func (v StructDecodeValue) JQValueToGoJQ() any {
	return v.JQValueToGoJQEx(func() *Options { return &Options{} })
}
