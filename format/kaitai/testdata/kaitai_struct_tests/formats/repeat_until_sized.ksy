meta:
  id: repeat_until_sized
  endian: le
seq:
  - id: records
    size: 5
    type: record
    repeat: until
    repeat-until: _.marker == 0xaa
types:
  record:
    seq:
      - id: marker
        type: u1
      - id: body
        type: u4
