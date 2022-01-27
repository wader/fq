package kaitai

// TODO: _io
// TODO: instance/seq array index
// TODO: add _meta?
// TODO: values tree somehow, ksexpr interface?
// TODO: prompt format name somehow?
// TODO: dump title/description?
// TODO: meta, endian, switch-on
// TODO: prompt, per format name?
// TODO: sizeof ceiled to bytes
// TODO: bitsizeof
// TODO: _index
// TODO: fq -o source=@format/kaitai/testdata/test.ksy -d kaitai d file
// TODO: fq -d format/kaitai/testdata/test.ksy d file
// TODO: error no default endianness
// TODO: ternary true/false same type
// TODO typed "no value"?

// resolve:

// enum:
//   <enum-name>
//   (<type-name>::)*<enum-name>

// type: <type>
// enum: <enum>
// value: <expr>
// switch-on: <expr>
//   cases:
//     <expr>: <type>
// <type>
//   (<type-name>::)*<type-name>
// <enum-value>
//   (<type-name>::)*<enum-name>::<name>
// <expr>
//   1+2
//   <enum-value>
//   <path>
// <path>
//   (<path>.)<path>
//   <instance>
//   <seq>
//   _parent
//   _root
//   _io

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"slices"
	"strings"

	"github.com/wader/fq/format"
	"github.com/wader/fq/format/kaitai/ksexpr"
	"github.com/wader/fq/format/kaitai/schema"
	"github.com/wader/fq/format/kaitai/schema/primitive"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
	"github.com/wader/fq/pkg/scalar"
)

func init() {
	interp.RegisterFormat(
		format.Kaitai,
		&decode.Format{
			Description:  "Kaitai struct declarative language",
			DecodeFn:     kaitaiDecode,
			Groups:       []*decode.Group{format.Probe_Args},
			DefaultInArg: format.Kaitai_In{},
		})
}

func addStrNonEmpty(d *decode.D, name string, v string) {
	if v != "" {
		d.FieldValueStr(name, v)
	}
}

func decodeEndian(c decode.Endian, e primitive.Endianess) decode.Endian {
	switch e {
	case primitive.CurrentEndian:
		return c
	case primitive.LE:
		return decode.LittleEndian
	case primitive.BE:
		return decode.BigEndian
	default:
		panic("unreachable")
	}
}

type typeInstance struct {
	schemaType *schema.Type
	parent     *typeInstance
	root       *typeInstance

	d *decode.D

	// instances map[string]any

	fields map[string]any

	// TODO: used by array "_"
	last any

	// seq    map[string]any
	// repeat []any // TODO: hmm
}

func (ti *typeInstance) Fatalf(format string, a ...any) {
	ti.d.Fatalf("%s: "+format, append([]any{ti.path()}, a...)...)
}

func (ti *typeInstance) Printf(format string, v ...any) {
	log.Printf("%s: "+format, append([]any{ti.path()}, v...)...)
}

// TODO: array index? not /?
func (ti *typeInstance) path() string {
	var ps []string
	for c := ti; c.parent != nil; c = c.parent {
		ps = append(ps, c.schemaType.ID)
	}
	slices.Reverse(ps)
	return "/" + strings.Join(ps, "/")
}

func ksExprField(d *decode.D, name string, v any) {
	switch v := v.(type) {
	case ksexpr.Integer:
		d.FieldValueSint(name, int64(v))
	case ksexpr.BigInt:
		d.FieldValueBigInt(name, v.V)
	case ksexpr.Boolean:
		d.FieldValueBool(name, bool(v))
	case ksexpr.Float:
		d.FieldValueFlt(name, float64(v))
	case ksexpr.String:
		d.FieldValueStr(name, string(v))
	case ksexpr.Array:
		d.FieldArray(name, func(d *decode.D) {
			for _, ve := range v {
				ksExprField(d, name, ve)
			}
		})
	case ksexpr.Object:
		// TODO: not possible?
		d.FieldStruct(name, func(d *decode.D) {
			for k, ve := range v {
				ksExprField(d, k, ve)
			}
		})
	case ksexpr.Enum:
		// TODO: ns

		switch vv := v.V.(type) {
		case ksexpr.Integer:
			d.FieldValueSint(name, int64(vv), scalar.SintFn(func(s scalar.Sint) (scalar.Sint, error) {
				s.Sym = v.Name
				return s, nil
			}))
		default:
			// TODO: bigint
			panic("unreachable")
		}

	case nil:
		// TODO: show nil option, keep value in other tree?
		d.FieldValueAny(name, nil)

	default:
		panic(fmt.Sprintf("unreachable %#+v", v))
	}
}

