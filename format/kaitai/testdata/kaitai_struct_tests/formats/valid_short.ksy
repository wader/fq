meta:
  id: valid_short
  encoding: utf-8
  endian: le
seq:
  - id: magic1
    size: 6
    valid: '[0x50, 0x41, 0x43, 0x4b, 0x2d, 0x31]'
  - id: uint8
    type: u1
    valid: 255
  - id: sint8
    type: s1
    valid: -1
  - id: magic_uint
    type: str
    size: 10
    valid: '"PACK-U-DEF"'
  - id: uint16
    type: u2
    valid: 65535
  - id: uint32
    type: u4
    valid: 4294967295
  - id: uint64
    type: u8
    valid: 18446744073709551615
  - id: magic_sint
    type: str
    size: 10
    valid: '"PACK-S-DEF"'
  - id: sint16
    type: s2
    valid: -1
  - id: sint32
    type: s4
    valid: -1
  - id: sint64
    type: s8
    valid: -1
