meta:
  id: expr_bytes_non_literal
seq:
  - id: one
    type: u1
  - id: two
    type: u1
instances:
  calc_bytes:
    value: '[one, two].as<bytes>'
