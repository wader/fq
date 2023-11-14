// Package rangecoer implements a range coder as per  3.8.1. Range Coding Mode
// of draft-ietf-cellar-ffv1.
package rangecoder

import (
	"log"

	"github.com/wader/fq/internal/mathex"
)

// Cross-references are to
// https://tools.ietf.org/id/draft-ietf-cellar-ffv1-17

// Coder is an instance of a range coder, as defined in:
//
//	Martin, G. Nigel N., "Range encoding: an algorithm for
//	removing redundancy from a digitised message.", July 1979.
type Coder struct {
	readByte   func() (byte, error)
	byteLen    int
	pos        int
	low        uint16
	rng        uint16
	cur_byte   int32
	zero_state [256]uint8
	one_state  [256]uint8
}

// NewCoder creates a new range coder instance.
//
// See: 3.8.1. Range Coding Mode
func NewCoder(low uint16, readByte func() (byte, error), byteLen int) (*Coder, error) {
	ret := new(Coder)
	ret.readByte = readByte
	ret.byteLen = byteLen
	// Figure 15.
	ret.pos = 0 // 2
	// Figure 14.
	ret.low = low
	// Figure 13.
	ret.rng = 0xFF00
	ret.cur_byte = -1
	if ret.low >= ret.rng {
		ret.low = ret.rng
		ret.pos = byteLen - 1
	}

	// 3.8.1.3. Initial Values for the Context Model
	ret.SetTable(DefaultStateTransition)

	return ret, nil
}

// Refills the buffer
func (c *Coder) refill() error {
	log.Printf("refill c.rng: %#+v\n", c.rng)
	// Figure 12.
	if c.rng < 0x100 {
		c.rng = c.rng << 8
		c.low = c.low << 8
		if c.pos < c.byteLen {
			b, err := c.readByte()
			if err != nil {
				return err
			}

			c.low += uint16(b)
			c.pos++
		}
	}
	return nil
}

// Gets the next boolean state
func (c *Coder) get(state *uint8) (bool, error) {
	// Figure 10.
	rangeoff := uint16((uint32(c.rng) * uint32((*state))) >> 8)
	c.rng -= rangeoff
	if c.low < c.rng {
		*state = c.zero_state[int(*state)]
		if err := c.refill(); err != nil {
			return false, err
		}
		return false, nil
	} else {
		c.low -= c.rng
		*state = c.one_state[int(*state)]
		c.rng = rangeoff
		if err := c.refill(); err != nil {
			return false, err
		}
		return true, nil
	}
}

// UR gets the next range coded unsigned scalar symbol.
//
// See: 4. Bitstream
func (c *Coder) UR(state []uint8) (uint32, error) {
	n, err := c.symbol(state, false)
	return uint32(n), err
}

// SR gets the next range coded signed scalar symbol.
//
// See: 4. Bitstream
func (c *Coder) SR(state []uint8) (int32, error) {
	return c.symbol(state, true)
}

// BR gets the next range coded Boolean symbol.
//
// See: 4. Bitstream
func (c *Coder) BR(state []uint8) (bool, error) {
	return c.get(&state[0])
}

// Gets the next range coded symbol.
//
// See: 3.8.1.2. Range Non Binary Values
func (c *Coder) symbol(state []uint8, signed bool) (int32, error) {
	b, err := c.get(&state[0])
	if err != nil {
		return 0, err
	}
	if b {
		return 0, nil
	}

	e := int32(0)
	for {
		b, err := c.get(&state[1+mathex.Min(e, 9)])
		if err != nil {
			return 0, err
		}
		if !b {
			break
		}
		e++
		if e > 31 {
			panic("WTF range coder!")
		}
	}

	a := uint32(1)
	for i := e - 1; i >= 0; i-- {
		a = a * 2
		b, err := c.get(&state[22+mathex.Min(i, 9)])
		if err != nil {
			return 0, err

		}
		if b {
			a++
		}
	}

	if signed {
		b, err := c.get(&state[11+mathex.Min(e, 10)])
		if err != nil {
			return 0, err
		}
		if b {
			return -(int32(a)), nil
		}
	}

	return int32(a), nil
}

func (c *Coder) SetTable(table [256]uint8) {
	// 3.8.1.4. State Transition Table

	// Figure 17.
	for i := 0; i < 256; i++ {
		c.one_state[i] = table[i]
	}
	// Figure 18.
	for i := 1; i < 255; i++ {
		c.zero_state[i] = uint8(uint16(256) - uint16(c.one_state[256-i]))
	}
}

// SentinalEnd ends the current range coder.
//
// See: 3.8.1.1.1. Termination
//   - Sentinal Mode
func (c *Coder) SentinalEnd() error {
	state := uint8(129)
	_, err := c.get(&state)
	return err
}

// GetPos gets the current position in the bitstream.
func (c *Coder) GetPos() int {
	if c.rng < 0x100 {
		return c.pos - 1
	}
	return c.pos
}
