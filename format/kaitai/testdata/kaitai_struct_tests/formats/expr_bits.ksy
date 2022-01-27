meta:
  id: expr_bits
seq:
  - id: enum_seq
    type: b2
    enum: items
  - id: a
    type: b3
  - id: byte_size
    size: a
  - id: repeat_expr
    type: s1
    repeat: expr
    repeat-expr: a
  - id: switch_on_type
    type:
      switch-on: a
      cases:
        2: s1
  - id: switch_on_endian
    type: endian_switch
instances:
  enum_inst:
    value: a
    enum: items
  inst_pos:
    pos: a
    type: s1
types:
  endian_switch:
    meta:
      endian:
        switch-on: _parent.a
        cases:
          1: le
          2: be
    seq:
      - id: foo
        type: s2
enums:
  items:
    1: foo
    2: bar
