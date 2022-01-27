meta:
  id: enum_deep_literals
  endian: le
seq:
  - id: pet_1
    type: u4
    enum: container1::animal
  - id: pet_2
    type: u4
    enum: container1::container2::animal
instances:
  is_pet_1_ok:
    value: pet_1 == container1::animal::cat
  is_pet_2_ok:
    value: pet_2 == container1::container2::animal::hare
types:
  container1:
    enums:
      animal:
        4: dog
        7: cat
        12: chicken
    types:
      container2:
        enums:
          animal:
            4: canary
            7: turtle
            12: hare
