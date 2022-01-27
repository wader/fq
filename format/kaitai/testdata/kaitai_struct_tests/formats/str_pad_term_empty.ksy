# Same as "str_pad_term", but used with different input file that is
# meant to test fully empty strings
meta:
  id: str_pad_term_empty
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
