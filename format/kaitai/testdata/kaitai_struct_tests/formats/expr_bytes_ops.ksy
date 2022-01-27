meta:
  id: expr_bytes_ops
seq:
  - id: one
    size: 3
instances:
  two:
    value: '[0x41, 0xff, 0x4b]'

  one_size:
    value: one.size
  one_first:
    value: one.first
  one_mid:
    value: one[1]
  one_last:
    value: one.last
  one_min:
    value: one.min
  one_max:
    value: one.max

  two_size:
    value: two.size
  two_first:
    value: two.first
  two_mid:
    value: two[1]
  two_last:
    value: two.last
  two_min:
    value: two.min
  two_max:
    value: two.max
