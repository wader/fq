meta:
  id: process_xor_value
  endian: le
seq:
  - id: key
    type: u1
  - id: buf
    size-eos: true
    process: xor(key)
