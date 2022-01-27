meta:
  id: enum_fancy
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
    4:
      id: dog
      doc: A member of genus Canis.
      -orig-id: MH_CANINE
    7:
      id: cat
      doc: Small, typically furry, carnivorous mammal.
      -orig-id: MH_FELINE
    12: chicken
