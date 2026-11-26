package pdf

// https://www.adobe.com/content/dam/acom/en/devnet/pdf/pdfs/PDF32000_2008.pdf
// https://feliam.wordpress.com/2010/08/14/pdf-a-broken-spec/
// https://github.com/modesty/pdf2json

// TODO: parse-from-end if possible?
// TODO: more unescape?
// TODO: streams filters
// TODO: refs
// TODO: EOL between object number generation "obj"?
// Ex:
// 202 0
// obj
// endobj

import (
	"bytes"
	"encoding/hex"
	"errors"
	"io"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/wader/fq/format"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
	"github.com/wader/fq/pkg/scalar"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/unicode"
)

var imageFormat decode.Group

func init() {
	interp.RegisterFormat(
		format.PDF,
		&decode.Format{
			Description: "Portable Document Format",
			Groups:      []*decode.Group{format.Probe},
			DecodeFn:    pdfDecode,
			Dependencies: []decode.Dependency{
				{Groups: []*decode.Group{format.Image}, Out: &imageFormat},
			},
		})
}

var eofLine = `%%EOF`
var whitespaceChars = []byte(" \n\r")
var delimiterChars = []byte(" \n\r()<>[]{}/%")

func isWhitespace(b byte) bool { return bytes.IndexByte(whitespaceChars, b) != -1 }

// \r, \n or \r\n
// TODO: figure someting nicer
func findLineEnd(d *decode.D) int64 {
	p := d.Pos()
	ei, v, err := d.TryPeekFind(8, 8, -1, func(v uint64) bool {
		return v == '\n' || v == '\r'
	})
	if errors.Is(err, io.EOF) {
		// to support %%EOF not having a line ending
		eofLineLen := int64(len(eofLine))
		if d.BitsLeft()*8 >= eofLineLen {
			return eofLineLen
		}
		return -1
	}
	if ei == -1 {
		return ei
	}
	ei += 8
	d.SeekRel(ei)
	if v == '\r' && d.PeekUintBits(8) == '\n' {
		ei += 8
	}
	d.SeekAbs(p)
	return ei / 8
}

var strToNumber = scalar.StrFn(func(s scalar.Str) (scalar.Str, error) {
	if strings.Contains(s.Actual, ".") {
		// real
		n, err := strconv.ParseFloat(s.Actual, 64)
		if err != nil {
			return s, err
		}
		s.Sym = n
	} else {
		// integer
		n, err := strconv.ParseInt(s.Actual, 10, 64)
		if err != nil {
			return s, err
		}
		s.Sym = n
	}

	return s, nil
})

func decodeLineStr(d *decode.D) string {
	ei := findLineEnd(d)
	if ei == -1 {
		d.Errorf("could not find line ending")
	}
	return strings.TrimRight(d.UTF8(int(ei)), "\r\n")
}

func decodeStrWhitespace(d *decode.D) string {
	ei, _ := d.PeekFind(8, 8, -1, func(v uint64) bool {
		return !isWhitespace(byte(v))
	})
	if ei == -1 {
		d.Errorf("could not find line ending")
	}
	return d.UTF8(int(ei / 8))
}

func decodeStrUntil(b []byte) func(d *decode.D) string {
	return func(d *decode.D) string {
		ei, _ := d.PeekFind(8, 8, -1, func(v uint64) bool {
			return bytes.IndexByte(b, byte(v)) != -1
		})
		if ei == -1 {
			d.Errorf("could not find ending %d", b)
		}
		return d.UTF8(int(ei / 8))
	}
}

var utf16BEExpect = unicode.UTF16(unicode.BigEndian, unicode.ExpectBOM)
var pdfDocEncoding = charmap.ISO8859_1 // TODO: correc?

// TODO: unescape \N \NN \NNN?
func decodeText(bs []byte) string {
	if s, err := utf16BEExpect.NewDecoder().String(string(bs)); err == nil {
		return s
	}
	if s, err := pdfDocEncoding.NewDecoder().String(string(bs)); err == nil {
		return s
	}
	return string(bs)
}

// "\" + "\n" or "\r" or "\r\n"
var eolEscapeRe = regexp.MustCompile("\\\\(?:\n|\r|\r\n)")

