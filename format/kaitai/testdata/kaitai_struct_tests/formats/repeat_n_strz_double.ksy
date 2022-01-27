meta:
  id: repeat_n_strz_double
  endian: le
seq:
  - id: qty
    type: u4
  - id: lines1
    type: strz
    encoding: UTF-8
    repeat: expr
    repeat-expr: qty / 2
  - id: lines2
    type: strz
    encoding: UTF-8
    repeat: expr
    repeat-expr: qty / 2
