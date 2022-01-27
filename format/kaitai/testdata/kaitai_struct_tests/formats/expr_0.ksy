meta:
  id: expr_0
  endian: le
seq:
  - id: len_of_1
    type: u2
instances:
  must_be_f7:
    value: 7 + 0xf0
  must_be_abc123:
    value: '"abc" + "123"'
