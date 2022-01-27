# Tests string concat, comparisons and boolean instance results
meta:
  id: expr_3
seq:
  - id: one
    type: u1
  - id: two
    type: str
    encoding: ASCII
    size: 3
instances:
  three:
    value: '"@" + two'
  four:
    value: '"_" + two + "_"'
  is_str_eq:
    value: two == "ACK"
  is_str_ne:
    value: two != "ACK"
  is_str_lt:
    value: two < "ACK2"
  is_str_gt:
    value: two > "ACK2"
  is_str_le:
    value: two <= "ACK2"
  is_str_ge:
    value: two >= "ACK2"
  is_str_lt2:
    value: three < two
  test_not:
    value: not false
