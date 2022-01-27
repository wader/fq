meta:
  id: bits_enum
seq:
  # byte 0
  - id: one
    type: b4
    enum: animal
  # byte 0-1
  - id: two
    type: b8
    enum: animal
  # byte 1
  - id: three
    type: b1
    enum: animal
enums:
  animal:
    0: cat
    1: dog
    4: horse
    5: platypus
