package tzif

import (
	"embed"
	"time"

	"github.com/wader/fq/format"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
	"github.com/wader/fq/pkg/scalar"
)

//go:embed tzif.md
var tzifFS embed.FS

func init() {
	interp.RegisterFormat(
		format.Tzif,
		&decode.Format{
			Description: "Time Zone Information Format",
			DecodeFn:    decodeTZIF,
			Groups:      []*decode.Group{format.Probe},
		})
	interp.RegisterFS(tzifFS)
}

func decodeTZIF(d *decode.D) any {
	d.Endian = decode.BigEndian

	v1h := decodeTZifHeader(d, "v1header")
	decodeTZifDataBlock(d, v1h, 1, "v1datablock")

	if v1h.ver >= 2 {
		v2h := decodeTZifHeader(d, "v2plusheader")
		decodeTZifDataBlock(d, v2h, 2, "v2plusdatablock")
		decodeTZifFooter(d)
	}

	return nil
}

type tzifHeader struct {
	magic    uint32
	ver      uint8
	isutcnt  uint32
	isstdcnt uint32
	leapcnt  uint32
	timecnt  uint32
	typecnt  uint32
	charcnt  uint32
}

var versionToSymMapper = scalar.UintMapSymStr{
	0x00: "1",
	0x32: "2",
	0x33: "3",
}

func decodeTZifHeader(d *decode.D, name string) tzifHeader {
	var h tzifHeader

	d.FieldStruct(name, func(d *decode.D) {
		h.magic = uint32(d.FieldU32("magic", scalar.UintHex, d.UintAssert(0x545a6966)))
		h.ver = uint8(d.FieldU8("ver", d.UintAssert(0x00, 0x32, 0x33), scalar.UintHex, versionToSymMapper))
		d.FieldRawLen("reserved", 15*8)
		h.isutcnt = uint32(d.FieldU32("isutcnt"))
		h.isstdcnt = uint32(d.FieldU32("isstdcnt"))
		h.leapcnt = uint32(d.FieldU32("leapcnt"))
		h.timecnt = uint32(d.FieldU32("timecnt"))
		h.typecnt = uint32(d.FieldU32("typecnt"))
		h.charcnt = uint32(d.FieldU32("charcnt"))
	})

	if h.isutcnt != 0 && h.isutcnt != h.typecnt {
		d.Fatalf("invalid isutcnt")
	}
	if h.isstdcnt != 0 && h.isstdcnt != h.typecnt {
		d.Fatalf("invalid isstdcnt")
	}
	if h.typecnt == 0 {
		d.Fatalf("invalid typecnt")
	}
	if h.charcnt == 0 {
		d.Fatalf("invalid charcnt")
	}

	return h
}

var unixTimeToStr = scalar.SintFn(func(s scalar.Sint) (scalar.Sint, error) {
	s.Sym = time.Unix(s.Actual, 0).UTC().Format(time.RFC3339)
	return s, nil
})

func decodeTZifDataBlock(d *decode.D, h tzifHeader, decodeAsVer int, name string) {
	timeSize := 8 * 8
	if decodeAsVer == 1 {
		timeSize = 4 * 8
	}

	d.FieldStruct(name, func(d *decode.D) {

		d.FieldArray("transition_times", func(d *decode.D) {
			for i := uint32(0); i < h.timecnt; i++ {
				t := d.FieldS("transition_time", timeSize, unixTimeToStr)
				if t < -576460752303423488 {
					d.Fatalf("transition time value should be at least -2^59 (-576460752303423488), but: %d (%0#16x)", t, t)
				}
			}
		})

		d.FieldArray("transition_types", func(d *decode.D) {
			for i := uint32(0); i < h.timecnt; i++ {
				t := uint8(d.FieldU8("transition_type"))
				if uint32(t) >= h.typecnt {
					d.Fatalf("transition type must be in the range [0, %d]", h.typecnt-1)
				}
			}
		})

		d.FieldArray("local_time_type_records", func(d *decode.D) {
			for i := uint32(0); i < h.typecnt; i++ {
				d.FieldStruct("local_time_type", func(d *decode.D) {
					d.FieldS32("utoff", d.SintAssertRange(-89999, 93599))
					d.FieldU8("dst", d.UintAssert(0, 1))
					d.FieldU8("idx", d.UintAssertRange(0, uint64(h.charcnt)-1))
				})
			}
		})

		d.FieldArray("time_zone_designations", func(d *decode.D) {
			i := int(h.charcnt)
			for {
				s := d.FieldUTF8Null("time_zone_designation")
				i -= len(s) + 1
				if i <= 0 {
					break
				}
			}
		})

		d.FieldArray("leap_second_records", func(d *decode.D) {
			prevOccur := int64(0)
			prevCorr := int64(0)

			for i := uint32(0); i < h.leapcnt; i++ {
				d.FieldStruct("leap_second_record", func(d *decode.D) {
					occur := d.FieldS("occur", timeSize, unixTimeToStr)
					corr := d.FieldS32("corr")

					if i == 0 && occur < 0 {
						d.Fatalf("the first value of occur must be nonnegative")
					}
					if i > 0 && occur-prevOccur < 2419199 {
						d.Fatalf("occur must be at least 2419199 greater than the previous value")
					}
					if i == 0 && corr != 1 && corr != -1 {
						d.Fatalf("the first value of corr must be either 1 or -1")
					}
					diff := corr - prevCorr
					if i > 0 && diff != 1 && diff != -1 {
						d.Fatalf("corr must differ by exactly 1 from the previous value: diff = %d, current corr = %d, previous corr = %d", diff, corr, prevCorr)
					}

					prevOccur = occur
					prevCorr = corr
				})
			}
		})

		d.FieldArray("standard_wall_indicators", func(d *decode.D) {
			for i := uint32(0); i < h.isstdcnt; i++ {
				d.FieldU8("standard_wall_indicator", d.UintAssert(0, 1))
			}
		})

		d.FieldArray("ut_local_indicators", func(d *decode.D) {
			for i := uint32(0); i < h.isutcnt; i++ {
				d.FieldU8("ut_local_indicator", d.UintAssert(0, 1))
			}
		})
	})
}

func decodeTZifFooter(d *decode.D) {
	d.FieldStruct("footer", func(d *decode.D) {
		d.FieldU8("nl1", d.UintAssert(0x0a))
		n := d.PeekFindByte(0x0a, d.BitsLeft()/8)
		d.FieldScalarUTF8("tz_string", int(n))
		d.FieldU8("nl2", d.UintAssert(0x0a))
	})
}
