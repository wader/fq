package safetensors

// https://huggingface.co/docs/safetensors/en/index

import (
	"fmt"
	"math"
	"sort"

	"github.com/wader/fq/format"
	"github.com/wader/fq/internal/mapstruct"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
	"github.com/wader/fq/pkg/scalar"
)

var jsonFormat decode.Group

type TensorInfo struct {
	Dtype       string `mapstruct:"dtype"`
	Shape       []int  `mapstruct:"shape"`
	DataOffsets []int  `mapstruct:"data_offsets"`
}

type SafeTensorsHeader struct {
	Tensors  map[string]TensorInfo `mapstruct:",remain"`
	Metadata map[string]any        `mapstruct:"__metadata__"`
}

func init() {
	interp.RegisterFormat(
		format.SAFETENSORS,
		&decode.Format{
			Description: "SafeTensors",
			DecodeFn:    decodeSafeTensors,
			Dependencies: []decode.Dependency{
				{Groups: []*decode.Group{format.JSON}, Out: &jsonFormat},
			},
		})
}

func parseHeader(dv *decode.Value) (*SafeTensorsHeader, error) {
	actualVal, ok := dv.V.(*scalar.Any)
	if !ok {
		return nil, fmt.Errorf("expected scalar.Any, got %T", dv.V)
	}

	headerMap, ok := actualVal.Actual.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("expected map[string]any, got %T", actualVal.Actual)
	}

	var header SafeTensorsHeader
	if err := mapstruct.ToStruct(headerMap, &header); err != nil {
		return nil, fmt.Errorf("failed to parse header: %w", err)
	}

	return &header, nil
}

// https://en.wikipedia.org/wiki/Bfloat16_floating-point_format
// https://en.wikipedia.org/wiki/Single-precision_floating-point_format
// float32:  1 sign bit, 8 exponent bits, 23 fraction bits
// bfloat16: 1 sign bit, 8 exponent bits, 7 fraction bits
// To convert bfloat16 to float32, we can shift the bits to the left by 16.
func bfloat16_bits_to_float(bits uint16) float32 {
	return math.Float32frombits(uint32(bits) << 16)
}

var dataDecoders = map[string]func(d *decode.D){
	"F64": func(d *decode.D) { d.FieldF64("x") },
	"F32": func(d *decode.D) { d.FieldF32("x") },
	"F16": func(d *decode.D) { d.FieldF16("x") },
	"BF16": func(d *decode.D) {
		d.FieldFltFn("x", func(d *decode.D) float64 {
			return float64(bfloat16_bits_to_float(uint16(d.U16())))
		})
	},
	"I64":  func(d *decode.D) { d.FieldS64("x") },
	"I32":  func(d *decode.D) { d.FieldS32("x") },
	"I16":  func(d *decode.D) { d.FieldS16("x") },
	"I8":   func(d *decode.D) { d.FieldS8("x") },
	"U8":   func(d *decode.D) { d.FieldU8("x") },
	"BOOL": func(d *decode.D) { d.FieldBool("x") },
}

func decodeSafeTensors(d *decode.D) any {
	d.Endian = decode.LittleEndian

	headerSize := d.FieldU64("header size")

	var dv *decode.Value

	d.LimitedFn(8*int64(headerSize), func(d *decode.D) {
		dv, _ = d.FieldFormat("header", &jsonFormat, nil)
	})

	d.FieldStruct("tensors", func(d *decode.D) {
		header, err := parseHeader(dv)
		if err != nil {
			d.Fatalf("failed to parse header: %v", err)
			return
		}

		// Get tensor names and sort them for deterministic output
		tensorNames := make([]string, 0, len(header.Tensors))
		for tensorName := range header.Tensors {
			tensorNames = append(tensorNames, tensorName)
		}
		sort.Strings(tensorNames)

		for _, tensorName := range tensorNames {
			tensorInfo := header.Tensors[tensorName]

			decoder, exists := dataDecoders[tensorInfo.Dtype]
			if !exists {
				d.Fatalf("unsupported dtype: %s", tensorInfo.Dtype)
				continue
			}

			if len(tensorInfo.DataOffsets) < 2 {
				d.Fatalf("invalid data_offsets for tensor %s: %v", tensorName, tensorInfo.DataOffsets)
				continue
			}

			begin := tensorInfo.DataOffsets[0]

			d.FieldStruct(tensorName, func(d *decode.D) {
				d.FieldArray("shape", func(d *decode.D) {
					for _, s := range tensorInfo.Shape {
						d.FieldValueSint("dim", int64(s))
					}
				})

				if len(tensorInfo.Shape) == 0 {
					return
				}

				d.SeekAbs(8*(8+int64(headerSize)+int64(begin)), func(d *decode.D) {
					var reshape func(d *decode.D, i int)
					reshape = func(d *decode.D, i int) {
						d.FieldArray("data", func(d *decode.D) {
							if i == len(tensorInfo.Shape)-1 {
								for range tensorInfo.Shape[i] {
									decoder(d)
								}
							} else {
								for range tensorInfo.Shape[i] {
									reshape(d, i+1)
								}
							}
						})
					}
					reshape(d, 0)
				})

			})
		}
	})

	return nil
}
