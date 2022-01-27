meta:
  id: switch_repeat_expr
  endian: le
seq:
  - id: code
    type: u1
  - id: size
    type: u4
  - id: body
    repeat: expr
    repeat-expr: 1
    size: size
    type:
      switch-on: code
      cases:
        0x11: one
        0x22: two
types:
  one:
    seq:
      - id: first
        size-eos: true
  two:
    seq:
      - id: second
        size-eos: true
