meta:
  id: process_xor4_value
seq:
  - id: key
    size: 4
  - id: buf
    size-eos: true
    process: xor(key)