type ksexprType struct {
	ti  *typeInstance
	typ *schema.Type
}

func (kt ksexprType) KSExprCall(ns []string, name string, args []any) (any, error) {
	kt.ti.Printf("ksexprType KSExprCall ns: %#+v name=%s\n", ns, name)

	// TODO: refactor parser?
	vns := append([]string{}, ns...)
	vns = append(vns, name)

	v, ok := kt.typ.ResolveType(vns)
	if !ok {
		// log.Println("  not found")
		return nil, fmt.Errorf("%s::%s", strings.Join(ns, "::"), name)
	}
	// log.Printf("  found v: %#+v\n", v)

	return v, nil
}

type ksexprEnum struct {
	ti  *typeInstance
	typ *schema.Type
}

func (ke ksexprEnum) KSExprCall(ns []string, name string, args []any) (any, error) {
	ke.ti.Printf("ksexprEnum KSExprCall ns: %#+v name=%s\n", ns, name)

	// TODO: refactor parser?
	vns := append([]string{}, ns...)
	vns = append(vns, name)

	v, ok := ke.typ.ResolveEnum(vns)
	if !ok {
		// log.Println("  not found")
		return nil, fmt.Errorf("%s::%s", strings.Join(ns, "::"), name)
	}
	// log.Printf("  found v: %#+v\n", v)

	return v, nil
}

func (ti *typeInstance) KSExprCall(ns []string, name string, args []any) (any, error) {
	ti.Printf("typeInstance KSExprCall ns %v name: %#+v\n", ns, name)

	if ns != nil {

		var tns []string
		tns = append(tns, ns...)
		tns = append(tns, name)

		ti.Printf("  tns: %#+v\n", tns)

		// tt, ok := ti.schemaType.ResolveType(tns)
		// if ok {
		// 	ti.Printf(" FOUND TYPE tt: %#+v\n", tt)
		// 	return tt, nil
		// }
		et, ok := ti.schemaType.ResolveEnum(ns)
		if ok {
			// ti.Printf(" FOUND ENUM et: %#+v\n", et)

			if v, ok := et.FromID[name]; ok {
				return ksexpr.Enum{Name: name, V: v}, nil
			}

			return et, nil
		}

		panic("hmm")

	}

	switch name {
	case "_":
		ti.Printf("ti.parent: %#+v\n", ti.parent)
		ti.Printf("ti.last: %#+v\n", ti.last)

		if ti.last == nil {
			ti.Fatalf("no last")
		}

		ti.Printf("  return ti.last: %#+v\n", ti.last)

		return ti.last, nil

	case "_io":
		// TODO: some io object?
		return ti, nil
	case "_parent":
		// TODO: no parent?
		return ti.parent, nil
	case "_root":
		ti.Printf("  return ti.root: %#+v\n", ti.root)

		return ti.root, nil
	case "size":
		return ti.d.Len() / 8, nil
	case "eof":
		ti.Printf(" k.d.End(): %#+v\n", ti.d.End())
		return ti.d.End(), nil
	default:

		ti.Printf("ti: %#+v\n", ti)

		ti.Printf("  ti.fields: %#+v\n", ti.fields)

		// TODO: KSExprIndex

		if v, ok := ti.fields[name]; ok {
			ti.Printf("  found as field: %#+v\n", name)
			return v, nil
		}

		// TODO: detect loop
		// TODO: in seq and in instace? look in parents?
		// TODO: if parent == nil?

		tst := ti.schemaType
		t, ok := tst.Instances[name]
		if !ok {
			ti.Printf("  NOT FOUND\n")
			ti.Printf("   tst.Instances: %#+v\n", tst.Instances)
			return nil, fmt.Errorf("instance %q not found", name)
		}

		tti := &typeInstance{
			schemaType: t,
			parent:     ti,
			root:       ti.root,
			d:          ti.d,

			fields: map[string]any{},
		}
		// TODO: why not ti.parent.d here and above?
		v := tti.decode(ti.d)
		// TODO: already exist? lint check?
		// ti.fields[t.ID] = v

		return v, nil
	}
}

