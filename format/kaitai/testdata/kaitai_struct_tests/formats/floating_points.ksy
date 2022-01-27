meta:
  id: floating_points
  endian: le
seq:
  - id: single_value
    type: f4
  - id: double_value
    type: f8
  - id: single_value_be
    type: f4be
  - id: double_value_be
    type: f8be
  - id: approximate_value
    type: f4
instances:
  single_value_plus_int:
    value: single_value + 1
  single_value_plus_float:
    value: single_value + 0.5
  double_value_plus_float:
    value: double_value + 0.05
