package ogg

// https://xiph.org/ogg/doc/framing.html

import (
	"fq/internal/bitio"
	"fq/pkg/decode"
	"fq/pkg/format"
	"log"
)

func init() {
	format.MustRegister(&decode.Format{
		Name:     format.OGG,
		Groups:   []string{format.PROBE},
		MIMEs:    []string{"audio/ogg"},
		DecodeFn: decodeOgg,
		Deps: []decode.Dep{
			{Names: []string{format.OGG_PAGE}, Formats: &oggPage},
			{Names: []string{format.OGG_PACKET}, Formats: &oggPacket},
		},
	})
}

var oggPage []*decode.Format
var oggPacket []*decode.Format

type stream struct {
	firstBit   int64
	sequenceNo uint32
	packetBuf  []byte
}

func decodeOgg(d *decode.D) interface{} {
	validPages := 0
	streams := map[uint32]*stream{}

	packets := d.FieldArray("packet")

	d.FieldArrayFn("page", func(d *decode.D) {
		for !d.End() {
			// TODO: FieldTryDecode return field and decoder? handle error?
			_, dv, errs := d.FieldTryDecode("page", oggPage)
			if errs != nil {
				break
			}
			p, _ := dv.(*page)
			if p == nil {
				// TODO: hmm
				break
			}

			s, sFound := streams[p.StreamSerialNumber]
			if !sFound {
				s = &stream{sequenceNo: p.SequenceNo}
				streams[p.StreamSerialNumber] = s
			}

			if !sFound && !p.IsFirstPage {
				// TODO: not first page and we haven't seen the stream before
				log.Println("not first page and we haven't seen the stream before")
			}
			hasData := len(s.packetBuf) > 0
			if p.IsContinuedPacket && !hasData {
				// TODO: continuation but we haven't seen any packet data yet
				log.Println("continuation but we haven't seen any packet data yet")
			}
			if !p.IsFirstPage && s.sequenceNo+1 != p.SequenceNo {
				// TODO: page gap
				log.Println("page gap")
			}

			for _, ps := range p.Segments {
				if s.packetBuf == nil {
					s.firstBit = d.Pos()
				}

				// TODO: decoder buffer api that panics?
				psLen, _ := ps.Len()

				// TODO: cleanup
				b, _ := ps.BytesRange(0, int(psLen/8))
				s.packetBuf = append(s.packetBuf, b...)
				if psLen/8 < 255 { // TODO: list range maps of demuxed packets?
					bb := bitio.NewBufferFromBytes(s.packetBuf, -1)
					packets.FieldDecodeBitBuf("packet", bb, oggPacket)
					s.packetBuf = nil
				}
			}

			s.sequenceNo = p.SequenceNo
			if p.IsLastPage {
				delete(streams, p.StreamSerialNumber)
			}

			validPages++
		}
	})

	if validPages == 0 {
		d.Invalid("no pages found")
	}

	return nil
}
