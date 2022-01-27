# Enum declared in top-level type, used in subtype
meta:
  id: enum_1
  endian: le
seq:
  - id: main
    type: main_obj
types:
  main_obj:
    seq:
      - id: submain
        type: submain_obj
    types:
      submain_obj:
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
