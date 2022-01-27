meta:
  id: enum_to_i
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
  pet_1_i:
    value: pet_1.to_i
  pet_1_mod:
    value: pet_1.to_i + 0x8000
  one_lt_two:
    value: pet_1.to_i < pet_2.to_i
  pet_1_eq_int:
    value: pet_1.to_i == 7
  pet_2_eq_int:
    value: pet_2.to_i == 5
