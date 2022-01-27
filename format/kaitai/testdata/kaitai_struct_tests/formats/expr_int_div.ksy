# Tests division operation, both positive and negative
# See https://github.com/kaitai-io/kaitai_struct/issues/746
#  => the KS division operation `a / b` should do `floor(a / b)`
meta:
  id: expr_int_div
  endian: le
seq:
  - id: int_u
    type: u4
  - id: int_s
    type: s4
instances:
  div_pos_const:
    value: 9837 / 13
  div_neg_const:
    value: -9837 / 13
  div_pos_seq:
    value: int_u / 13
  div_neg_seq:
    value: int_s / 13
