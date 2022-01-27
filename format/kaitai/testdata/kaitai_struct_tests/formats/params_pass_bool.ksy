meta:
  id: params_pass_bool
seq:
  - id: s_false
    type: b1
  - id: s_true
    type: b1
  - id: seq_b1
    type: param_type_b1(s_true)
  - id: seq_bool
    type: param_type_bool(s_false)
  - id: literal_b1
    type: param_type_b1(false)
  - id: literal_bool
    type: param_type_bool(true)
  - id: inst_b1
    type: param_type_b1(v_true)
  - id: inst_bool
    type: param_type_bool(v_false)
instances:
  v_false:
    value: false
  v_true:
    value: true
types:
  param_type_b1:
    params:
      - id: arg
        type: b1
    seq:
      - id: foo
        size: 'arg ? 1 : 2'
  param_type_bool:
    params:
      - id: arg
        type: bool
    seq:
      - id: foo
        size: 'arg ? 1 : 2'