func decodeStrLiteralString(d *decode.D) string {
	ei, _ := d.PeekFind(16, 8, -1, func(v uint64) bool {
		c1 := v >> 8
		c0 := v & 0xff
		// find non-escaped ")"
		return c0 != '\\' && c1 == ')'
	})
	if ei == -1 {
		d.Errorf("could not find literal string ending")
	}

	bs := d.BytesLen(int(ei / 8))
	bs = eolEscapeRe.ReplaceAll(bs, []byte{})

	return decodeText(bs)
}

var whitespaceRE = regexp.MustCompile(`\s`)

func decodeStrHexString(d *decode.D) string {
	ei, _ := d.PeekFind(8, 8, -1, func(v uint64) bool {
		return v == '>'
	})
	if ei == -1 {
		d.Errorf("could not find literal string ending")
	}
	s := d.UTF8(int(ei / 8))
	// whitespace should be ignored inside hex string
	hexStr := whitespaceRE.ReplaceAllString(s, "")
	buf, err := hex.DecodeString(hexStr)
	if err != nil {
		d.IOPanic(err, "hex.DecodeString")
	}
	return decodeText(buf)
}

func decodeValue(d *decode.D) any {
	var r any

	if isWhitespace(byte(d.PeekUintBits(8))) {
		d.FieldStrFn("heading_whitespace", decodeStrWhitespace)
	}

	// currently only returns value for types whose values we are interested in, atm
	// only dictionary, name, string and number so we can get stream length and filter

	s := string(d.PeekBytes(2))
	switch {
	case (s[0] >= '0' && s[0] <= '9') || s[0] == '+' || s[0] == '-':
		num := d.FieldStrFn("value", decodeStrUntil(delimiterChars), strToNumber)
		r = num
	case s == "<<":
		d.FieldUTF8("start", 2)
		d.FieldArray("pairs", func(d *decode.D) {
			pairs := map[string]any{}
			for string(d.PeekBytes(2)) != ">>" {
				d.FieldStruct("pair", func(d *decode.D) {
					var keyV any
					var valueV any
					d.FieldStruct("key", func(d *decode.D) { keyV = decodeValue(d) })
					d.FieldStruct("value", func(d *decode.D) { valueV = decodeValue(d) })

					if s, ok := keyV.(string); ok {
						pairs[s] = valueV
					}
				})
			}
			r = pairs
		})
		d.FieldUTF8("end", 2)
	case s[0] == '[':
		d.FieldUTF8("start", 1)
		d.FieldArray("objects", func(d *decode.D) {
			var objectsV []any
			for d.PeekUintBits(8) != ']' {
				var valueV any
				d.FieldStruct("object", func(d *decode.D) { valueV = decodeValue(d) })
				objectsV = append(objectsV, valueV)
			}
			r = objectsV
		})
		d.FieldUTF8("end", 1)
	case s[0] == '/':
		d.FieldUTF8("start", 1)
		name := d.FieldStrFn("value", decodeStrUntil(delimiterChars))
		r = name
	case s[0] == '<':
		d.FieldUTF8("start", 1)
		s := d.FieldStrFn("value", decodeStrHexString)
		d.FieldUTF8("end", 1)
		r = s
	case s[0] == '(':
		d.FieldUTF8("start", 1)
		s := d.FieldStrFn("value", decodeStrLiteralString)
		d.FieldUTF8("end", 1)
		r = s
	case s[0] == 'R':
		// TODO: should handle references differently?
		d.FieldUTF8("reference", 1)
	case s[0] == 't':
		d.FieldUTF8("value", 4, scalar.StrSym(true))
	case s[0] == 'f':
		d.FieldUTF8("value", 5, scalar.StrSym(false))
	case s[0] == 'n':
		d.FieldUTF8("value", 4, scalar.StrSym(nil))
	default:
		d.Fatalf("unknown type %q", s)
	}

	if isWhitespace(byte(d.PeekUintBits(8))) {
		d.FieldStrFn("tailing_whitespace", decodeStrWhitespace)
	}

	return r
}

