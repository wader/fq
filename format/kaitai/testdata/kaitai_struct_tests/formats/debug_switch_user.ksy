meta:
  id: debug_switch_user
  endian: le
  ks-debug: true
seq:
  - id: code
    type: u1
  - id: data
    type:
      switch-on: code
      cases:
        1: one
        2: two
types:
  one:
    seq:
      - id: val
        type: s2
  two:
    seq:
      - id: val
        type: u2
