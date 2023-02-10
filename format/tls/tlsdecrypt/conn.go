// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// TLS low level connection and record layer

package tlsdecrypt

import (
	"crypto/cipher"
	"crypto/subtle"
	"hash"
	"sync"
)

// A halfConn represents one direction of the record layer
// connection, either sending or receiving.
type halfConn struct {
	sync.Mutex

	err     error  // first permanent error
	version uint16 // protocol version
	cipher  any    // cipher algorithm
	mac     hash.Hash
	seq     [8]byte // 64-bit sequence number

	scratchBuf [13]byte // to avoid allocs; interface method args escape

	nextCipher any       // next encryption state
	nextMac    hash.Hash // next MAC algorithm

	trafficSecret []byte // current TLS 1.3 traffic secret
}

// incSeq increments the sequence number.
func (hc *halfConn) incSeq() {
	for i := 7; i >= 0; i-- {
		hc.seq[i]++
		if hc.seq[i] != 0 {
			return
		}
	}

	// Not allowed to let sequence number wrap.
	// Instead, must renegotiate before it does.
	// Not likely enough to bother.
	panic("TLS: sequence number wraparound")
}

// explicitNonceLen returns the number of bytes of explicit nonce or IV included
// in each record. Explicit nonces are present only in CBC modes after TLS 1.0
// and in certain AEAD modes in TLS 1.2.
func (hc *halfConn) explicitNonceLen() int {
	if hc.cipher == nil {
		return 0
	}

	switch c := hc.cipher.(type) {
	case cipher.Stream:
		return 0
	case aead:
		return c.explicitNonceLen()
	case cbcMode:
		// TLS 1.1 introduced a per-record explicit IV to fix the BEAST attack.
		if hc.version >= VersionTLS11 {
			return c.BlockSize()
		}
		return 0
	default:
		panic("unknown cipher type")
	}
}

// extractPadding returns, in constant time, the length of the padding to remove
// from the end of payload. It also returns a byte which is equal to 255 if the
// padding was valid and 0 otherwise. See RFC 2246, Section 6.2.3.2.
func extractPadding(payload []byte) (toRemove int, good byte) {
	if len(payload) < 1 {
		return 0, 0
	}

	paddingLen := payload[len(payload)-1]
	t := uint(len(payload)-1) - uint(paddingLen)
	// if len(payload) >= (paddingLen - 1) then the MSB of t is zero
	good = byte(int32(^t) >> 31)

	// The maximum possible padding length plus the actual length field
	toCheck := 256
	// The length of the padded data is public, so we can use an if here
	if toCheck > len(payload) {
		toCheck = len(payload)
	}

	for i := 0; i < toCheck; i++ {
		t := uint(paddingLen) - uint(i)
		// if i <= paddingLen then the MSB of t is zero
		mask := byte(int32(^t) >> 31)
		b := payload[len(payload)-1-i]
		good &^= mask&paddingLen ^ mask&b
	}

	// We AND together the bits of good and replicate the result across
	// all the bits.
	good &= good << 4
	good &= good << 2
	good &= good << 1
	good = uint8(int8(good) >> 7)

	// Zero the padding length on error. This ensures any unchecked bytes
	// are included in the MAC. Otherwise, an attacker that could
	// distinguish MAC failures from padding failures could mount an attack
	// similar to POODLE in SSL 3.0: given a good ciphertext that uses a
	// full block's worth of padding, replace the final block with another
	// block. If the MAC check passed but the padding check failed, the
	// last byte of that block decrypted to the block size.
	//
	// See also macAndPaddingGood logic below.
	paddingLen &= good

	toRemove = int(paddingLen) + 1
	return
}

func roundUp(a, b int) int {
	return a + (b-a%b)%b
}

// cbcMode is an interface for block ciphers using cipher block chaining.
type cbcMode interface {
	cipher.BlockMode
	SetIV([]byte)
}

