meta:
  id: to_string_custom
  encoding: UTF-8
seq:
  - id: s1
    type: str
    terminator: 0x7c
  - id: s2
    type: str
    terminator: 0x7c
to-string: |
  "s1 = " + s1 + ", s2 = " + s2
