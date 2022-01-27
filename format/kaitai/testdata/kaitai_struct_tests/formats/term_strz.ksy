meta:
  id: term_strz
  endian: le
seq:
  - id: s1
    type: str
    encoding: UTF-8
    terminator: 0x7c
  - id: s2
    type: str
    encoding: UTF-8
    terminator: 0x7c
    consume: false
  - id: s3
    type: str
    encoding: UTF-8
    terminator: 0x40
    include: true
