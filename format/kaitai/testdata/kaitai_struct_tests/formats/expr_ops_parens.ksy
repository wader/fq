meta:
  id: expr_ops_parens
instances:
  i_42:
    value: 42
  i_m13:
    value: -13
  i_sum_to_str:
    value: (i_42 + i_m13).to_s

  f_2pi:
    value: 6.28
  f_e:
    value: 2.72
  f_sum_to_int:
    value: (f_2pi + f_e).to_i

  str_0_to_4:
    value: '"01234"'
  str_5_to_9:
    value: '"56789"'
  str_concat_len:
    value: (str_0_to_4 + str_5_to_9).length
  str_concat_rev:
    value: (str_0_to_4 + str_5_to_9).reverse
  str_concat_substr_2_to_7:
    value: (str_0_to_4 + str_5_to_9).substring(2, 7)
  str_concat_to_i:
    value: (str_0_to_4 + str_5_to_9).to_i

  bool_eq:
    value: (false == true).to_i
  bool_and:
    value: (false and true).to_i
  bool_or:
    value: (not false or false).to_i
