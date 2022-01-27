meta:
  id: test
  endian: le
seq:
  - id: numbers
    type: u1
    repeat: until
    repeat-until: _ == 1
  - id: some_type
    type: bla
  - id: some_str
    type: str
    size: 10
instances:
  inst1:
    value: 1+122
  a:
    value: '"aaa"'
  b:
    value: inst1+123
  b2:
    value: inst1+123

types:
  bla:
    seq:
      - id: test
        type: u2
