# Tests bytearray comparisons
meta:
  id: expr_bytes_cmp
seq:
  - id: one
    size: 1
  - id: two
    size: 3
instances:
  ack:
    value: '[65, 67, 75]'
  ack2:
    value: '[65, 67, 75, 50]'
  hi_val:
    value: '[0x90, 67]'
  is_eq:
    value: two == ack
  is_ne:
    value: two != ack
  is_lt:
    value: two < ack2
  is_gt:
    value: two > ack2
  is_le:
    value: two <= ack2
  is_ge:
    value: two >= ack2
  is_lt2:
    value: one < two
  is_gt2:
    value: hi_val > two
