# Assert that enum id's can handle values in Int's range
meta:
  id: enum_int_range_s
  endian:  be

enums:
  constants:
    -2147483648: int_min
    0: zero
    2147483647: int_max

seq:
  - id: f1
    type: s4
    enum: constants
  - id: f2
    type: s4
    enum: constants
  - id: f3
    type: s4
    enum: constants
