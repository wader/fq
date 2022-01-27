# Tests "pad-right" and "terminator" functionality in fixed-length byte arrays
meta:
  id: bytes_pad_term
seq:
  - id: str_pad
    size: 20
    pad-right: 0x40
  - id: str_term
    size: 20
    terminator: 0x40
  - id: str_term_and_pad
    size: 20
    terminator: 0x40
    pad-right: 0x2b
  - id: str_term_include
    size: 20
    terminator: 0x40
    include: true
