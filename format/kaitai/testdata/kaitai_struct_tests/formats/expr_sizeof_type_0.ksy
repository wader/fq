meta:
  id: expr_sizeof_type_0
  endian: le
types:
  block:
    seq:
      - id: a
        type: u1
      - id: b
        type: u4
      - id: c
        size: 2
instances:
  sizeof_block:
    value: sizeof<block>
