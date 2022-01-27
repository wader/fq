meta:
  id: valid_long
  encoding: utf-8
  endian: le
seq:
  - id: magic1
    size: 6
    valid:
      eq: '[0x50, 0x41, 0x43, 0x4b, 0x2d, 0x31]'
  - id: uint8
    type: u1
    valid:
      eq: 255
  - id: sint8
    type: s1
    valid:
      eq: -1
  - id: magic_uint
    type: str
    size: 10
    valid:
      eq: '"PACK-U-DEF"'
  - id: uint16
    type: u2
    valid:
      eq: 65535
  - id: uint32
    type: u4
    valid:
      eq: 4294967295
  - id: uint64
    type: u8
    valid:
      eq: 18446744073709551615
  - id: magic_sint
    type: str
    size: 10
    valid:
      eq: '"PACK-S-DEF"'
  - id: sint16
    type: s2
    valid:
      eq: -1
  - id: sint32
    type: s4
    valid:
      eq: -1
  - id: sint64
    type: s8
    valid:
      eq: -1
