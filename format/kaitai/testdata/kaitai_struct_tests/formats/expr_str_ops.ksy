meta:
  id: expr_str_ops
  encoding: ASCII
seq:
  - id: one
    type: str
    size: 5
instances:
  one_len:
    value: one.length
  one_rev:
    value: one.reverse
  one_substr_0_to_3:
    value: one.substring(0, 3)
  one_substr_2_to_5:
    value: one.substring(2, 5)
  one_substr_3_to_3:
    value: one.substring(3, 3)

  two:
    value: '"0123456789"'
  two_len:
    value: two.length
  two_rev:
    value: two.reverse
  two_substr_0_to_7:
    value: two.substring(0, 7)
  two_substr_4_to_10:
    value: two.substring(4, 10)
  two_substr_0_to_10:
    value: two.substring(0, 10)

  to_i_attr:
    value: '"9173".to_i'
  to_i_r10:
    value: '"-072".to_i(10)'
  to_i_r2:
    value: '"1010110".to_i(2)'
  to_i_r8:
    value: '"721".to_i(8)'
  to_i_r16:
    value: '"47cf".to_i(16)'
