# Tests reading un/signed 64-bit integers in languages representing integers as 64-bit floats ("double"s)
# It is intended especially for JavaScript.
meta:
  id: integers_double_overflow
seq:
  - id: signed_safe_min_be # 0x00
    type: s8be
  - id: signed_safe_min_le # 0x08
    type: s8le
  - id: signed_safe_max_be # 0x10
    type: s8be
  - id: signed_safe_max_le # 0x18
    type: s8le
  - id: signed_unsafe_neg_be # 0x20
    type: s8be
  - id: signed_unsafe_neg_le # 0x28
    type: s8le
  - id: signed_unsafe_pos_be # 0x30
    type: s8be
  - id: signed_unsafe_pos_le # 0x38
    type: s8le
instances:
  unsigned_safe_max_be:
    pos: 0x10
    type: u8be
  unsigned_safe_max_le:
    pos: 0x18
    type: u8le
  unsigned_unsafe_pos_be:
    pos: 0x30
    type: u8be
  unsigned_unsafe_pos_le:
    pos: 0x38
    type: u8le
