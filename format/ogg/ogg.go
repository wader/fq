package ogg

// https://xiph.org/ogg/doc/framing.html
// TODO: audio/ogg"

import (
	"bytes"
	"fq/format"
	"fq/pkg/bitio"
	"fq/pkg/decode"
)

var oggPage []*decode.Format
var vorbisPacket []*decode.Format
var opusPacket []*decode.Format

func init() {
	format.MustRegister(&decode.Format{
		Name:        format.OGG,
		Description: "OGG file",
		Groups:      []string{format.PROBE},
		DecodeFn:    decodeOgg,
		Dependencies: []decode.Dependency{
			{Names: []string{format.OGG_PAGE}, Formats: &oggPage},
			{Names: []string{format.VORBIS_PACKET}, Formats: &vorbisPacket},
			{Names: []string{format.OPUS_PACKET}, Formats: &opusPacket},
		},
	})
}

var (
	vorbisIdentification = []byte("\x01vorbis")
	opusIdentification   = []byte("OpusHead")
)

type streamCodec int

const (
	codecUnknown streamCodec = iota
	codecVorbis
	codecOpus
)

type stream struct {
	firstBit   int64
	sequenceNo uint32
	packetBuf  []byte
	packetD    *decode.D
	codec      streamCodec
}

func decodeOgg(d *decode.D, in interface{}) interface{} {
	validPages := 0
	streams := map[uint32]*stream{}
	streamsD := d.FieldArray("streams")

	d.FieldArrayFn("pages", func(d *decode.D) {
		for !d.End() {
			_, dv, _ := d.FieldTryDecode("page", oggPage)
			if dv == nil {
				break
			}
			oggPageOut, ok := dv.(format.OggPageOut)
			if !ok {
				panic("page decode is not a oggPageOut")
			}

			s, sFound := streams[oggPageOut.StreamSerialNumber]
			if !sFound {
				var packetsD *decode.D
				streamsD.FieldStructFn("stream", func(d *decode.D) {
					d.FieldValueU("serial_number", uint64(oggPageOut.StreamSerialNumber), "")
					packetsD = d.FieldArray("packets")
				})
				s = &stream{
					sequenceNo: oggPageOut.SequenceNo,
					packetD:    packetsD,
					codec:      codecUnknown,
				}
				streams[oggPageOut.StreamSerialNumber] = s
			}

			if !sFound && !oggPageOut.IsFirstPage {
				// TODO: not first page and we haven't seen the stream before
				// log.Println("not first page and we haven't seen the stream before")
			}
			hasData := len(s.packetBuf) > 0
			if oggPageOut.IsContinuedPacket && !hasData {
				// TODO: continuation but we haven't seen any packet data yet
				// log.Println("continuation but we haven't seen any packet data yet")
			}
			if !oggPageOut.IsFirstPage && s.sequenceNo+1 != oggPageOut.SequenceNo {
				// TODO: page gap
				// log.Println("page gap")
			}

			for _, ps := range oggPageOut.Segments {
				if s.packetBuf == nil {
					s.firstBit = d.Pos()
				}

				// TODO: decoder buffer api that panics?
				psBytes := ps.Len() / 8

				// TODO: cleanup
				b, _ := ps.BytesRange(0, int(psBytes))
				s.packetBuf = append(s.packetBuf, b...)
				if psBytes < 255 { // TODO: list range maps of demuxed packets?
					bb := bitio.NewBufferFromBytes(s.packetBuf, -1)

					if s.codec == codecUnknown {
						if b, err := bb.PeekBytes(len(vorbisIdentification)); err == nil && bytes.Equal(b, vorbisIdentification) {
							s.codec = codecVorbis
						} else if b, err := bb.PeekBytes(len(opusIdentification)); err == nil && bytes.Equal(b, opusIdentification) {
							s.codec = codecOpus
						}
					}

					switch s.codec {
					case codecVorbis:
						s.packetD.FieldTryDecodeBitBuf("packet", bb, vorbisPacket)
					case codecOpus:
						s.packetD.FieldTryDecodeBitBuf("packet", bb, opusPacket)
					case codecUnknown:
						s.packetD.FieldBitBuf("packet", bb)
					}

					s.packetBuf = nil
				}
			}

			s.sequenceNo = oggPageOut.SequenceNo
			if oggPageOut.IsLastPage {
				delete(streams, oggPageOut.StreamSerialNumber)
			}

			validPages++
		}
	})

	if validPages == 0 {
		d.Invalid("no pages found")
	}

	return nil
}
