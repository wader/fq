meta:
  id: params_pass_struct
seq:
  - id: first
    type: block
  - id: one
    type: struct_type(first)
types:
  block:
    seq:
      - id: foo
        type: u1
  struct_type:
    params:
      - id: foo
        type: struct
    seq:
      - id: bar
        type: baz(foo)
    types:
      baz:
        params:
          - id: foo
            type: struct
        seq:
          - id: qux
            type: u1
