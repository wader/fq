meta:
  id: params_pass_array_int
  endian: le
seq:
  - id: ints
    type: u2
    repeat: expr
    repeat-expr: 3
  - id: pass_ints
    type: wants_ints(ints)
  - id: pass_ints_calc
    type: wants_ints(ints_calc)
types:
  wants_ints:
    params:
      - id: nums
        type: u2[]
instances:
  ints_calc:
    value: '[27643, 7].as<u2[]>'
