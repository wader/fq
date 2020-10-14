package ogg

// https://xiph.org/ogg/doc/framing.html

import (
	"fq/pkg/decode"
	"fq/pkg/format"
)

var oggPage []*decode.Format

var File = format.MustRegister(&decode.Format{
	Name:  "ogg",
	MIMEs: []string{"audio/ogg"},
	New: func() decode.Decoder {
		return &FileDecoder{
			streams: map[uint32]*stream{},
		}
	},
	Deps: []decode.Dep{
		{Names: []string{"ogg_page"}, Formats: &oggPage},
	},
})

type stream struct {
	firstBit   int64
	sequenceNo uint32
	packetBuf  []byte
}

// Decoder is a ogg decoder
type FileDecoder struct {
	decode.Common

	streams map[uint32]*stream
}

// Decode ogg
func (d *FileDecoder) Decode() {
	validPages := 0

	for !d.End() {
		// TODO: FieldTryDecode return field and decoder?
		_, errs := d.FieldTryDecode("page", oggPage)
		if errs != nil {
			break
		}
		// p, _ := pageDecoder.(*PageDecoder)
		// if p == nil {
		// 	// TODO: hmm
		// 	break
		// }

		// s, sFound := d.streams[p.StreamSerialNumber]
		// if !sFound {
		// 	s = &stream{sequenceNo: p.SequenceNo}
		// 	d.streams[p.StreamSerialNumber] = s
		// }

		// if !sFound && !p.IsFirstPage {
		// 	// TODO: not first page and we haven't seen the stream before
		// 	log.Println("not first page and we haven't seen the stream before")
		// }
		// hasData := len(s.packetBuf) > 0
		// if p.IsContinuedPacket && !hasData {
		// 	// TODO: continuation but we haven't seen any packet data yet
		// 	log.Println("continuation but we haven't seen any packet data yet")
		// }
		// if !p.IsFirstPage && s.sequenceNo+1 != p.SequenceNo {
		// 	// TODO: page gap
		// 	log.Println("page gap")
		// }

		// log.Printf("p.SequenceNo: %#+v\n", p.SequenceNo)

		// for _, ps := range p.Segments {
		// 	if s.packetBuf == nil {
		// 		s.firstBit = d.Pos()
		// 	}
		// 	if ps.Len/8 < 255 { // TODO: list range maps of demuxed packets?
		// 		bb, err := bitbuf.NewFromBytes(s.packetBuf, 0)
		// 		if err != nil {
		// 			panic(err) // TODO: fixme
		// 		}
		// 		d.FieldDecodeBitBuf("packet", s.firstBit, d.Pos(), bb, vorbis.Packet)
		// 		s.packetBuf = nil
		// 	}
		// }

		// s.sequenceNo = p.SequenceNo
		// if p.IsLastPage {
		// 	delete(d.streams, p.StreamSerialNumber)
		// }

		validPages++
	}

	if validPages == 0 {
		d.Invalid("no frames found")
	}
}