// decrypt authenticates and decrypts the record if protection is active at
// this stage. The returned plaintext might overlap with the input.
func (hc *halfConn) decrypt(record []byte) ([]byte, recordType, error) {
	var plaintext []byte
	typ := recordType(record[0])
	payload := record[recordHeaderLen:]

	// In TLS 1.3, change_cipher_spec messages are to be ignored without being
	// decrypted. See RFC 8446, Appendix D.4.
	if hc.version == VersionTLS13 && typ == recordTypeChangeCipherSpec {
		return payload, typ, nil
	}

	paddingGood := byte(255)
	paddingLen := 0

	explicitNonceLen := hc.explicitNonceLen()

	if hc.cipher != nil {
		switch c := hc.cipher.(type) {
		case cipher.Stream:
			c.XORKeyStream(payload, payload)
		case aead:
			if len(payload) < explicitNonceLen {
				return nil, 0, alertBadRecordMAC
			}
			nonce := payload[:explicitNonceLen]
			if len(nonce) == 0 {
				nonce = hc.seq[:]
			}
			payload = payload[explicitNonceLen:]

			var additionalData []byte
			if hc.version == VersionTLS13 {
				additionalData = record[:recordHeaderLen]
			} else {
				additionalData = append(hc.scratchBuf[:0], hc.seq[:]...)
				additionalData = append(additionalData, record[:3]...)
				n := len(payload) - c.Overhead()
				additionalData = append(additionalData, byte(n>>8), byte(n))
			}

			var err error
			plaintext, err = c.Open(payload[:0], nonce, payload, additionalData)
			if err != nil {
				return nil, 0, alertBadRecordMAC
			}
		case cbcMode:
			blockSize := c.BlockSize()
			minPayload := explicitNonceLen + roundUp(hc.mac.Size()+1, blockSize)
			if len(payload)%blockSize != 0 || len(payload) < minPayload {
				return nil, 0, alertBadRecordMAC
			}

			if explicitNonceLen > 0 {
				c.SetIV(payload[:explicitNonceLen])
				payload = payload[explicitNonceLen:]
			}
			c.CryptBlocks(payload, payload)

			// In a limited attempt to protect against CBC padding oracles like
			// Lucky13, the data past paddingLen (which is secret) is passed to
			// the MAC function as extra data, to be fed into the HMAC after
			// computing the digest. This makes the MAC roughly constant time as
			// long as the digest computation is constant time and does not
			// affect the subsequent write, modulo cache effects.
			paddingLen, paddingGood = extractPadding(payload)
		default:
			panic("unknown cipher type")
		}

		if hc.version == VersionTLS13 {
			if typ != recordTypeApplicationData {
				return nil, 0, alertUnexpectedMessage
			}
			if len(plaintext) > maxPlaintext+1 {
				return nil, 0, alertRecordOverflow
			}
			// Remove padding and find the ContentType scanning from the end.
			for i := len(plaintext) - 1; i >= 0; i-- {
				if plaintext[i] != 0 {
					typ = recordType(plaintext[i])
					plaintext = plaintext[:i]
					break
				}
				if i == 0 {
					return nil, 0, alertUnexpectedMessage
				}
			}
		}
	} else {
		plaintext = payload
	}

	if hc.mac != nil {
		macSize := hc.mac.Size()
		if len(payload) < macSize {
			return nil, 0, alertBadRecordMAC
		}

		n := len(payload) - macSize - paddingLen
		n = subtle.ConstantTimeSelect(int(uint32(n)>>31), 0, n) // if n < 0 { n = 0 }
		record[3] = byte(n >> 8)
		record[4] = byte(n)
		remoteMAC := payload[n : n+macSize]
		localMAC := tls10MAC(hc.mac, hc.scratchBuf[:0], hc.seq[:], record[:recordHeaderLen], payload[:n], payload[n+macSize:])

		// This is equivalent to checking the MACs and paddingGood
		// separately, but in constant-time to prevent distinguishing
		// padding failures from MAC failures. Depending on what value
		// of paddingLen was returned on bad padding, distinguishing
		// bad MAC from bad padding can lead to an attack.
		//
		// See also the logic at the end of extractPadding.
		macAndPaddingGood := subtle.ConstantTimeCompare(localMAC, remoteMAC) & int(paddingGood)
		if macAndPaddingGood != 1 {
			return nil, 0, alertBadRecordMAC
		}

		plaintext = payload[:n]
	}

	hc.incSeq()
	return plaintext, typ, nil
}