func (ti *typeInstance) eval(name string, exprSource string, e *schema.Expr) (any, error) {
	v, err := e.KSExpr.Eval(ti)
	if err != nil {
		return nil, fmt.Errorf("%s: %s: %s: %w", name, exprSource, e.Str, err)
	}
	return v, nil
}

func (ti *typeInstance) mustEval(name string, exprSource string, e *schema.Expr) any {
	v, err := ti.eval(name, exprSource, e)
	if err != nil {
		ti.Fatalf("%s", err)
	}
	return v
}

func (ti *typeInstance) evalInt(name string, exprSource string, e *schema.Expr) (int, error) {
	sv, err := ti.eval(name, exprSource, e)
	if err != nil {
		return 0, err
	}
	s, ok := ksexpr.ToInt(sv)
	if !ok {
		return 0, fmt.Errorf("%s: %s: %s: did not evaluate to an integer: %v", name, exprSource, e.Str, s)
	}
	return s, nil
}

func (ti *typeInstance) mustEvalInt(name string, exprSource string, e *schema.Expr) int {
	ti.Printf("  mustEvalInt %s exprSource: %#+v %q\n", name, exprSource, e.Str)
	s, err := ti.evalInt(name, exprSource, e)
	if err != nil {
		ti.Fatalf("%s", err)
		panic("unreachable")
	}
	return s
}

func (ti *typeInstance) evalInt64(name string, exprSource string, e *schema.Expr) (int64, error) {
	sv, err := ti.eval(name, exprSource, e)
	if err != nil {
		return 0, err
	}
	s, ok := ksexpr.ToInt64(sv)
	if !ok {
		return 0, fmt.Errorf("%s: %s: %s: did not evaluate to an integer: %v", name, exprSource, e.Str, s)
	}
	return s, nil
}

func (ti *typeInstance) mustEvalInt64(name string, exprSource string, e *schema.Expr) int64 {
	ti.Printf("  mustEvalInt %s exprSource: %#+v %q\n", name, exprSource, e.Str)
	s, err := ti.evalInt64(name, exprSource, e)
	if err != nil {
		ti.Fatalf("%s", err)
		panic("unreachable")
	}
	return s
}

func (ti *typeInstance) mustEvalBool(name string, exprSource string, e *schema.Expr) bool {
	s := ti.parent.mustEval(name, exprSource, e)
	switch s := s.(type) {
	case ksexpr.Boolean:
		return bool(s)
	default:
		ti.Fatalf("%s: %s: %s: did not evaluated to a boolean: %s", name, exprSource, s)
		panic("unreachable")
	}
}

func contentsByteSize(c []any) (int, error) {
	s := 0

	for _, e := range c {
		switch e := e.(type) {
		case string:
			s += len(e)
		case int:
			if e < 0 || e > 255 {
				return 0, fmt.Errorf("contents array has invalid non-byte integer: %d", e)
			}
			s++
		default:
			return 0, fmt.Errorf("contents array as invalid value: %v", e)
		}
	}

	return s, nil
}

