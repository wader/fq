package bencode

// https://wiki.theory.org/BitTorrentSpecification#Bencoding

import (
	"embed"
	"strconv"

	"github.com/wader/fq/format"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
	"github.com/wader/fq/pkg/scalar"
)

//go:embed bencode.jq
//go:embed bencode.md
var bencodeFS embed.FS

func init() {
	interp.RegisterFormat(
		format.Bencode,
		&decode.Format{
			Description: "BitTorrent bencoding",
			DecodeFn:    decodeBencode,
			Functions:   []string{"torepr"},
		})
	interp.RegisterFS(bencodeFS)
}

var typeToNames = scalar.StrMapSymStr{
	"d": "dictionary",
	"i": "integer",
	"l": "list",
	"0": "string",
	"1": "string",
	"2": "string",
	"3": "string",
	"4": "string",
	"5": "string",
	"6": "string",
	"7": "string",
	"8": "string",
	"9": "string",
}

func decodeStrIntUntil(b byte) func(d *decode.D) int64 {
	return func(d *decode.D) int64 {
		// 21 is sign + longest 64 bit in base 10
		i := d.PeekFindByte(b, 21)
		if i == -1 {
			d.Fatalf("decodeStrIntUntil: failed to find %v", b)
		}
		s := d.UTF8(int(i))
		n, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			d.Fatalf("decodeStrIntUntil: %q: %s", s, err)
		}
		return n
	}
}

func decodeBencodeValue(d *decode.D) {
	typ := d.FieldUTF8("type", 1, typeToNames)
	switch typ {
	case "0", "1", "2", "3", "4", "5", "6", "7", "8", "9":
		d.SeekRel(-8)
		length := d.FieldSintFn("length", decodeStrIntUntil(':'))
		d.FieldUTF8("separator", 1, d.StrAssert(":"))
		d.FieldUTF8("value", int(length))
	case "i":
		d.FieldSintFn("value", decodeStrIntUntil('e'))
		d.FieldUTF8("end", 1, d.StrAssert("e"))
	case "l":
		d.FieldArray("values", func(d *decode.D) {
			for d.PeekUintBits(8) != 'e' {
				d.FieldStruct("value", decodeBencodeValue)
			}
		})
		d.FieldUTF8("end", 1, d.StrAssert("e"))
	case "d":
		d.FieldArray("pairs", func(d *decode.D) {
			for d.PeekUintBits(8) != 'e' {
				d.FieldStruct("pair", func(d *decode.D) {
					d.FieldStruct("key", decodeBencodeValue)
					d.FieldStruct("value", decodeBencodeValue)
				})
			}
		})
		d.FieldUTF8("end", 1, d.StrAssert("e"))
	default:
		d.Fatalf("unknown type %v", typ)
	}
}

func decodeBencode(d *decode.D) any {
	decodeBencodeValue(d)
	return nil
}
