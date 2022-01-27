meta:
  id: repeat_n_struct
  endian: le
seq:
  - id: qty
    type: u4
  - id: chunks
    type: chunk
    repeat: expr
    repeat-expr: qty
types:
  chunk:
    seq:
      - id: offset
        type: u4
      - id: len
        type: u4
