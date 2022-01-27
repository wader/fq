meta:
  id: repeat_until_s4
  endian: le
seq:
  - id: entries
    type: s4
    repeat: until
    repeat-until: _ == -1
  - id: afterall
    type: strz
    encoding: ASCII
