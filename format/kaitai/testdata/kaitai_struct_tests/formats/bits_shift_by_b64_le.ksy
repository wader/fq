# If the particular runtime library implements `read_bits_int_le()` according to the common
# algorithm introduced in https://github.com/kaitai-io/kaitai_struct/issues/949 (as all runtime
# libs should), this test will make it store `res >> (bitsNeeded: 64)` in `bits` for subsequent
# bit integers. However, the behavior of `x >> 64` is often problematic in languages with
# 64-bit operators - this case usually has to be handled with special care. If it's not, this
# test will reveal it.
meta:
  id: bits_shift_by_b64_le
  bit-endian: le
seq:
  - id: a
    type: b64
  - id: b
    type: b8
