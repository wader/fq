package jq

import (
	"fmt"
	"fq/pkg/decode"
	"fq/pkg/format"
	"io/ioutil"
	"log"
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

				d := c.(*decode.D)

				in := make([]reflect.Value, method.Type.NumIn())
				in[0] = reflect.ValueOf(d)

				for i := 1; i < method.Type.NumIn(); i++ {
					//t := method.Type.In(i)
					object := a[i-1]
					fmt.Println(i, "->", object)
					in[i] = reflect.ValueOf(object)
				}

				return method.Func.Call(in)[0].Interface()
			},
		}
	}

	//log.Printf("dFuncs: %#+v\n", dFuncs)

	code, err := gojq.Compile(query, gojq.WithExtraFunctions(dFuncs))

	// code, err := gojq.Compile(query, gojq.WithExtraFunctions(map[string]gojq.Function{
	// 	"Struct": {
	// 		Argcount: gojq.Argcount1,
	// 		Callback: func(c interface{}, a []interface{}) interface{} {
	// 			d := c.(*decode.D)
	// 			name := a[0].(string)
	// 			return d.FieldStruct(name)
	// 		},
	// 	},

	// 	"Array": {
	// 		Argcount: gojq.Argcount1,
	// 		Callback: func(c interface{}, a []interface{}) interface{} {
	// 			d := c.(*decode.D)
	// 			name := a[0].(string)
	// 			return d.FieldArray(name)
	// 		},
	// 	},

	// 	"FieldU": {
	// 		Argcount: gojq.Argcount2,
	// 		Callback: func(c interface{}, a []interface{}) interface{} {
	// 			log.Printf("d: %#+v\n", d)
	// 			d := c.(*decode.D)
	// 			name := a[0].(string)
	// 			bits := a[1].(int)

	// 			return (&big.Int{}).SetUint64(d.FieldU(name, bits))
	// 		},
	// 	},

	// 	"U": {
	// 		Argcount: gojq.Argcount1,
	// 		Callback: func(c interface{}, a []interface{}) interface{} {
	// 			d := c.(*decode.D)
	// 			bits := a[0].(int)

	// 			return (&big.Int{}).SetUint64(d.U(bits))
	// 		},
	// 	},
	// }))
	if err != nil {
		panic(err)
	}

	iter := code.Run(d)
	for {
		v, ok := iter.Next()
		if !ok {
			break
		}
		log.Printf("v: %#+v\n", v)
	}

	return nil
}
