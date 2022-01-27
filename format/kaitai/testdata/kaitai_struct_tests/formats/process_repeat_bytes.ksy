meta:
  id: process_repeat_bytes
seq:
  - id: bufs
    size: 5
    repeat: expr
    repeat-expr: 2
    process: xor(0x9e)
