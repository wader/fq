meta:
  id: expr_1
  endian: le
seq:
  - id: len_of_1
    type: u2
  - id: str1
    type: str
    size: len_of_1_mod
    encoding: ASCII
instances:
  len_of_1_mod:
    value: len_of_1 - 2
  str1_len:
    value: str1.length
