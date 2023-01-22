package mpeg

import (
	"github.com/wader/fq/format"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
	"github.com/wader/fq/pkg/scalar"
)

func init() {
	interp.RegisterFormat(
		format.MPEG_TS_Packet,
		&decode.Format{
			Description: "MPEG Transport Stream Packet",
			DecodeFn:    tsPacketDecode,
		})
}

func tsPacketDecode(d *decode.D) any {
	var mtpi format.MpegTsPacketIn

	if d.ArgAs(&mtpi) {
		mtpi.ProgramMap = map[int]format.MpegTsProgram{}
		mtpi.StreamMap = map[int]format.MpegTsStream{}
		mtpi.ContinuityMap = map[int]int{}
	}

	var mtpo format.MpegTsPacketOut

	d.FieldU8("sync", scalar.UintHex, d.UintAssert(0x47))
	mtpo.TransportErrorIndicator = d.FieldBool("transport_error_indicator", d.BoolAssert(false))
	mtpo.PayloadUnitStart = d.FieldBool("payload_unit_start")
	d.FieldBool("transport_priority")
	pid := d.FieldU13("pid", tsPidMap, scalar.UintHex)
	if p, ok := mtpi.ProgramMap[int(pid)]; ok {
		d.FieldValueUint("program", uint64(p.Number), scalar.UintHex)
	} else if s, ok := mtpi.StreamMap[int(pid)]; ok {
		if p, ok := mtpi.ProgramMap[s.ProgramPid]; ok {
			d.FieldValueUint("program", uint64(p.Number), scalar.UintHex)
		}
		d.FieldValueUint("stream_type", uint64(s.Type), tsStreamTypeMap)
	}
	mtpo.Pid = int(pid)
	mtpo.TransportScramblingControl = int(d.FieldU2("transport_scrambling_control", scalar.UintMapSymStr{
		0b00: "not_scrambled",
		0b01: "reserved",
		0b10: "even_key",
		0b11: "odd_key",
	}))
	adaptationFieldControl := d.FieldU2("adaptation_field_control", scalar.UintMapSymStr{
		0b00:                              "reserved",
		adaptationFieldControlPayloadOnly: "payload_only",
		adaptationFieldControlAdaptationFieldOnly:       "adaptation_field_only",
		adaptationFieldControlAdaptationFieldAndPayload: "adaptation_and_payload",
	})
	mtpo.ContinuityCounter = int(d.FieldU4("continuity_counter", scalar.UintFn(func(s scalar.Uint) (scalar.Uint, error) {
		prev, prevFound := mtpi.ContinuityMap[int(pid)]
		current := int(s.Actual)

		switch {
		case prevFound && (prev+1)&0xf == current:
			s.Description = "continuous"
		case prevFound:
			s.Description = "non-continuous"
		default:
			s.Description = "unknown"
		}

		return s, nil
	})))

	switch adaptationFieldControl {
	case adaptationFieldControlAdaptationFieldOnly,
		adaptationFieldControlAdaptationFieldAndPayload:
		d.FieldStruct("adaptation_field", func(d *decode.D) {
			length := d.FieldU8("length") // Number of bytes in the adaptation field immediately following this byte
			d.FramedFn(int64(length)*8, func(d *decode.D) {
				d.FieldBool("discontinuity_indicator")                                               // Set if current TS packet is in a discontinuity state with respect to either the continuity counter or the program clock reference
				d.FieldBool("random_access_indicator")                                               // Set when the stream may be decoded without errors from this point
				d.FieldBool("elementary_stream_priority_indicator")                                  // Set when this stream should be considered "high priority"
				pcrPresent := d.FieldBool("pcr_present")                                             // Set when PCR field is present
				opcrPresent := d.FieldBool("opcr_present")                                           // Set when OPCR field is present
				splicingPointPresent := d.FieldBool("splicing_point_present")                        // Set when splice countdown field is present
				transportPrivatePresent := d.FieldBool("transport_private_present")                  // Set when transport private data is present
				adaptationFieldExtensionPresent := d.FieldBool("adaptation_field_extension_present") // Set when adaptation extension data is present
				if pcrPresent {
					d.FieldU("pcr", 48)
				}
				if opcrPresent {
					d.FieldU("opcr", 48)
				}
				if splicingPointPresent {
					d.FieldU8("splicing_point")
				}
				if transportPrivatePresent {
					d.FieldStruct("transport_private", func(d *decode.D) {
						length := d.FieldU8("length")
						d.FieldRawLen("data", int64(length)*8)
					})
				}
				if adaptationFieldExtensionPresent {
					d.FieldStruct("adaptation_extension", func(d *decode.D) {
						length := d.FieldU8("length")
						d.FramedFn(int64(length)*8, func(d *decode.D) {
							d.FieldBool("legal_time_window")
							d.FieldBool("piecewise_rate")
							d.FieldBool("seamless_splice")
							d.FieldU5("reserved", scalar.UintHex)
							d.FieldRawLen("data", d.BitsLeft())
						})
					})

					// Optional fields
					// LTW flag set (2 bytes)
					// LTW valid flag	1	0x8000
					// LTW offset	15	0x7fff	Extra information for rebroadcasters to determine the state of buffers when packets may be missing.
					// Piecewise flag set (3 bytes)
					// Reserved	2	0xc00000
					// Piecewise rate	22	0x3fffff	The rate of the stream, measured in 188-byte packets, to define the end-time of the LTW.
					// Seamless splice flag set (5 bytes)
					// Splice type	4	0xf000000000	Indicates the parameters of the H.262 splice.
					// DTS next access unit	36	0x0efffefffe	The PES DTS of the splice point. Split up as multiple fields, 1 marker bit (0x1), 15 bits, 1 marker bit, 15 bits, and 1 marker bit, for 33 data bits total.
				}
				if d.BitsLeft() > 0 {
					d.FieldRawLen("stuffing", d.BitsLeft())
				}
			})
		})
	}

	isTable := tsPidIsTable(mtpo.Pid, mtpi.ProgramMap)
	if isTable {
		var payloadPointer uint64
		if mtpo.PayloadUnitStart {
			payloadPointer = d.FieldU8("payload_pointer")
		}
		if payloadPointer > 0 {
			d.FieldRawLen("stuffing", int64(payloadPointer)*8)
		}
	}

	switch adaptationFieldControl {
	case adaptationFieldControlPayloadOnly,
		adaptationFieldControlAdaptationFieldAndPayload:
		payload := d.FieldRawLen("payload", d.BitsLeft())
		mtpo.Payload = d.ReadAllBits(payload)
	default:
		// TODO: unknown adaption control flags
		d.FieldRawLen("unknown", d.BitsLeft())
	}

	return mtpo
}
