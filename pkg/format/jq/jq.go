package jq

import (
	"fmt"
	"fq/pkg/bitio"
	"fq/pkg/decode"
	"fq/pkg/format"
	"io/ioutil"
	"math/big"
	"reflect"

	"github.com/itchyny/gojq"
)

func init() {
	format.MustRegister(&decode.Format{
		Name:     "jq",
		DecodeFn: jqDecode,
	})
}

func jqDecode(d *decode.D) interface{} {

	s, err := ioutil.ReadFile("/Users/wader/src/fq/pkg/format/jq/test.jq")
	if err != nil {
		panic(err)
	}
	query, err := gojq.Parse(string(s))
	if err != nil {
		panic(err)
	}

	dFuncs := map[string]gojq.Function{}

	dType := reflect.TypeOf(d)
	for i := 0; i < dType.NumMethod(); i++ {
		method := dType.Method(i)

		dFuncs[method.Name] = gojq.Function{
			Argcount: 1 << (method.Type.NumIn() - 1),
			Callback: func(c interface{}, a []interface{}) interface{} {
				method := method

				d, ok := c.(*decode.D)
				if !ok {
					return fmt.Errorf("expected decoder got %v", d)
				}

				in := make([]reflect.Value, method.Type.NumIn())
				in[0] = reflect.ValueOf(d)

				for i := 1; i < method.Type.NumIn(); i++ {
					//t := method.Type.In(i)
					object := a[i-1]
					// fmt.Println(i, "->", object)
					in[i] = reflect.ValueOf(object)
				}

				rv := method.Func.Call(in)[0].Interface()
				switch vv := rv.(type) {
				case int, bool, float64, string, nil, *decode.D:
					return rv
				case int64:
					return big.NewInt(vv)
				case uint64:
					return big.NewInt(int64(vv))
				case []byte:
					return string(vv)
				case *bitio.Buffer:
					// TODO:
					panic("bitbuf")
				default:
					panic("unreachable")
				}
			},
		}
	}

	code, err := gojq.Compile(query, gojq.WithExtraFunctions(dFuncs))

	if err != nil {
		panic(err)
	}

	iter := code.Run(d)
	for {
		v, ok := iter.Next()
		if !ok {
			break
		}
		if err, ok := v.(error); ok {
			d.Invalid(err.Error())
		}
	}

	return nil
}