func pdfObjBodyDecode(d *decode.D) {
	d.FieldValueStr("type", "body")
	d.FieldArray("objects", func(d *decode.D) {
		d.FieldStruct("object", func(d *decode.D) {
			d.FieldStrFn("start", decodeLineStr)
			var dictV any
			d.FieldStruct("dictionary", func(d *decode.D) { dictV = decodeValue(d) })

			log.Printf("dictV: %#+v\n", dictV)

			endObj := false
			for !endObj {
				e := d.Pos()
				line := decodeLineStr(d)
				d.SeekAbs(e)

				switch line {
				case "endobj":
					d.FieldStrFn("end", decodeLineStr)
					endObj = true
				case "stream":
					d.FieldStruct("stream", func(d *decode.D) {

						d.FieldStrFn("start", decodeLineStr)

						// TODO: proper string find
						ei, _ := d.PeekFind(64, 8, -1, func(v uint64) bool {
							return v == (0 |
								'e'<<56 |
								'n'<<48 |
								'd'<<40 |
								's'<<32 |
								't'<<24 |
								'r'<<16 |
								'e'<<8 |
								'a')
						})

						// _ = d.FieldRawLen("data", ei)

						if dv, _, _ := d.TryFieldFormatLen("data", ei, &imageFormat, nil); dv == nil {
							_ = d.FieldRawLen("data", ei)
						}

						// d.FieldFormatReaderLen("uncompressed", dataLen, zlib.NewReader, iccProfileFormat)

						// var rFn func(r io.Reader) io.Reader
						// switch compressionMethod {
						// case delfateMethod:
						// *bitio.Buffer implements io.ByteReader so hat deflate don't do own
						// // buffering and might read more than needed messing up knowing compressed size
						// rFn = func(r io.Reader) io.Reader { return flate.NewReader(r) }
						// // }

						// if rFn != nil {
						// 	readCompressedSize, uncompressedBB, dv, _, err := d.TryFieldReaderRangeFormat("uncompressed", d.Pos(), ei, rFn, imageFormat, nil)
						// 	log.Printf("err: %#+v\n", err)
						// 	if uncompressedBB != nil {
						// 		if dv == nil {
						// 			d.FieldRootBitBuf("uncompressed", uncompressedBB)
						// 		}
						// 		d.FieldRawLen("compressed", readCompressedSize)
						// 	}
						// }

						d.FieldStrFn("end", decodeLineStr)
					})
				default:
					d.Fatalf("bla")
				}
			}
		})
	})
}

func pdfDecode(d *decode.D) any {
	d.FieldStruct("header", func(d *decode.D) {
		d.FieldUTF8("version", 8, d.StrAssert(
			"%PDF-1.0",
			"%PDF-1.1",
			"%PDF-1.2",
			"%PDF-1.3",
			"%PDF-1.4",
			"%PDF-1.5",
			"%PDF-1.6",
			"%PDF-1.7",
		))

		d.SeekAbs(0)

		// TODO: has binary comment, byte >= 128

		d.FieldArray("comments", func(d *decode.D) {
			for d.PeekUintBits(8) == '%' {
				d.FieldStrFn("comment", decodeLineStr)
			}
		})
	})

	d.FieldArray("parts", func(d *decode.D) {
		for !d.End() {
			p := d.Pos()

			line := decodeLineStr(d)
			d.SeekAbs(p)

			switch {
			case strings.HasPrefix(line, "%"):
				d.FieldStruct("comment", func(d *decode.D) {
					d.FieldValueStr("type", "comment")
					d.FieldStrFn("line", decodeLineStr)
				})
			case strings.TrimSpace(line) == "":
				d.FieldStruct("whitespace", func(d *decode.D) {
					d.FieldValueStr("type", "whitespace")
					d.FieldStrFn("line", decodeLineStr)
				})
			case strings.HasSuffix(line, "obj"):
				d.FieldStruct("body", pdfObjBodyDecode)
			case line == "xref":
				d.FieldStruct("xref", func(d *decode.D) {
					d.FieldValueStr("type", "xref")
					d.FieldStrFn("start", decodeLineStr)
					d.FieldArray("lines", func(d *decode.D) {
						for {
							b := byte(d.PeekUintBits(8))
							if !(b >= '0' && b <= '9') {
								break
							}
							d.FieldStrFn("start", decodeLineStr)
						}
					})
				})
			case line == "trailer":
				d.FieldStruct("trailer", func(d *decode.D) {
					d.FieldValueStr("type", "trailer")
					d.FieldStrFn("start", decodeLineStr)
					d.FieldStruct("dictionary", func(d *decode.D) { decodeValue(d) })
				})
			case line == "startxref":
				d.FieldStruct("startxref", func(d *decode.D) {
					d.FieldValueStr("type", "startxref")
					d.FieldStrFn("start", decodeLineStr)
					d.FieldStrFn("offset", decodeLineStr, strToNumber)
				})
			default:
				d.Fatalf("unknown line %q", line)
			}
		}
	})

	return nil
}
