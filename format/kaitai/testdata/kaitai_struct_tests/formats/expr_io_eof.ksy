# Tests _io.eof operation
meta:
  id: expr_io_eof
  endian: le
seq:
  - id: substream1
    size: 4
    type: one_or_two
  - id: substream2
    size: 8
    type: one_or_two
types:
  one_or_two:
    seq:
      - id: one
        type: u4
      - id: two
        type: u4
        if: not _io.eof
    instances:
      reflect_eof:
        value: _io.eof
