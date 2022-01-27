meta:
  id: params_pass_array_usertype
seq:
  - id: blocks
    type: block
    repeat: expr
    repeat-expr: 2
  - id: pass_blocks
    type: param_type(blocks)
types:
  block:
    seq:
      - id: foo
        type: u1
  param_type:
    params:
      - id: bar
        type: block[]
    seq:
      - id: one
        size: bar[0].foo
      - id: two
        size: bar[1].foo
