meta:
  id: enum_deep
  endian: le
seq:
  - id: pet_1
    type: u4
    enum: container1::animal
  - id: pet_2
    type: u4
    enum: container1::container2::animal
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
