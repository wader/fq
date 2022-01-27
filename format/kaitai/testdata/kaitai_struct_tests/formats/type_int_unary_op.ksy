meta:
  id: type_int_unary_op
  endian: le
seq:
  - id: value_s2
    type: s2
  - id: value_s8
    type: s8
instances:
  unary_s2:
    value: -value_s2
  unary_s8:
    value: -value_s8
