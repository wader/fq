package jq

import (
	"fq/pkg/decode"
	"fq/pkg/format"
)

func init() {
	format.MustRegister(&decode.Format{
		Name:     "jq",
		DecodeFn: jqDecode,
	})
}

func jqDecode(d *decode.D) interface{} {
	// script, ok := d.Options["script"].(string)
	// if !ok {
	// 	d.Invalid("no script in options")
	// }

	/*
		query, err := gojq.Parse(script)
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

						var ar reflect.Value
						switch aa := a[i-1].(type) {
						case *big.Int:
							ar = reflect.ValueOf(aa.Int64())
						default:
							ar = reflect.ValueOf(object)
						}

						in[i] = ar
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
	*/

	return nil
}
