meta:
  id: repeat_until
seq:
  - id: count
    type: u1
  - id: numbers
    type: u1
    repeat: until
    repeat-until: _ == 2