func (ti *typeInstance) decodePrimitiveType(d *decode.D, tst *schema.Type, pt primitive.Type) any {
	ti.Printf("  typ is primitive %d\n", d.Pos())

	if pt.BitAlign != 0 {
		d.BitEndianAlign()
	}

	var typEnum *schema.Enum
	if tst.Enum != nil {
		enumExprV, err := tst.Enum.KSExpr.Eval(ksexprEnum{ti, tst})
		if err != nil {
			ti.Fatalf("failed to resolved enum: %s", err)
		}

		var ok bool
		typEnum, ok = enumExprV.(*schema.Enum)
		if !ok {
			ti.Fatalf("enum resolved to non-enum: %#+v", enumExprV)
		}

		ti.Printf("  se: %#+v\n", typEnum)
	}

	switch pt.Encoding {
	case primitive.Bool:
		return d.FieldBoolE(tst.ID, decodeEndian(d.BitEndian, pt.Endian))

	case primitive.Bits,
		primitive.Unsigned:
		var mappers []scalar.UintMapper

		e := d.Endian
		if pt.Encoding == primitive.Bits {
			e = d.BitEndian
		}

		if typEnum != nil {
			mappers = append(mappers, scalar.UintFn(func(s scalar.Uint) (scalar.Uint, error) {
				if v, ok := typEnum.ToID[ksexpr.ToValue(s.Actual)]; ok {
					// s.Sym = fmt.Sprintf("%s::%s::%s", tst.Root.Meta.ID, tst.Enum.Str, v)
					s.Sym = v
				}
				return s, nil
			}))
		}

		v := d.FieldUE(tst.ID, pt.BitSize, decodeEndian(e, pt.Endian), mappers...)

		if typEnum != nil {
			if name, ok := typEnum.ToID[ksexpr.ToValue(v)]; ok {
				return ksexpr.Enum{Name: name, V: v}
			}
		}

		return v

	case primitive.Signed:
		var mappers []scalar.SintMapper

		if typEnum != nil {
			mappers = append(mappers, scalar.SintFn(func(s scalar.Sint) (scalar.Sint, error) {
				if v, ok := typEnum.ToID[ksexpr.ToValue(s.Actual)]; ok {
					// s.Sym = fmt.Sprintf("%s::%s::%s", tst.Root.Meta.ID, tst.Enum.Str, v)
					s.Sym = v
				}
				return s, nil
			}))
		}

		return d.FieldSE(tst.ID, pt.BitSize, decodeEndian(d.Endian, pt.Endian), mappers...)

	case primitive.Float:
		return d.FieldFE(tst.ID, pt.BitSize, decodeEndian(d.Endian, pt.Endian))

	case primitive.Bytes:
		switch {
		case tst.Size != nil:
			// TODO: parent?
			s := ti.parent.mustEvalInt(tst.ID, "size", tst.Size)
			ti.Printf("  %s: SIZE: s: %#+v\n", tst.ID, s)
			return d.FieldRawLen(tst.ID, int64(s)*8)

		case tst.SizeEOS:
			// TODO: error not byte aligned?
			return d.FieldRawLen(tst.ID, d.BitsLeft())

		case tst.Contents != nil:
			switch c := tst.Contents.(type) {
			case []any:
				// TODO: move schema parse if constant?
				l, err := contentsByteSize(c)
				if err != nil {
					d.Fatalf("%s: %s", tst.ID, err)
				}

				return d.FieldRawLen(tst.ID, int64(l)*8)

			case string:
				l := len(c)

				return d.FieldRawLen(tst.ID, int64(l)*8)

			default:
				panic("unreachable")
			}
		default:
			panic("unreachable")
		}

	case primitive.Str:
		s := 0
		switch {
		case tst.Size != nil:
			s = ti.parent.mustEvalInt(tst.ID, "size", tst.Size)
		case tst.SizeEOS:
			// TODO: error if not byte aligned?
			s = int(d.BitsLeft() / 8)

			ti.Printf("  SIZE-EOS: %#+v\n", s)
		}

		return d.FieldUTF8(tst.ID, s)

	case primitive.StrTerminated:
		// TODO: config
		// TODO: encoding
		return d.FieldUTF8Null(tst.ID)

	default:
		panic("unreachable")
	}
}

