package fairplay

// https://github.com/easonlin404/ksm/blob/master/ksm.go

import (
	"github.com/wader/fq/format"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
)

func init() {
	interp.RegisterFormat(
		format.Fairplay_SPC,
		&decode.Format{
			Description: "FairPlay Server Playback Context",
			DecodeFn:    fairPlaySPCDecode,
		})
}

func fairPlaySPCDecode(d *decode.D) any {
	d.FieldU32("version")
	d.FieldRawLen("reserved", 32)
	d.FieldRawLen("iv", 16*8)
	d.FieldRawLen("aes_key_oaep", 128*8)
	d.FieldRawLen("certificate_hash", 20*8)
	payloadLen := d.FieldU32("payload_length")
	d.FieldRawLen("payload", int64(payloadLen)*8)

	return nil
}
