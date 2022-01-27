meta:
  id: multiple_use
  endian: le
seq:
  - id: t1
    type: type_1
  - id: t2
    type: type_2
types:
  multi:
    seq:
      - id: value
        type: s4
  type_1:
    seq:
      - id: first_use
        type: multi  # parent type = type_1
  type_2:
    instances:
      second_use:
        pos: 0
        type: multi  # parent type = type_2
