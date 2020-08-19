package hexdump

const hextable = "0123456789abcdef"

func Hexpairs(b []byte, width int, offset int) string {
	if len(b) == 0 {
		return ""
	}
	t := offset + len(b)
	s := make([]byte, t*3-1)

	for i := 0; i < t; i++ {
		if i < offset {
			s[i*3+0] = ' '
			s[i*3+1] = ' '
		} else {
			v := b[i-offset]
			s[i*3+0] = hextable[v>>4]
			s[i*3+1] = hextable[v&0xf]
		}
		if i != t-1 {
			if i > 0 && i%width == width-1 {
				s[i*3+2] = '\n'
			} else {
				s[i*3+2] = ' '
			}
		}
	}

	return string(s[0 : t*3-1])
}

func Printable(b []byte, offset int) string {
	t := offset + len(b)
	s := make([]byte, t)

	for i := 0; i < t; i++ {
		if i < offset {
			s[i] = ' '
		} else {
			v := b[i-offset]
			if v < 32 || v > 126 {
				s[i] = '.'
			} else {
				s[i] = v
			}
		}
	}

	return string(s)
}
