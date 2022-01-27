# Tests modulo operation, both positive and negative
meta:
  id: expr_mod
  endian: le
seq:
  - id: int_u
    type: u4
  - id: int_s
    type: s4
instances:
  mod_pos_const:
    value: 9837 % 13
  mod_neg_const:
    value: -9837 % 13
  mod_pos_seq:
    value: int_u % 13
  mod_neg_seq:
    value: int_s % 13
