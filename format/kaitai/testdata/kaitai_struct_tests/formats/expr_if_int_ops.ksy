meta:
  id: expr_if_int_ops
  endian: le
seq:
  - id: skip
    size: 2
  - id: it
    type: s2
    if: true
  - id: boxed
    type: s2
    if: true
instances:
  is_eq_prim:
    value: it == 0x4141
  is_eq_boxed:
    value: it == boxed
