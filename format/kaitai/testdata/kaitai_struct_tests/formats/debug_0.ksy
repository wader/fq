# Tests various capabilities in --debug mode
meta:
  id: debug_0
  ks-debug: true
seq:
  - id: one
    type: u1
  - id: array_of_ints
    type: u1
    repeat: expr
    repeat-expr: 3
  - type: u1 # anonymous, numbered field
