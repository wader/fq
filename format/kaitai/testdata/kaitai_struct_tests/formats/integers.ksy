meta:
  id: integers
  endian: le
seq:
  - id: magic1
    contents: 'PACK-1'
  - id: uint8
    type: u1
  - id: sint8
    type: s1
  - id: magic_uint
    contents: 'PACK-U-DEF'
  - id: uint16
    type: u2
  - id: uint32
    type: u4
  - id: uint64
    type: u8
  - id: magic_sint
    contents: 'PACK-S-DEF'
  - id: sint16
    type: s2
  - id: sint32
    type: s4
  - id: sint64
    type: s8
  - id: magic_uint_le
    contents: 'PACK-U-LE'
  - id: uint16le
    type: u2le
  - id: uint32le
    type: u4le
  - id: uint64le
    type: u8le
  - id: magic_sint_le
    contents: 'PACK-S-LE'
  - id: sint16le
    type: s2le
  - id: sint32le
    type: s4le
  - id: sint64le
    type: s8le
  - id: magic_uint_be
    contents: 'PACK-U-BE'
  - id: uint16be
    type: u2be
  - id: uint32be
    type: u4be
  - id: uint64be
    type: u8be
  - id: magic_sint_be
    contents: 'PACK-S-BE'
  - id: sint16be
    type: s2be
  - id: sint32be
    type: s4be
  - id: sint64be
    type: s8be
