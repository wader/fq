package checksum

// IPv4 implements hash.Hash
type IPv4 struct {
	sum uint
	odd bool
}

func (c *IPv4) Write(p []byte) (n int, err error) {
	for _, b := range p {
		if c.odd {
			c.sum += uint(b)
			if c.sum > 0xffff {
				c.sum++
				c.sum &= 0xffff
			}
		} else {
			c.sum += uint(b) << 8
		}
		c.odd = !c.odd
	}
	return len(p), nil
}

func (c *IPv4) Sum(b []byte) []byte {
	s := c.sum
	if c.odd {
		if s > 0xffff {
			s++
			s &= 0xffff
		}

	}
	s ^= 0xffff
	return append(b, byte(s>>8), byte(s))
}
func (c *IPv4) Reset()         { c.sum = 0 }
func (c *IPv4) Size() int      { return 2 }
func (c *IPv4) BlockSize() int { return 2 }
