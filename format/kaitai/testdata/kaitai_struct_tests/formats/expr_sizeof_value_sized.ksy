meta:
  id: expr_sizeof_value_sized
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
seq:
  - id: block1
    type: block
    size: 12
  - id: more
    type: u2
instances:
  self_sizeof:
    value: _sizeof
  sizeof_block:
    value: block1._sizeof
  sizeof_block_a:
    value: block1.a._sizeof
  sizeof_block_b:
    value: block1.b._sizeof
  sizeof_block_c:
    value: block1.c._sizeof
