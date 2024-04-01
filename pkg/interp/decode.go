package interp

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"math/big"
	"reflect"
	"time"

	"github.com/mitchellh/copystructure"
	"github.com/wader/fq/internal/bitiox"
	"github.com/wader/fq/internal/gojqx"
	"github.com/wader/fq/internal/iox"
	"github.com/wader/fq/internal/mapstruct"
	"github.com/wader/fq/pkg/bitio"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/scalar"

	"github.com/wader/gojq"
)

func init() {
	RegisterFunc0("_registry", (*Interp)._registry)
	RegisterFunc1("_tovalue", (*Interp)._toValue)
	RegisterFunc2("_decode", (*Interp)._decode)
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
			vf["decode_in_arg"] = gojqx.Normalize(args)
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
	opts, err := OptionsFromValue(om)
	if err != nil {
		return err
	}

	v, err := toValue(func() (*Options, error) { return opts, nil }, c)
	if err != nil {
		return err
	}
	return v
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
	// tricky to keep track of absolute positions in the underlying readers
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
					EvalOpts{output: iox.DiscardCtxWriter{Ctx: i.EvalInstance.Ctx}},
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

func valueOrFallbackKey(name string, baseKey func(name string) any, valueHas func(key any) any, valueKey func(name string) any) any {
	v := valueHas(name)
	if b, ok := v.(bool); ok && b {
		return valueKey(name)
	}
	return baseKey(name)
}
func valueOrFallbackHas(key any, baseHas func(key any) any, valueHas func(key any) any) any {
	v := valueHas(key)
	if b, ok := v.(bool); ok && !b {
		return baseHas(key)
	}
	return v
}

// TODO: make more efficient somehow? shallow values but might be hard
// when things like tovalue.key should behave like a jq value and not a decode value etc
func toValue(optsFn func() (*Options, error), v any) (any, error) {
	return gojqx.ToGoJQValueFn(v, func(v any) (any, error) {
		switch v := v.(type) {
		case JQValueEx:
			if optsFn == nil {
				return v.JQValueToGoJQ(), nil
			}
			return v.JQValueToGoJQEx(optsFn), nil
		case gojq.JQValue:
			return v.JQValueToGoJQ(), nil
		default:
			return v, nil
		}
	})
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

	case scalar.Scalarable:
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
				JQValue: &gojqx.Lazy{
					Type:     "string",
					IsScalar: true,
					Fn: func() (gojq.JQValue, error) {
						buf := &bytes.Buffer{}
						vvvC, err := bitio.CloneReader(vvv)
						if err != nil {
							return nil, err
						}
						if _, err := bitiox.CopyBits(buf, vvvC); err != nil {
							return nil, err
						}
						return gojqx.String([]rune(buf.String())), nil
					},
				},
				decodeValueBase: decodeValueBase{dv: dv},
				isRaw:           true,
			}
		case bool:
			return decodeValue{
				JQValue:         gojqx.Boolean(vvv),
				decodeValueBase: decodeValueBase{dv: dv},
			}
		case int:
			return decodeValue{
				JQValue:         gojqx.Number{V: vvv},
				decodeValueBase: decodeValueBase{dv: dv},
			}
		case int64:
			return decodeValue{
				JQValue:         gojqx.Number{V: big.NewInt(vvv)},
				decodeValueBase: decodeValueBase{dv: dv},
			}
		case uint64:
			return decodeValue{
				JQValue:         gojqx.Number{V: new(big.Int).SetUint64(vvv)},
				decodeValueBase: decodeValueBase{dv: dv},
			}
		case float64:
			return decodeValue{
				JQValue:         gojqx.Number{V: vvv},
				decodeValueBase: decodeValueBase{dv: dv},
			}
		case string:
			return decodeValue{
				JQValue:         gojqx.String(vvv),
				decodeValueBase: decodeValueBase{dv: dv},
			}
		case []any:
			return decodeValue{
				JQValue:         gojqx.Array(vvv),
				decodeValueBase: decodeValueBase{dv: dv},
			}
		case map[string]any:
			return decodeValue{
				JQValue:         gojqx.Object(vvv),
				decodeValueBase: decodeValueBase{dv: dv},
			}
		case nil:
			return decodeValue{
				JQValue:         gojqx.Null{},
				decodeValueBase: decodeValueBase{dv: dv},
			}
		case *big.Int:
			return decodeValue{
				JQValue:         gojqx.Number{V: vvv},
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

func (dvb decodeValueBase) Display(w io.Writer, opts *Options) error { return dump(dvb.dv, w, opts) }
func (dvb decodeValueBase) ToBinary() (Binary, error) {
	if s, ok := dvb.dv.V.(scalar.Scalarable); ok && s.ScalarFlags().IsSynthetic() {
		return Binary{}, fmt.Errorf("synthetic value can't be a binary")
	}
	return Binary{br: dvb.dv.RootReader, r: dvb.dv.InnerRange(), unit: 8}, nil
}
func (decodeValueBase) ExtType() string { return "decode_value" }
func (dvb decodeValueBase) ExtKeys() []string {
	return []string{
		"_actual",
		"_bits",
		"_buffer_root",
		"_bytes",
		"_description",
		"_error",
		"_format_root",
		"_format",
		"_gap",
		"_index",
		"_len",
		"_name",
		"_out",
		"_parent",
		"_path",
		"_root",
		"_start",
		"_stop",
		"_sym",
	}
}

func (dvb decodeValueBase) JQValueHas(key any) any {
	name, ok := key.(string)
	if !ok {
		return false
	}

	switch name {
	case "_actual",
		"_bits",
		"_buffer_root",
		"_bytes",
		"_description",
		"_error",
		"_format_root",
		"_format",
		"_gap",
		"_index",
		"_len",
		"_name",
		"_out",
		"_parent",
		"_path",
		"_root",
		"_start",
		"_stop",
		"_sym":
		return true
	}

	return false
}

func (dvb decodeValueBase) JQValueKey(name string) any {
	dv := dvb.dv

	switch name {
	case "_actual":
		switch dv.V.(type) {
		case scalar.Scalarable:
			return makeDecodeValue(dv, decodeValueActual)
		default:
			return nil
		}
	case "_bits":
		if s, ok := dv.V.(scalar.Scalarable); ok && s.ScalarFlags().IsSynthetic() {
			return nil
		}
		return Binary{
			br:   dv.RootReader,
			r:    dv.Range,
			unit: 1,
		}
	case "_buffer_root":
		// TODO: rename?
		return makeDecodeValue(dv.BufferRoot(), decodeValueValue)
	case "_bytes":
		if s, ok := dv.V.(scalar.Scalarable); ok && s.ScalarFlags().IsSynthetic() {
			return nil
		}
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
		case scalar.Scalarable:
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
		case scalar.Scalarable:
			return vv.ScalarFlags().IsGap()
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
		case scalar.Scalarable:
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

	return nil
}

var _ DecodeValue = decodeValue{}

type decodeValue struct {
	gojq.JQValue
	decodeValueBase
	isRaw bool
}

func (v decodeValue) JQValueKey(name string) any {
	return valueOrFallbackKey(name, v.decodeValueBase.JQValueKey, v.JQValue.JQValueHas, v.JQValue.JQValueKey)
}
func (v decodeValue) JQValueHas(key any) any {
	return valueOrFallbackHas(key, v.decodeValueBase.JQValueHas, v.JQValue.JQValueHas)
}
func (v decodeValue) JQValueToGoJQEx(optsFn func() (*Options, error)) any {
	if !v.isRaw {
		return v.JQValueToGoJQ()
	}

	if s, ok := v.dv.V.(scalar.Scalarable); ok && !s.ScalarFlags().IsSynthetic() {
		bv, err := v.ToBinary()
		if err != nil {
			return err
		}
		return bv.JQValueToGoJQEx(optsFn)
	}

	return v.JQValueToGoJQ()

}

// decode value array

var _ DecodeValue = ArrayDecodeValue{}

type ArrayDecodeValue struct {
	gojqx.Base
	decodeValueBase
	*decode.Compound
}

func NewArrayDecodeValue(dv *decode.Value, out any, c *decode.Compound) ArrayDecodeValue {
	return ArrayDecodeValue{
		decodeValueBase: decodeValueBase{dv: dv, out: out},
		Base:            gojqx.Base{Typ: gojq.JQTypeArray},
		Compound:        c,
	}
}

func (v ArrayDecodeValue) JQValueKey(name string) any {
	return valueOrFallbackKey(name, v.decodeValueBase.JQValueKey, v.Base.JQValueHas, v.Base.JQValueKey)
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
	return valueOrFallbackHas(
		key,
		v.decodeValueBase.JQValueHas,
		func(key any) any {
			intKey, ok := key.(int)
			if !ok {
				return gojqx.HasKeyTypeError{L: gojq.JQTypeArray, R: fmt.Sprintf("%v", key)}
			}
			return intKey >= 0 && intKey < len(v.Compound.Children)
		})
}
func (v ArrayDecodeValue) JQValueToGoJQEx(optsFn func() (*Options, error)) any {
	opts, err := optsFn()
	if err != nil {
		return err
	}

	vs := make([]any, 0, len(v.Compound.Children))
	for _, f := range v.Compound.Children {
		switch s := f.V.(type) {
		case scalar.Scalarable:
			if s.ScalarFlags().IsGap() && opts.SkipGaps {
				// skip, note for arrays this will affect indexes
				continue
			}
		}

		vs = append(vs, makeDecodeValue(f, decodeValueValue))
	}
	return vs
}
func (v ArrayDecodeValue) JQValueToGoJQ() any {
	return v.JQValueToGoJQEx(func() (*Options, error) { return &Options{}, nil })
}

// decode value struct

var _ DecodeValue = StructDecodeValue{}

type StructDecodeValue struct {
	gojqx.Base
	decodeValueBase
	*decode.Compound
}

func NewStructDecodeValue(dv *decode.Value, out any, c *decode.Compound) StructDecodeValue {
	return StructDecodeValue{
		decodeValueBase: decodeValueBase{dv: dv, out: out},
		Base:            gojqx.Base{Typ: gojq.JQTypeObject},
		Compound:        c,
	}
}

func (v StructDecodeValue) JQValueLength() any   { return len(v.Compound.Children) }
func (v StructDecodeValue) JQValueSliceLen() any { return len(v.Compound.Children) }
func (v StructDecodeValue) JQValueKey(name string) any {
	return valueOrFallbackKey(
		name,
		v.decodeValueBase.JQValueKey,
		func(key any) any {
			stringKey, ok := key.(string)
			if !ok {
				return false
			}
			if v.Compound.ByName != nil {
				if _, ok := v.Compound.ByName[stringKey]; ok {
					return true
				}
			}
			return false
		},
		func(name string) any {
			if v.Compound.ByName != nil {
				if f, ok := v.Compound.ByName[name]; ok {
					return makeDecodeValue(f, decodeValueValue)
				}
			}

			return nil
		},
	)
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
	return valueOrFallbackHas(
		key,
		v.decodeValueBase.JQValueHas,
		func(key any) any {
			stringKey, ok := key.(string)
			if !ok {
				return gojqx.HasKeyTypeError{L: gojq.JQTypeObject, R: fmt.Sprintf("%v", key)}
			}

			if v.Compound.ByName != nil {
				if _, ok := v.Compound.ByName[stringKey]; ok {
					return true
				}
			}

			return false
		},
	)
}
func (v StructDecodeValue) JQValueToGoJQEx(optsFn func() (*Options, error)) any {
	opts, err := optsFn()
	if err != nil {
		return err
	}

	vm := make(map[string]any, len(v.Compound.Children))
	for _, f := range v.Compound.Children {
		switch s := f.V.(type) {
		case scalar.Scalarable:
			if s.ScalarFlags().IsGap() && opts.SkipGaps {
				continue
			}
		}

		vm[f.Name] = makeDecodeValue(f, decodeValueValue)
	}
	return vm
}
func (v StructDecodeValue) JQValueToGoJQ() any {
	return v.JQValueToGoJQEx(func() (*Options, error) { return &Options{}, nil })
}
