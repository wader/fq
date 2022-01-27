meta:
  id: process_repeat_usertype
  endian: le
seq:
  - id: blocks
    size: 5
    type: block
    repeat: expr
    repeat-expr: 2
    process: xor(0x9e)
types:
  block:
    seq:
      - id: a
        type: s4
      - id: b
        type: s1
