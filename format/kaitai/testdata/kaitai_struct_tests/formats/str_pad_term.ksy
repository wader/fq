# Tests "pad-right" and "terminator" functionality in fixed-length strings ("str" with "size")
meta:
  id: str_pad_term
  encoding: UTF-8
seq:
  - id: str_pad
    type: str
    size: 20
    pad-right: 0x40
  - id: str_term
    type: str
    size: 20
    terminator: 0x40
  - id: str_term_and_pad
    type: str
    size: 20
    terminator: 0x40
    pad-right: 0x2b
  - id: str_term_include
    type: str
    size: 20
    terminator: 0x40
    include: true
  # "consume" is pointless to test: it will be always consumed anyway, as we have a fixed-length string
