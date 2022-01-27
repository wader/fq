# https://github.com/kaitai-io/kaitai_struct/issues/494
meta:
  id: switch_repeat_expr_invalid
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
        0xff: one # there is actually 0x11 in the file
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
