meta:
  id: recursive_one
seq:
  - id: one
    type: u1
  - id: next
    type:
      switch-on: one & 3
      cases:
        0: recursive_one
        1: recursive_one
        2: recursive_one
        3: fini
types:
  fini:
    seq:
      - id: finisher
        type: u2le
