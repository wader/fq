package ogg

// https://xiph.org/ogg/doc/framing.html

import (
	"fq/internal/bitbuf"
	"fq/internal/decode"
	"log"
)

var Register = &decode.Register{
	Name: "ogg",
	MIME: "",
	New: func(common decode.Common) decode.Decoder {
		return &Decoder{
			Common:  common,
			streams: map[uint32]*stream{},
		}
	},
}

type stream struct {
	sequenceNo uint32
	packetBuf  []byte
}

// Decoder is a ogg decoder
type Decoder struct {
	decode.Common

	streams map[uint32]*stream
}

// Decode ogg
func (d *Decoder) Decode(opts decode.Options) {
	for !d.End() {
		d.FieldNoneFn("page", func() {
			// TODO: validate bits left
			d.FieldValidateString("capture_pattern", "OggS")
			d.FieldValidateUFn("stream_structure_version", 0, d.U8)
			d.FieldU5("unused_flags")
			isLastPage := d.FieldBool("last_page")
			isFirstPage := d.FieldBool("first_page")
			isContinuedPacket := d.FieldBool("continued_packet")
			d.FieldU64LE("absolute_granule_position")
			streamSerialNumber := uint32(d.FieldU32LE("stream_serial_number"))
			pageSequenceNo := uint32(d.FieldU32LE("page_sequence_no"))
			d.FieldU32("page_checksum")
			pageSegments := d.FieldU8("page_segments")
			segmentTable := d.FieldBytesLen("segment_table", uint64(pageSegments))

			s, sFound := d.streams[streamSerialNumber]
			if !sFound {
				s = &stream{sequenceNo: pageSequenceNo}
				d.streams[streamSerialNumber] = s
			}

			if !sFound && !isFirstPage {
				// TODO: not first page and we haven't seen the stream before
				log.Println("not first page and we haven't seen the stream before")
			}
			hasData := len(s.packetBuf) > 0
			if isContinuedPacket && !hasData {
				// TODO: continuation but we haven't seen any packet data yet
				log.Println("continuation but we haven't seen any packet data yet")
			}
			if !isFirstPage && s.sequenceNo+1 != pageSequenceNo {
				// TODO: page gap
				log.Println("page gap")
			}

			for _, ss := range segmentTable {
				bs := d.FieldBytesLen("segment", uint64(ss))
				s.packetBuf = append(s.packetBuf, bs...)
				if len(bs) < 255 {
					d.FieldDecodeBitBuf("packet", 0, 0, bitbuf.NewFromBytes(s.packetBuf), []string{"vorbis"})
					s.packetBuf = nil
				}
			}

			s.sequenceNo = pageSequenceNo
			if isLastPage {
				delete(d.streams, streamSerialNumber)
			}
		})
	}
}
