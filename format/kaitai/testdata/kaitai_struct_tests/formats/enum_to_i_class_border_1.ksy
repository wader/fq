# https://github.com/kaitai-io/kaitai_struct/issues/552
meta:
  id: enum_to_i_class_border_1
  endian: le
  imports:
    - enum_to_i_class_border_2

seq:
  - id:   pet_1
    type: u4
    enum: animal
  - id:   pet_2
    type: u4
    enum: animal

enums:
  animal:
    4:  dog
    7:  cat
    12: chicken

instances:
  some_dog:
    value: 4
    enum:  animal

  checker:
    pos:  0
    type: enum_to_i_class_border_2(_root)
