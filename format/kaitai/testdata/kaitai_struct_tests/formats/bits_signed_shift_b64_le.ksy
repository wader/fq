# Make sure that languages with 64-bit operators (Java, PHP) use *unsigned* (aka
# zero-fill) right shift, not signed (sign-propagating) right shift. This test
# will fail if a signed shift is used.
meta:
  id: bits_signed_shift_b64_le
  bit-endian: le
seq:
  - id: a
    type: b63
  - id: b
    type: b9
