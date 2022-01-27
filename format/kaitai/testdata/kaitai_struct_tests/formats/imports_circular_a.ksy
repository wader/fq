meta:
  id: imports_circular_a
  imports:
    - imports_circular_b
seq:
  - id: code
    type: u1
  - id: two
    type: imports_circular_b
