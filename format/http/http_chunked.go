package http

import (
	"log"
	"strconv"

	"github.com/wader/fq/format"
	"github.com/wader/fq/internal/lazyre"
	"github.com/wader/fq/pkg/bitio"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
	"github.com/wader/fq/pkg/scalar"
)

var httpChunkedContentTypeGroup decode.Group
var httpChunkedGzipGroup decode.Group

func init() {
	interp.RegisterFormat(
		format.HTTP_Chunked,
		&decode.Format{
			Description:  "HTTP chunked encoding",
			DecodeFn:     httpChunkedDecode,
			DefaultInArg: format.Http_Chunked_In{},
			Dependencies: []decode.Dependency{
				{Groups: []*decode.Group{format.Content_Type}, Out: &httpChunkedContentTypeGroup},
				{Groups: []*decode.Group{format.Gzip}, Out: &httpChunkedGzipGroup},
			},
		})
}

var chunkStartLineRE = &lazyre.RE{S: `(?P<length>.*\r?\n)`}

func httpChunkedDecode(d *decode.D) any {
	var hci format.Http_Chunked_In
	hciOk := d.ArgAs(&hci)

	var chunkBRs []bitio.ReadAtSeeker

	d.FieldArray("chunks", func(d *decode.D) {
		seenEnd := false
		for !seenEnd {
			d.FieldStruct("chunk", func(d *decode.D) {
				// TODO: chunk extension
				cm := map[string]string{}
				d.FieldRE(chunkStartLineRE.Must(), &cm, scalar.ActualTrimSpace)

				lengthStr := cm["length"]
				length, err := strconv.ParseInt(lengthStr, 16, 64)
				if err != nil {
					d.Fatalf("failed to parse length %q", lengthStr)
				}

				br := d.FieldRawLen("data", length*8)
				chunkBRs = append(chunkBRs, br)
				d.FieldUTF8("new_line", 2)
				if length == 0 {
					// TODO: trailer
					seenEnd = true
					return
				}
			})
		}
	})

	mbr, err := bitio.NewMultiReader(chunkBRs...)
	if err != nil {
		d.IOPanic(err, "bitio.NewMultiReader", "bitio.NewMultiReader")
	}

	if hciOk {
		log.Printf("chunked bodyGroupInArg: %#+v\n", hci)

		d.FieldStructRootBitBufFn("data", mbr, func(d *decode.D) {
			// TODO: http content encoding group?
			switch hci.ContentEncoding {
			case "gzip":
				d.FieldFormatOrRaw("body", &httpChunkedGzipGroup, format.Gzip_In{
					ContentType: hci.ContentType,
					Pairs:       hci.Pairs,
				})
			default:
				d.FieldFormatOrRaw("body", &httpChunkedContentTypeGroup, format.Content_Type_In{
					ContentType: hci.ContentType,
					Pairs:       hci.Pairs,
				})
			}
		})

	} else {
		d.FieldRootBitBuf("data", mbr)
	}

	return nil
}
