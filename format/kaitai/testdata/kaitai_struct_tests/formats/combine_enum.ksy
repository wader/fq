meta:
  id: combine_enum
  endian: le
seq:
  - id: enum_u4
    type: u4
    enum: animal
  - id: enum_u2
    type: u2
    enum: animal
enums:
  animal:
    7: pig
    12: horse
instances:
  enum_u4_u2:
    value: 'false ? enum_u4 : enum_u2'
