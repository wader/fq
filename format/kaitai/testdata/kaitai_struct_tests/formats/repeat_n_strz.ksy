meta:
  id: repeat_n_strz
  endian: le
seq:
  - id: qty
    type: u4
  - id: lines
    type: strz
    encoding: UTF-8
    repeat: expr
    repeat-expr: qty
