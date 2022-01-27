meta:
  id: params_pass_usertype
seq:
  - id: first
    type: block
  - id: one
    type: param_type(first)
types:
  block:
    seq:
      - id: foo
        type: u1
  param_type:
    params:
      - id: foo
        type: block
    seq:
      - id: buf
        size: foo.foo
