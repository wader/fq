# Tests _io.pos operation
meta:
  id: expr_io_pos
  endian: le
seq:
  - id: substream1
    size: 16
    type: all_plus_number
  - id: substream2
    size: 14
    type: all_plus_number
types:
  all_plus_number:
    seq:
      - id: my_str
        type: strz
        encoding: UTF-8
      - id: body
        size: _io.size - _io.pos - 2
      - id: number
        type: u2
