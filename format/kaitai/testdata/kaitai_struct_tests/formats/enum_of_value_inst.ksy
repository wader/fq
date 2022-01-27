meta:
  id: enum_of_value_inst
  endian: le
seq:
  - id: pet_1
    type: u4
    enum: animal
  - id: pet_2
    type: u4
    enum: animal
enums:
  animal:
    4: dog
    7: cat
    12: chicken
instances:
  pet_3:
    value: "pet_1 == animal::cat ? 4 : 12"
    enum: animal
  pet_4:
    value: "pet_1 == animal::cat ? animal::dog : animal::chicken"
