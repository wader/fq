# Assert that enum id's can handle values in Long's range
meta:
  id: enum_long_range_s
  endian:  be

enums:
  constants:
    -9223372036854775808: long_min
    -2147483649: int_below_min
    -2147483648: int_min
    0: zero
    2147483647: int_max
    2147483648: int_over_max
    9223372036854775807: long_max

seq:
  - id: f1
    type: s8
    enum: constants
  - id: f2
    type: s8
    enum: constants
  - id: f3
    type: s8
    enum: constants
  - id: f4
    type: s8
    enum: constants
  - id: f5
    type: s8
    enum: constants
  - id: f6
    type: s8
    enum: constants
  - id: f7
    type: s8
    enum: constants