func (ti *typeInstance) decode(d *decode.D) any {
	tst := ti.schemaType

	ti.Printf("decode: tst.ID: %v\n", tst.ID)

	if ti.parent != nil && ti.parent.schemaType.Repeat == "" {
		if v, ok := ti.parent.fields[tst.ID]; ok {
			ti.Printf("  already exist\n")
			return v
		}
	}

	// TODO: only seq/instance
	// if: <expr>
	if tst.If != nil {
		ti.Printf("  If: %#+v\n", tst.If.Str)
		if !ti.mustEvalBool(tst.ID, "if", tst.If) {
			ti.Printf("    false -> nil\n")
			// ti.d.FieldValueAny(tst.ID, nil)
			return nil
		}
	}

	// value: ...
	if tst.Value != nil {
		v := ti.parent.mustEval(tst.ID, "value", tst.Value)
		ksExprField(d, tst.ID, v)

		ti.Printf("  value instance: v: %#+v\n", v)

		return v
	}

	startPos := d.Pos()
	sizeBits := d.BitsLeft()

	ti.Printf("start startPos: %#+v\n", startPos)
	ti.Printf("start sizeBits: %#+v\n", sizeBits)

	// TODO: only instance
	if tst.Pos != nil {
		ti.Printf("  Pos: %#+v\n", tst.Pos)
		v := ti.mustEvalInt64(tst.ID, "pos", tst.Pos)
		ti.Printf("    pos: %#+v\n", v)
		startPos += v * 8
		sizeBits = d.Len() - startPos
	}

	// TODO: refactor out to byteSize/bitsSize function?
	// TODO: only seq/instance
	if tst.Repeat == "" && tst.Size != nil {
		ti.Printf("  Size: %#+v\n", tst.Size.Str)
		v := ti.parent.mustEvalInt64(tst.ID, "size", tst.Size)
		ti.Printf("    size: %#+v\n", v)
		sizeBits = v * 8
	}

	ti.Printf("after startPos: %#+v\n", startPos)
	ti.Printf("after sizeBits: %#+v\n", sizeBits)

	var decodeFn func(d *decode.D) any

	if tst.Repeat != "" {
		// repeat: until|expr|eos
		// repeat-expr: <expr> # number of loops
		// repeat-until: <expr> # until false
		ti.Printf("  repeat\n")
		decodeFn = func(d *decode.D) any {
			var v any
			d.FieldArray(tst.ID, func(d *decode.D) {
				v = ti.decodeRepeat(d)
			})
			return v
		}
	} else if tst.Seq != nil || tst.Instances != nil {
		// TODO: pass is seq or instance? pos etc
		// seq:
		//   ...
		// instances:
		//   ...
		ti.Printf("  seq/instance\n")
		decodeFn = func(d *decode.D) any {
			var v any
			d.FieldStruct(tst.ID, func(d *decode.D) {
				v = ti.decodeSeq(d)
			})
			return v
		}
	} else {

		// cases:
		// no type -> bytes
		//
		// type:
		//   switch-on: <expr>
		//   cases:
		//     <expr>: <type>
		//
		// type: u1 -> primitive
		//
		// type: a -> user defined a
		resolvedTyp := &schema.Type{ID: "bytes", Primitive: primitive.Types["bytes"]}
		if tst.Type != nil {
			var typExpr ksexpr.Node

			switch {
			case tst.Type.Value != nil:
				typExpr = tst.Type.Value.KSExpr
			case tst.Type.SwitchOn != nil:
				ti.Printf("  SWITCH-ON")

				// TODO: types int vs int64 etc, ksexpr helper?

				tv := ti.parent.mustEval(tst.ID, "switch-on", tst.Type.SwitchOn)

				ti.Printf("  tv: %#+v\n", tv)

				for _, ce := range tst.Type.CasesExprs {
					// default case
					if ce.Expr.Str == "_" {
						typExpr = ce.Value.KSExpr
						break
					}

					kv := ti.parent.mustEval(tst.ID, "case", ce.Expr)
					ti.Printf("    kv: %#+v\n", kv)
					// TODO: ignore error?
					if v, _ := ksexpr.IsEqual(tv, kv); v {
						typExpr = ce.Value.KSExpr
						break
					}
				}

				// TODO: bette way?
				// no type match and there and no size, skip it
				if typExpr == nil && tst.Size == nil {
					ti.Printf("  no case match and no size, skipping")
					return nil
				}

				// TODO: not found should be bytes?
			default:
				panic("unreachable")
			}

			if typExpr != nil {
				typExprV, err := typExpr.Eval(ksexprType{ti, tst.Parent})
				if err != nil {
					ti.Fatalf("failed to resolved type: %s", err)
				}

				var ok bool
				resolvedTyp, ok = typExprV.(*schema.Type)
				if !ok {
					ti.Fatalf("type resolved to non-type: %#+v", typExprV)
				}
			}
		}

		decodeFn = func(d *decode.D) any {
			ti.Printf("decodeType: typ=%#+v\n", resolvedTyp)

			if resolvedTyp.Primitive != nil {
				v := ti.decodePrimitiveType(d, tst, *resolvedTyp.Primitive)
				ti.parent.fields[tst.ID] = v
				return v
			} else {
				ti.Printf("  user defined\n")

				// TODO: hmm
				rt := *resolvedTyp
				rt.ID = tst.ID

				tti := &typeInstance{
					schemaType: &rt,
					parent:     ti,
					root:       ti.root,
					d:          d,

					fields: map[string]any{},
				}

				return tti.decode(d)
			}
		}
	}

	if tst.Size != nil {
		var r any
		d.RangeFn(startPos, sizeBits, func(d *decode.D) {
			r = decodeFn(d)
		})
		d.SeekRel(sizeBits)
		return r
	} else {
		if tst.Pos != nil {
			d.SeekAbs(startPos)
		}
		return decodeFn(d)
	}
}

