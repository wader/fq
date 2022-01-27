# Test languages without built-in ternary operator support (Lua) if they're not blindly using `cond and if_true or if_false`
# (doesn't work properly when if_true == false)
meta:
  id: type_ternary_2nd_falsy
seq:
  - id: int_truthy
    type: u1
  - id: ut
    type: foo
  - id: int_array
    type: u1
    repeat: expr
    repeat-expr: 2
  - id: int_array_empty
    type: u1
    repeat: expr
    repeat-expr: 0
types:
  foo:
    seq:
      - id: m
        type: u1
instances:
  t:
    value: true
  null_ut:
    value: ut
    if: false
  v_false:
    value: 't ? false : true'
  v_int_zero:
    value: 't ? 0 : 10'
  v_int_neg_zero:
    value: 't ? -0 : -20'
  v_float_zero:
    value: 't ? 0.0 : 3.14'
  v_float_neg_zero:
    value: 't ? -0.0 : -2.72'
  v_str_w_zero: # "0" is falsy in PHP, might be also in other languages
    value: 't ? "0" : "30"'
  v_null_ut:
    value: 't ? null_ut : ut'
  v_str_empty:
    value: 't ? "" : "kaitai"'
  v_int_array_empty:
    value: 't ? int_array_empty : int_array'
