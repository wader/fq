meta:
  id: process_custom
seq:
  - id: buf1
    size: 5
    process: my_custom_fx(7, true, [0x20, 0x30, 0x40])
  - id: buf2
    size: 5
    process: nested.deeply.custom_fx(7)
  - id: key
    type: u1
  - id: buf3
    size: 5
    process: my_custom_fx(key, false, [0x00])
