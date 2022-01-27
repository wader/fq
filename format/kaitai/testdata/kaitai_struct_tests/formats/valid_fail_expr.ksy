meta:
  id: valid_fail_expr
seq:
  - id: foo
    type: u1
    valid: # should pass
      expr: _ == 1
  - id: bar
    type: s2le
    valid: # there's actually -190 in the file
      expr: _ < -190 or _ > -190
