package negentropy

// https://github.com/hoytech/negentropy

import (
	"embed"
	"math"
	"time"

	"github.com/wader/fq/format"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
	"github.com/wader/fq/pkg/scalar"
)

//go:embed negentropy.md
var negFS embed.FS

func init() {
	interp.RegisterFormat(
		format.Negentropy,
		&decode.Format{
			Description: "Negentropy message",
			DecodeFn:    decodeNegentropyMessage,
			Groups:      []*decode.Group{},
		})
	interp.RegisterFS(negFS)
}

const (
	version         = 0x61
	fingerprintSize = 16
)

var modeMapper = scalar.SintMapSymStr{
	0: "skip",
	1: "fingerprint",
	2: "idlist",
}

type timestampTranslator struct{ last time.Time }

func (tt *timestampTranslator) MapSint(s scalar.Sint) (scalar.Sint, error) {
	if s.Actual == 0 {
		s.Description = "infinity"
		tt.last = time.Unix(math.MaxInt64, 0)
		return s, nil
	} else {
		timestamp := tt.last.Add(time.Second * time.Duration(s.Actual-1))
		s.Description = timestamp.UTC().Format(time.DateTime + " UTC")
		s.Sym = timestamp.Unix()
		tt.last = timestamp
		return s, nil
	}
}

func decodeNegentropyMessage(d *decode.D) any {
	tt := &timestampTranslator{last: time.Unix(0, 0)}

	d.Endian = decode.BigEndian

	v := d.FieldU8("version")
	if v != version {
		d.Fatalf("unexpected version %d (expected %d), is this really a negentropy message?", v, version)
	}

	d.FieldStructArrayLoop("bounds", "bound", d.NotEnd, func(d *decode.D) {
		d.FieldSintFn("timestamp", decodeVarInt, tt)

		size := d.FieldSintFn("id_prefix_size", decodeVarInt)
		if size > 32 {
			d.Fatalf("unexpected id prefix size bigger than 32: %d", size)
		}
		if size > 0 {
			d.FieldRawLen("id_prefix", size*8, scalar.RawHex)
		}

		mode := d.FieldSintFn("mode", decodeVarInt, modeMapper)
		switch mode {
		case 0:
			return
		case 1:
			d.FieldRawLen("fingerprint", fingerprintSize*8, scalar.RawHex)
			return
		case 2:
			d.FieldStruct("idlist", func(d *decode.D) {
				num := d.FieldSintFn("num", decodeVarInt)
				d.FieldArray("ids", func(d *decode.D) {
					for i := 0; i < int(num); i++ {
						d.FieldRawLen("id", 32*8, scalar.RawHex)
					}
				})
			})
		default:
			d.Fatalf("unexpected mode %d", mode)
		}
	})

	return nil
}

func decodeVarInt(d *decode.D) int64 {
	var res int64 = 0

	for {
		b := d.U8()
		res = (res << 7) | (int64(b) & 127)
		if (b & 128) == 0 {
			break
		}
	}

	return res
}
