# Make sure that languages with 32-bit operators (JavaScript) use *unsigned*
# (aka zero-fill) right shift, not signed (sign-propagating) right shift. This
# test will fail if a signed shift is used.
meta:
  id: bits_signed_shift_b32_le
  bit-endian: le
seq:
  - id: a
    type: b31
  - id: b
    type: b9
