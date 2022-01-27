meta:
  id: bits_byte_aligned
seq:
  - id: one
    type: b6
  # skips 2 bits
  - id: byte_1
    type: u1
  - id: two
    type: b3
  - id: three
    type: b1
  # skips 4 bits
  - id: byte_2
    type: u1
  - id: four
    type: b14
  # skips 2 bits
  - id: byte_3
    size: 1
  - id: full_byte
    type: b8
  - id: byte_4
    type: u1