func (ti *typeInstance) decodeSeq(d *decode.D) any {
	tst := ti.schemaType

	if ti.parent != nil {
		ti.parent.fields[tst.ID] = ti
	}

	// TODO: move

	if tst.Meta != nil {
		// ti.d.FieldStruct("_meta", func(d *decode.D) {
		// 	addStrNonEmpty(d, "id", tst.Meta.ID)
		// 	addStrNonEmpty(d, "title", tst.Meta.Title)
		// 	addStrNonEmpty(d, "endian", tst.Meta.Endian)
		// })

		if tst.Meta.Endian != nil {
			// TODO: switch-on
			d.Endian = decodeEndian(d.Endian, primitive.Endianess(*tst.Meta.Endian))
		}
		if tst.Meta.BitEndian != nil {
			// TODO: switch-on
			d.BitEndian = decodeEndian(d.Endian, primitive.Endianess(*tst.Meta.BitEndian))
		}
	}

	ti.Printf("decodeSeq\n")
	for _, t := range tst.Seq {
		ti.Printf("  SEQ t.ID: %#+v\n", t.ID)

		tti := &typeInstance{
			schemaType: t,
			parent:     ti,
			root:       ti.root,
			d:          d,

			fields: map[string]any{},
		}
		tti.decode(d)
		// TODO: already exist? lint check?
		// ti.fields[t.ID] = v
		// ti.last = v
	}

	ti.Printf("decodeSeq Instances")
	for id, t := range tst.Instances {
		ti.Printf("  id: %#+v\n", id)

		tti := &typeInstance{
			schemaType: t,
			parent:     ti,
			root:       ti.root,
			d:          d,

			fields: map[string]any{},
		}
		tti.decode(d)
		// TODO: already exist? lint check?
		// ti.fields[t.ID] = v
		// ti.last = v

		// if _, err := ti.resolveInstance(id); err != nil {
		// 	ti.Fatalf("%s", err)
		// }
	}

	return ti
}

