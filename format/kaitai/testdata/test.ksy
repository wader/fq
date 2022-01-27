meta:
  id: test
seq:
  - id: numbers1
    type: u1

  - id: if_false
    type: u1
    if: false

  - id: repeat1
    type: u1
    repeat: expr
    repeat-expr: 3

  # - id: numbers3
  #   type: s8
  #   repeat: expr
  #   repeat-expr: 10
  # - id: test
  #   type: test

instances:
  a: { value: "1+1" }
  b: { value: numbers1 }
# types:
#   test:
#     seq:
#       - id: a
#         type: u1
