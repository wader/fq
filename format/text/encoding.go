package text

import (
	"bytes"
	"embed"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"strings"

	"github.com/wader/fq/pkg/bitio"
	"github.com/wader/fq/pkg/interp"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/unicode"
)

//go:embed encoding.jq
var textFS embed.FS

func init() {
	interp.RegisterFunc0("from_hex", func(_ *interp.Interp, c string) any {
		b, err := hex.DecodeString(c)
		if err != nil {
			return err
		}
		bb, err := interp.NewBinaryFromBitReader(bitio.NewBitReader(b, -1), 8, 0)
		if err != nil {
			return err
		}
		return bb
	})
	interp.RegisterFunc0("to_hex", func(_ *interp.Interp, c any) any {
		br, err := interp.ToBitReader(c)
		if err != nil {
			return err
		}
		buf := &bytes.Buffer{}
		if _, err := io.Copy(hex.NewEncoder(buf), bitio.NewIOReader(br)); err != nil {
			return err
		}
		return buf.String()
	})

	// TODO: other encodings and share?
	base64Encoding := func(enc string) *base64.Encoding {
		switch enc {
		case "url":
			return base64.URLEncoding
		case "rawstd":
			return base64.RawStdEncoding
		case "rawurl":
			return base64.RawURLEncoding
		default:
			return base64.StdEncoding
		}
	}
	type fromBase64Opts struct {
		Encoding string
	}
	interp.RegisterFunc1("_from_base64", func(_ *interp.Interp, c string, opts fromBase64Opts) any {
		b, err := base64Encoding(opts.Encoding).DecodeString(c)
		if err != nil {
			return err
		}
		bin, err := interp.NewBinaryFromBitReader(bitio.NewBitReader(b, -1), 8, 0)
		if err != nil {
			return err
		}
		return bin
	})
	type toBase64Opts struct {
		Encoding string
	}
	interp.RegisterFunc1("_to_base64", func(_ *interp.Interp, c any, opts toBase64Opts) any {
		br, err := interp.ToBitReader(c)
		if err != nil {
			return err
		}
		bb := &bytes.Buffer{}
		wc := base64.NewEncoder(base64Encoding(opts.Encoding), bb)
		if _, err := io.Copy(wc, bitio.NewIOReader(br)); err != nil {
			return err
		}
		wc.Close()
		return bb.String()
	})

	strEncoding := func(s string) encoding.Encoding {
		switch s {
		case "UTF8":
			return unicode.UTF8
		case "UTF16":
			return unicode.UTF16(unicode.LittleEndian, unicode.UseBOM)
		case "UTF16LE":
			return unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM)
		case "UTF16BE":
			return unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM)
		case "CodePage037":
			return charmap.CodePage037
		case "CodePage437":
			return charmap.CodePage437
		case "CodePage850":
			return charmap.CodePage850
		case "CodePage852":
			return charmap.CodePage852
		case "CodePage855":
			return charmap.CodePage855
		case "CodePage858":
			return charmap.CodePage858
		case "CodePage860":
			return charmap.CodePage860
		case "CodePage862":
			return charmap.CodePage862
		case "CodePage863":
			return charmap.CodePage863
		case "CodePage865":
			return charmap.CodePage865
		case "CodePage866":
			return charmap.CodePage866
		case "CodePage1047":
			return charmap.CodePage1047
		case "CodePage1140":
			return charmap.CodePage1140
		case "ISO8859_1":
			return charmap.ISO8859_1
		case "ISO8859_2":
			return charmap.ISO8859_2
		case "ISO8859_3":
			return charmap.ISO8859_3
		case "ISO8859_4":
			return charmap.ISO8859_4
		case "ISO8859_5":
			return charmap.ISO8859_5
		case "ISO8859_6":
			return charmap.ISO8859_6
		case "ISO8859_6E":
			return charmap.ISO8859_6E
		case "ISO8859_6I":
			return charmap.ISO8859_6I
		case "ISO8859_7":
			return charmap.ISO8859_7
		case "ISO8859_8":
			return charmap.ISO8859_8
		case "ISO8859_8E":
			return charmap.ISO8859_8E
		case "ISO8859_8I":
			return charmap.ISO8859_8I
		case "ISO8859_9":
			return charmap.ISO8859_9
		case "ISO8859_10":
			return charmap.ISO8859_10
		case "ISO8859_13":
			return charmap.ISO8859_13
		case "ISO8859_14":
			return charmap.ISO8859_14
		case "ISO8859_15":
			return charmap.ISO8859_15
		case "ISO8859_16":
			return charmap.ISO8859_16
		case "KOI8R":
			return charmap.KOI8R
		case "KOI8U":
			return charmap.KOI8U
		case "Macintosh":
			return charmap.Macintosh
		case "MacintoshCyrillic":
			return charmap.MacintoshCyrillic
		case "Windows874":
			return charmap.Windows874
		case "Windows1250":
			return charmap.Windows1250
		case "Windows1251":
			return charmap.Windows1251
		case "Windows1252":
			return charmap.Windows1252
		case "Windows1253":
			return charmap.Windows1253
		case "Windows1254":
			return charmap.Windows1254
		case "Windows1255":
			return charmap.Windows1255
		case "Windows1256":
			return charmap.Windows1256
		case "Windows1257":
			return charmap.Windows1257
		case "Windows1258":
			return charmap.Windows1258
		case "XUserDefined":
			return charmap.XUserDefined
		default:
			return nil
		}
	}

	type toStrEncodingOpts struct {
		Encoding string
	}
	interp.RegisterFunc1("_to_strencoding", func(_ *interp.Interp, c string, opts toStrEncodingOpts) any {
		h := strEncoding(opts.Encoding)
		if h == nil {
			return fmt.Errorf("unknown string encoding %s", opts.Encoding)
		}

		bb := &bytes.Buffer{}
		if _, err := io.Copy(h.NewEncoder().Writer(bb), strings.NewReader(c)); err != nil {
			return err
		}
		outBR := bitio.NewBitReader(bb.Bytes(), -1)
		bin, err := interp.NewBinaryFromBitReader(outBR, 8, 0)
		if err != nil {
			return err
		}

		return bin
	})

	type fromStrEncodingOpts struct {
		Encoding string
	}
	interp.RegisterFunc1("_from_strencoding", func(_ *interp.Interp, c any, opts fromStrEncodingOpts) any {
		inBR, err := interp.ToBitReader(c)
		if err != nil {
			return err
		}
		h := strEncoding(opts.Encoding)
		if h == nil {
			return fmt.Errorf("unknown string encoding %s", opts.Encoding)
		}

		bb := &bytes.Buffer{}
		if _, err := io.Copy(bb, h.NewDecoder().Reader(bitio.NewIOReader(inBR))); err != nil {

			return err
		}

		return bb.String()
	})

	interp.RegisterFS(textFS)
}