func (ti *typeInstance) decodeRepeat(d *decode.D) any {
	tst := ti.schemaType
	repeatTst := *ti.schemaType
	// TODO: better way and is correct?
	repeatTst.Repeat = ""
	repeatTst.Pos = nil

	var vs []any

	ti.parent.fields[repeatTst.ID] = vs

	switch tst.Repeat {
	case "eos":
		ti.Printf("  REPEAT-EOS:\n")

		for !d.End() {
			tti := &typeInstance{
				schemaType: &repeatTst,
				parent:     ti,
				root:       ti.root,
				d:          d,

				fields: map[string]any{},
			}

			v := tti.decode(d)
			ti.last = v
			vs = append(vs, v)

			ti.parent.fields[repeatTst.ID] = vs

			ti.Printf("REPEAT-EOS !ti.d.End(): %#+v pos=%d left=%d\n", !ti.d.End(), ti.d.Pos(), ti.d.BitsLeft())

		}
	case "until":
		if tst.RepeatUntil == nil {
			ti.Fatalf("%s: repeat: %s: without repeat-until", tst.ID, tst.Repeat)
		}

		ti.Printf("  REPEAT-UTIL: n: %#+v\n", tst.RepeatUntil.Str)

		for {
			tti := &typeInstance{
				schemaType: &repeatTst,
				parent:     ti,
				root:       ti.root,
				d:          d,

				fields: map[string]any{},
			}

			v := tti.decode(d)
			ti.last = v
			vs = append(vs, v)

			ti.parent.fields[repeatTst.ID] = vs

			// TODO: skip .parent, no new instance for repeat?
			if ti.parent.mustEvalBool(tst.ID, "repeat-until", tst.RepeatUntil) {
				break
			}

			// ti.repeat = append(ti.repeat, v)
		}
	case "expr":
		if tst.RepeatExpr == nil {
			ti.Fatalf("%s: repeat: %s: without repeat-expr", tst.ID, tst.Repeat)
		}

		// TODO: skip .parent, no new instance for repeat?
		n := ti.parent.mustEvalInt(tst.ID, "repeat-expr", tst.RepeatExpr)
		ti.Printf("  REPEAT-EXPR: n: %#+v\n", n)

		for i := 0; i < n; i++ {
			tti := &typeInstance{
				schemaType: &repeatTst,
				parent:     ti,
				root:       ti.root,
				d:          d,

				fields: map[string]any{},
			}
			v := tti.decode(d)
			ti.last = v
			vs = append(vs, v)

			ti.parent.fields[repeatTst.ID] = vs
		}
	default:
		// TODO: add verify in parser
		panic("unreachable")
	}

	return vs
}

func dumpKSTree(v any) any {
	switch v := v.(type) {
	case *typeInstance:
		m := map[string]any{}
		for k, v := range v.fields {
			m[k] = dumpKSTree(v)
		}
		return m
	case []any:
		var a []any
		for _, e := range v {
			a = append(a, dumpKSTree(e))
		}
		return a
	default:
		return v
	}
}

func kaitaiDecode(d *decode.D) any {
	var ki format.Kaitai_In
	var pai format.Probe_Args_In
	if !d.ArgAs(&ki) {
		d.Fatalf("no source option")
	}
	var r io.Reader
	r = strings.NewReader(ki.Source)
	if d.ArgAs(&pai) {
		// TODO: decode should be group aware? only one so fail probe_args return value
		f, err := d.Options.FS.Open(pai.DecodeGroup)
		if err != nil {
			d.Fatalf("fail to read source ksy: %s", err)
		}
		defer f.Close()
		r = f
	}

	t, err := schema.Parse(r)
	if err != nil {
		d.Fatalf("source: %v", err)
	}

	ti := &typeInstance{
		schemaType: t,
		parent:     nil,
		d:          d,

		fields: map[string]any{},
	}
	ti.root = ti
	// uses decodeSeq directly as d here is already a struct
	v := ti.decodeSeq(d)

	je := json.NewEncoder(os.Stdout)
	je.SetIndent("", "  ")
	_ = je.Encode(dumpKSTree(v))

	return nil
}
