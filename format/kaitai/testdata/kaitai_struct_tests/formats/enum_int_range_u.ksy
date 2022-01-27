# Assert that enum id's can handle values in Int's range
meta:
  id: enum_int_range_u
  endian:  be

enums:
  constants:
    0: zero
    4294967295: int_max

seq:
  - id: f1
    type: u4
    enum: constants
  - id: f2
    type: u4
    enum: constants
