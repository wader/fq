package http

import (
	"fmt"
	"log"
	"regexp"

	"github.com/wader/fq/format"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
)

var mimeMultipartTextprotoGroup decode.Group

func init() {
	interp.RegisterFormat(
		format.MIME_Multi_Part,
		&decode.Format{
			Description:  "MIME multipart",
			Groups:       []*decode.Group{format.Content_Type},
			DecodeFn:     mimeMultipartDecode,
			DefaultInArg: format.Mime_Multipart_In{},
			Dependencies: []decode.Dependency{
				{Groups: []*decode.Group{format.TextProto}, Out: &mimeMultipartTextprotoGroup},
			},
		})
}

const multipartContentType = "multipart/form-data"

func mimeMultipartDecode(d *decode.D) any {
	var boundary string

	log.Println("multipart:")

	var cti format.Content_Type_In
	var mmpi format.Mime_Multipart_In
	if d.ArgAs(&cti) {
		log.Printf("cti: %#+v\n", cti)
		if cti.ContentType != multipartContentType {
			d.Fatalf("content-type not " + multipartContentType)
		}
		boundary = cti.Pairs["boundary"]
	} else if d.ArgAs(&mmpi) {
		log.Printf("mmpi: %#+v\n", mmpi)
		boundary = mmpi.Boundary
	}

	if boundary == "" {
		d.Fatalf("no boundary set")
	}

	const boundaryREEndGroup = 1
	var boundaryRE = regexp.MustCompile(fmt.Sprintf("--%s(?P<end>--)?\r\n", regexp.QuoteMeta(boundary)))
	var endBoundaryLen int64

	firstBoundaryRs := d.RE(boundaryRE)
	if firstBoundaryRs == nil {
		d.Fatalf("first boundary %q not found", boundary)
	}
	firstBoundaryR := firstBoundaryRs[0]

	d.FieldUTF8("preamble", int(firstBoundaryR.Start/8))
	d.FieldArray("parts", func(d *decode.D) {
		for {
			boundaryStartRs := d.RE(boundaryRE)
			boundaryStartR := boundaryStartRs[0]
			boundaryStartEnd := boundaryStartRs[boundaryREEndGroup]

			if boundaryStartRs == nil {
				d.Fatalf("boundary %q not found", boundary)
			}
			if boundaryStartEnd.Start != -1 {
				// found a boundary with ending "--"
				endBoundaryLen = boundaryStartR.Len
				break
			}

			d.FieldStruct("part", func(d *decode.D) {
				d.FieldUTF8("start_boundary", int(boundaryStartR.Len/8))

				boundaryEndRs := d.RE(boundaryRE)
				if boundaryEndRs == nil {
					d.Fatalf("boundary end %q not found", boundary)
				}
				boundaryEndR := boundaryEndRs[0]

				partLen := (boundaryEndR.Start - boundaryStartR.Stop()) /* \r\n */
				d.FramedFn(partLen, func(d *decode.D) {
					d.FieldFormat("headers", &mimeMultipartTextprotoGroup, format.TextProto_In{Name: "header"})
					d.FieldUTF8("header_end", 2)
					d.FieldRawLen("data", d.BitsLeft()-16)
					d.FieldUTF8("data_end", 2)
				})

				d.SeekAbs(boundaryEndRs[0].Start)
			})
		}
	})

	d.FieldUTF8("end_boundary", int(endBoundaryLen/8))
	d.FieldUTF8("epilogue", int(d.BitsLeft()/8))

	return nil
}
