meta:
  id: index_sizes
  endian: le
  encoding: ASCII
seq:
  - id: qty
    type: u4
  - id: sizes
    type: u4
    repeat: expr
    repeat-expr: qty
  - id: bufs
    type: str
    size: sizes[_index]
    repeat: expr
    repeat-expr: qty
