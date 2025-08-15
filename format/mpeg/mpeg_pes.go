package mpeg

// TODO: probeable?

// http://dvdnav.mplayerhq.hu/dvdinfo/mpeghdrs.html

import (
	"github.com/wader/fq/format"
	"github.com/wader/fq/pkg/bitio"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
)

var pesPacketGroup decode.Group
var mpegSpuGroup decode.Group

func init() {
	interp.RegisterFormat(
		format.MPEG_PES,
		&decode.Format{
			Description: "MPEG Packetized elementary stream",
			DecodeFn:    pesDecode,
			RootArray:   true,
			RootName:    "packets",
			Dependencies: []decode.Dependency{
				{Groups: []*decode.Group{format.MPEG_PES_Packet}, Out: &pesPacketGroup},
				{Groups: []*decode.Group{format.MPEG_SPU}, Out: &mpegSpuGroup},
			},
		})
}

type subStream struct {
	b []byte
	l int
}

func pesDecode(d *decode.D) any {
	substreams := map[int]*subStream{}

	prefix := d.PeekUintBits(24)
	if prefix != 0b0000_0000_0000_0000_0000_0001 {
		d.Errorf("no pes prefix found")
	}

	i := 0

	spuD := d.FieldArrayValue("spus")

	for d.NotEnd() {
		dv, v, err := d.TryFieldFormat("packet", &pesPacketGroup, nil)
		if dv == nil || err != nil {
			break
		}

		switch dvv := v.(type) {
		case subStreamPacket:
			s, ok := substreams[dvv.number]
			if !ok {
				s = &subStream{}
				substreams[dvv.number] = s
			}
			s.b = append(s.b, dvv.buf...)

			if s.l == 0 && len(s.b) >= 2 {
				s.l = int(s.b[0])<<8 | int(s.b[1])
				// TODO: zero l?
			}

			// TODO: is this how spu end is signalled?
			if s.l == len(s.b) {
				spuD.FieldFormatBitBuf("spu", bitio.NewBitReader(s.b, -1), &mpegSpuGroup, nil)
				s.b = nil
				s.l = 0
			}
		}

		i++
	}

	return nil
}
