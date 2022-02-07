package interp

import (
	"io"
	"regexp"
	"strings"

	"github.com/wader/fq/internal/gojqextra"
	"github.com/wader/fq/internal/ioextra"
	"github.com/wader/fq/pkg/bitio"
	"github.com/wader/fq/pkg/ranges"
	"github.com/wader/gojq"
)

func init() {
	functionRegisterFns = append(functionRegisterFns, func(i *Interp) []Function {
		return []Function{
			{"_match_buffer", 1, 2, nil, i._bufferMatch},
		}
	})
}

func (i *Interp) _bufferMatch(c interface{}, a []interface{}) gojq.Iter {
	var ok bool

	bv, err := toBuffer(c)
	if err != nil {
		return gojq.NewIter(err)
	}

	var re string
	var byteRunes bool
	var global bool

	switch a0 := a[0].(type) {
	case string:
		re = a0
	default:
		reBuf, err := toBytes(a0)
		if err != nil {
			return gojq.NewIter(err)
		}
		var reRs []rune
		for _, b := range reBuf {
			reRs = append(reRs, rune(b))
		}
		byteRunes = true
		// escape paratheses runes etc
		re = regexp.QuoteMeta(string(reRs))
	}

	var flags string
	if len(a) > 1 {
		flags, ok = a[1].(string)
		if !ok {
			return gojq.NewIter(gojqextra.FuncTypeNameError{Name: "find", Typ: "string"})
		}
	}

	if strings.Contains(flags, "b") {
		byteRunes = true
	}
	global = strings.Contains(flags, "g")

	// TODO: err to string
	// TODO: extract to regexpextra? "all" FindReaderSubmatchIndex that can iter?
	sre, err := gojqextra.CompileRegexp(re, "gimb", flags)
	if err != nil {
		return gojq.NewIter(err)
	}
	sreNames := sre.SubexpNames()

	br, err := bv.toBuffer()
	if err != nil {
		return gojq.NewIter(err)
	}

	var rr interface {
		io.RuneReader
		io.Seeker
	}
	// raw bytes regexp matching is a bit tricky, what we do is to read each byte as a codepoint (ByteRuneReader)
	// and then we can use UTF-8 encoded codepoint to match a raw byte. So for example \u00ff (encoded as 0xc3 0xbf)
	// will match the byte \0xff
	if byteRunes {
		// byte mode, read each byte as a rune
		rr = ioextra.ByteRuneReader{RS: bitio.NewIOReadSeeker(br)}
	} else {
		rr = ioextra.RuneReadSeeker{RS: bitio.NewIOReadSeeker(br)}
	}

	var off int64
	prevOff := int64(-1)
	return iterFn(func() (interface{}, bool) {
		// TODO: correct way to handle empty match for buffer, move one byte forward?
		// > "asdasd" | [match(""; "g")], [(tobytes | match(""; "g"))] | length
		// 7
		// 1
		if prevOff == off {
			return nil, false
		}

		if prevOff != -1 && !global {
			return nil, false
		}

		_, err = rr.Seek(off, io.SeekStart)
		if err != nil {
			return err, false
		}

		l := sre.FindReaderSubmatchIndex(rr)
		if l == nil {
			return nil, false
		}

		var captures []interface{}
		var firstCapture map[string]interface{}

		for i := 0; i < len(l)/2; i++ {
			start, end := l[i*2], l[i*2+1]
			capture := map[string]interface{}{
				"offset": int(off) + start,
				"length": end - start,
			}

			if start != -1 {
				matchBitOff := (off + int64(start)) * 8
				matchLength := int64(end-start) * 8
				bbo := Buffer{
					br: bv.br,
					r: ranges.Range{
						Start: bv.r.Start + matchBitOff,
						Len:   matchLength,
					},
					unit: 8,
				}

				capture["string"] = bbo
			} else {
				capture["string"] = nil
			}

			if i > 0 {
				if sreNames[i] != "" {
					capture["name"] = sreNames[i]
				} else {
					capture["name"] = nil
				}
			}

			if i == 0 {
				firstCapture = capture
			}

			captures = append(captures, capture)
		}

		prevOff = off
		off = off + int64(l[1])

		firstCapture["captures"] = captures[1:]

		return firstCapture, true
	})
}
