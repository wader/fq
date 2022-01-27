meta:
  id: process_xor_const
  endian: le
seq:
  - id: key
    type: u1
  - id: buf
    size-eos: true
    process: xor(0xff)
