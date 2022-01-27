# https://github.com/kaitai-io/kaitai_struct/issues/552
meta:
  id: enum_to_i_class_border_2
  endian: le
  imports:
    - enum_to_i_class_border_1

params:
  - id:   parent
    type: enum_to_i_class_border_1

instances:
  is_dog:
    value: parent.some_dog.to_i == 4
