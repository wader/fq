meta:
  id: params_pass_array_struct
seq:
  - id: one
    type: foo
  - id: two
    type: bar
  - id: pass_structs
    type: struct_type(one_two)
instances:
  one_two:
    value: '[one, two]'
types:
  foo:
    seq:
      - id: f
        type: u1
  bar:
    seq:
      - id: b
        type: u1
  struct_type:
    params:
      - id: structs
        type: struct[]
