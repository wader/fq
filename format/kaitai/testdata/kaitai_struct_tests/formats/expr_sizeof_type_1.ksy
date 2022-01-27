meta:
  id: expr_sizeof_type_1
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
      - id: d
        type: subblock
    types:
      subblock:
        seq:
          - id: a
            size: 4
instances:
  sizeof_block:
    value: sizeof<block>
  sizeof_subblock:
    value: sizeof<block::subblock>
