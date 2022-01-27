meta:
  id: enum_import
  endian: le
  imports:
    - enum_0
    - enum_deep
seq:
  - id: pet_1
    type: u4
    enum: enum_0::animal
  - id: pet_2
    type: u4
    enum: enum_deep::container1::container2::animal
