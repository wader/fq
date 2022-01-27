meta:
  id: combine_bool
seq:
  - id: bool_bit
    type: b1
instances:
  bool_calc:
    value: false
  bool_calc_bit:
    value: 'true ? bool_calc : bool_bit'
