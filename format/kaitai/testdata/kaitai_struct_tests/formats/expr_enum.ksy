# Tests enum for value instances
meta:
  id: expr_enum
seq:
  - id: one
    type: u1
instances:
  const_dog:
    value: 4
    enum: animal
  derived_boom:
    value: one
    enum: animal
  derived_dog:
    value: one - 98
    enum: animal
enums:
  animal:
    4: dog
    7: cat
    12: chicken
    0x66: boom
