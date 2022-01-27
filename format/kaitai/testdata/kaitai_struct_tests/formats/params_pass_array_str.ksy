meta:
  id: params_pass_array_str
  encoding: ascii
seq:
  - id: str_array
    type: str
    size: 2
    repeat: expr
    repeat-expr: 3
  - id: pass_str_array
    type: wants_strs(str_array)
  - id: pass_str_array_calc
    type: wants_strs(str_array_calc)
types:
  wants_strs:
    params:
      - id: strs
        type: str[]
instances:
  str_array_calc:
    value: '["aB", "Cd"]'
