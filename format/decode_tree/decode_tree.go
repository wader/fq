package decode_tree

import (
	"github.com/wader/fq/format"
	"github.com/wader/fq/internal/gojqx"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
	"github.com/wader/fq/pkg/scalar"
	"github.com/wader/gojq"
)

func init() {
	interp.RegisterFormat(
		format.Decode_Tree,
		&decode.Format{
			Description:  "Decode tree",
			DecodeFn:     decodeTree,
			DefaultInArg: format.Decode_Tree_In{},
		})
}

func decodeValue(d *decode.D, name string, v map[string]any) {
	// log.Printf("dedeoValue: %s: %#+v\n", name, v)

	oa, oOk := v["o"]
	o, _ := gojqx.Cast[int](oa)
	la, lOk := v["l"]
	l, _ := gojqx.Cast[int](la)
	va, vOk := v["v"]
	sa, _ := v["s"]

	if !(oOk && lOk && vOk) {
		d.Fatalf("%s: o, l and v must be defined", name)
	}

	s := gojqx.Normalize(sa)

	o *= 8
	l *= 8

	// TODO: o/l json int is a bit limited

	switch v := va.(type) {
	case int:
		d.SeekAbs(int64(o))
		d.FieldSintScalarFn(name, func(d *decode.D) scalar.Sint {
			d.SeekRel(int64(l))
			return scalar.Sint{
				Actual: int64(v),
				Sym:    s,
			}
		})
	case string:
		d.SeekAbs(int64(o))
		d.FieldStrScalarFn(name, func(d *decode.D) scalar.Str {
			d.SeekRel(int64(l))
			return scalar.Str{
				Actual: v,
				Sym:    s,
			}
		})
	case nil:
		d.SeekAbs(int64(o))
		d.FieldScalarRawLen(name, int64(l)) // TODO: descirption?
	case map[string]any:
		// log.Println("BLA")
		d.FieldStruct(name, func(d *decode.D) {
			decodeObject(d, v)
		})
	case []any:
		d.FieldArray(name, func(d *decode.D) {
			decodeArray(d, v)
		})
	default:
		// log.Printf("v: %#+v %T\n", v, v)
	}
}

func decodeObject(d *decode.D, o map[string]any) {
	// log.Printf("dedeoObject: %#+v\n", o)

	for k, v := range o {
		// log.Printf("k: %#+v\n", k)

		vm, ok := v.(map[string]any)
		if !ok {
			continue
		}

		decodeValue(d, k, vm)
	}

}

func decodeArray(d *decode.D, a []any) {
	// log.Printf("decodeArray: %#+v\n", a)

	for _, v := range a {
		// log.Printf("i: %#+v\n", i)

		vm, ok := v.(map[string]any)
		if !ok {
			continue
		}

		decodeValue(d, "entry" /* TOOD */, vm)
	}

}

func decodeTree(d *decode.D) any {
	var dti format.Decode_Tree_In
	d.ArgAs(&dti)

	// TODO: cast tree

	switch v := dti.Tree.(type) {
	case map[string]any:
		decodeObject(d, v)
	default:
		d.Fatalf("root must be an object, is %s", gojq.TypeOf(v))
	}

	return nil
}
