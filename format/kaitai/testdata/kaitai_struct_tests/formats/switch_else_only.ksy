meta:
  id: switch_else_only
seq:
  - id: opcode
    type: s1
  - id: prim_byte
    type:
      switch-on: opcode
      cases:
        _: s1
  - id: indicator
    size: 4
  - id: ut
    type:
      switch-on: indicator
      cases:
        _: data
types:
  data:
    seq:
      - id: value
        size: 4
