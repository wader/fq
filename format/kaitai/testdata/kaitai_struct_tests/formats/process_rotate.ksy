meta:
  id: process_rotate
  endian: le
seq:
  - id: buf1
    size: 5
    process: rol(3)
  - id: buf2
    size: 5
    process: ror(3)
  - id: key
    type: u1
  - id: buf3
    size: 5
    process: rol(key)
