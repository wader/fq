package ogg

// https://xiph.org/ogg/doc/framing.html

import (
	"bytes"
	"fmt"

	"github.com/wader/fq/format"
	"github.com/wader/fq/pkg/bitio"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
)

var oggPageFormat decode.Group
var vorbisPacketFormat decode.Group
var opusPacketFormat decode.Group
var flacMetadatablockFormat decode.Group
var flacFrameFormat decode.Group

func init() {
	interp.RegisterFormat(decode.Format{
		Name:        format.OGG,
		Description: "OGG file",
		Groups:      []string{format.PROBE},
		DecodeFn:    decodeOgg,
		Dependencies: []decode.Dependency{
			{Names: []string{format.OGG_PAGE}, Group: &oggPageFormat},
			{Names: []string{format.VORBIS_PACKET}, Group: &vorbisPacketFormat},
			{Names: []string{format.OPUS_PACKET}, Group: &opusPacketFormat},
			{Names: []string{format.FLAC_METADATABLOCK}, Group: &flacMetadatablockFormat},
			{Names: []string{format.FLAC_FRAME}, Group: &flacFrameFormat},
		},
	})
}

var (
	vorbisIdentification = []byte("\x01vorbis")
	opusIdentification   = []byte("OpusHead")
	flacIdentification   = []byte("\x7fFLAC")
)

type streamCodec int

const (
	codecUnknown streamCodec = iota
	codecVorbis
	codecOpus
	codecFlac
)

type stream struct {
	sequenceNo     uint32
	packetBuf      []byte
	packetD        *decode.D
	codec          streamCodec
	flacStreamInfo format.FlacStreamInfo
}

func decodeOgg(d *decode.D, in any) any {
	validPages := 0
	streams := map[uint32]*stream{}
	streamsD := d.FieldArrayValue("streams")

	d.FieldArray("pages", func(d *decode.D) {
		for !d.End() {
			_, dv, _ := d.TryFieldFormat("page", oggPageFormat, nil)
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
				streamsD.FieldStruct("stream", func(d *decode.D) {
					d.FieldValueU("serial_number", uint64(oggPageOut.StreamSerialNumber))
					packetsD = d.FieldArrayValue("packets")
				})
				s = &stream{
					sequenceNo: oggPageOut.SequenceNo,
					packetD:    packetsD,
					codec:      codecUnknown,
				}
				streams[oggPageOut.StreamSerialNumber] = s
			}

			// if !sFound && !oggPageOut.IsFirstPage {
			// 	// TODO: not first page and we haven't seen the stream before
			// 	// log.Println("not first page and we haven't seen the stream before")
			// }
			// hasData := len(s.packetBuf) > 0
			// if oggPageOut.IsContinuedPacket && !hasData {
			// 	// TODO: continuation but we haven't seen any packet data yet
			// 	// log.Println("continuation but we haven't seen any packet data yet")
			// }
			// if !oggPageOut.IsFirstPage && s.sequenceNo+1 != oggPageOut.SequenceNo {
			// 	// TODO: page gap
			// 	// log.Println("page gap")
			// }

			for _, bs := range oggPageOut.Segments {
				s.packetBuf = append(s.packetBuf, bs...)
				if len(bs) < 255 {
					br := bitio.NewBitReader(s.packetBuf, -1)

					if s.codec == codecUnknown {
						if bytes.HasPrefix(s.packetBuf, vorbisIdentification) {
							s.codec = codecVorbis
						} else if bytes.HasPrefix(s.packetBuf, opusIdentification) {
							s.codec = codecOpus
						} else if bytes.HasPrefix(s.packetBuf, flacIdentification) {
							s.codec = codecFlac
						}
					}

					switch s.codec {
					case codecVorbis:
						// TODO: err
						if _, _, err := s.packetD.TryFieldFormatBitBuf("packet", br, vorbisPacketFormat, nil); err != nil {
							s.packetD.FieldRootBitBuf("packet", br)
						}
					case codecOpus:
						// TODO: err
						if _, _, err := s.packetD.TryFieldFormatBitBuf("packet", br, opusPacketFormat, nil); err != nil {
							s.packetD.FieldRootBitBuf("packet", br)
						}
					case codecFlac:
						if len(s.packetBuf) == 0 {
							return
						}

						switch {
						case s.packetBuf[0] == 0x7f:
							s.packetD.FieldStructRootBitBufFn("packet", br, func(d *decode.D) {
								d.FieldU8("type")
								d.FieldUTF8("signature", 4)
								d.FieldU8("major")
								d.FieldU8("minor")
								d.FieldU16("header_packets")
								d.FieldUTF8("flac_signature", 4)
								dv, v := d.FieldFormat("metadatablock", flacMetadatablockFormat, nil)
								flacMetadatablockOut, ok := v.(format.FlacMetadatablockOut)
								if dv != nil && !ok {
									panic(fmt.Sprintf("expected FlacMetadatablockOut, got %#+v", flacMetadatablockOut))
								}
								s.flacStreamInfo = flacMetadatablockOut.StreamInfo
							})
						case s.packetBuf[0] == 0xff:
							s.packetD.FieldFormatBitBuf("packet", br, flacFrameFormat, nil)
						default:
							s.packetD.FieldFormatBitBuf("packet", br, flacMetadatablockFormat, nil)
						}
					case codecUnknown:
						s.packetD.FieldRootBitBuf("packet", br)
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
		d.Fatalf("no pages found")
	}

	return nil
}
