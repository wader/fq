meta:
  id: if_values
  endian: le
seq:
  - id: codes
    type: code
    repeat: expr
    repeat-expr: 3
types:
  code:
    seq:
      - id: opcode
        type: u1
    instances:
      half_opcode:
        value: opcode / 2
        if: opcode % 2 == 0
