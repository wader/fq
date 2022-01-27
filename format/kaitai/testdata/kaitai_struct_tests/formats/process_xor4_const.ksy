meta:
  id: process_xor4_const
seq:
  - id: key
    size: 4
  - id: buf
    size-eos: true
    process: xor([0xec, 0xbb, 0xa3, 0x14])
