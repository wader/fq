package mpeg

// TODO: probeable?

// http://dvdnav.mplayerhq.hu/dvdinfo/mpeghdrs.html

import (
	"log"

	"github.com/wader/fq/format"
	"github.com/wader/fq/format/registry"
	"github.com/wader/fq/pkg/bitio"
	"github.com/wader/fq/pkg/decode"
)

var pesPacketFormat decode.Group
var spuFormat decode.Group

func init() {
	registry.MustRegister(decode.Format{
		Name:        format.MPEG_PES,
		Description: "MPEG Packetized elementary stream",
		DecodeFn:    pesDecode,
		RootArray:   true,
		RootName:    "packets",
		Dependencies: []decode.Dependency{
			{Names: []string{format.MPEG_PES_PACKET}, Group: &pesPacketFormat},
			{Names: []string{format.MPEG_SPU}, Group: &spuFormat},
		},
	})
}

type subStream struct {
	b []byte
	l int
}

func pesDecode(d *decode.D, in interface{}) interface{} {
	substreams := map[int]*subStream{}

	prefix := d.PeekBits(24)
	if prefix != 0b0000_0000_0000_0000_0000_0001 {
		d.Errorf("no pes prefix found")
	}

	i := 0

	spuD := d.FieldArrayValue("spus")

	for d.NotEnd() {
		dv, v, err := d.TryFieldFormat("packet", pesPacketFormat, nil)
		if dv == nil || err != nil {
			log.Printf("errs[0]: %#+v\n", err)
			break
		}

		switch dvv := v.(type) {
		case subStreamPacket:
			s, ok := substreams[dvv.number]
			if !ok {
				s = &subStream{}
				substreams[dvv.number] = s
			}
			b, _ := dvv.bb.BytesRange(0, int(dvv.bb.Len()/8))
			s.b = append(s.b, b...)

			if s.l == 0 && len(b) >= 2 {
				s.l = int(b[0])<<8 | int(b[1])
				// TODO: zero l?
			}

			// TODO: is this how spu end is signalled?
			if s.l == len(s.b) {
				spuD.FieldFormatBitBuf("spu", bitio.NewBufferFromBytes(s.b, -1), spuFormat, nil)
				s.b = nil
				s.l = 0
			}
		}

		i++
	}

	return nil
}
