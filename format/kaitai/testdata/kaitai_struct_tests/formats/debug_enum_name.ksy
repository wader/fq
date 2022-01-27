# Tests enum name generation in --debug mode
meta:
  id: debug_enum_name
  ks-debug: true
seq:
  - id: one
    type: u1
    enum: test_enum1
  - id: array_of_ints
    type: u1
    enum: test_enum2
    repeat: expr
    repeat-expr: 1
  - id: test_type
    type: test_subtype
types:
  test_subtype:
    seq:
      - id: field1
        type: u1
        enum: inner_enum1
      - id: field2
        type: u1
    instances:
      instance_field:
        value: "field2 & 0xf"
        enum: inner_enum2
    enums:
      inner_enum1:
        67: enum_value_67
      inner_enum2:
        11: enum_value_11
enums:
  test_enum1:
    80: enum_value_80
  test_enum2:
    65: enum_value_65
