meta:
  id: enum_test

seq:
  - id: tjo
    type: t

instances:
  cat:
    value: animal::cat
  cat_i:
    value: animal::cat.to_i
  cat_i1:
    value: animal::cat.to_i + 1

  pizza:
    value: t::dish::pizza
  pizza_i:
    value: t::dish::pizza.to_i
  pizza_i1:
    value: t::dish::pizza.to_i + 1

enums:
  animal:
    1: cat
    2: dog
types:
  t:
    seq:
      - id: a
        type: u1
        enum: dish
    enums:
      dish:
        1: pizza
        2: pasta
        71: aaa
