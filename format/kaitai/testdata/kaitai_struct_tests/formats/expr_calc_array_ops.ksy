meta:
  id: expr_calc_array_ops
instances:
  int_array:
    value: '[10, 25, 50, 100, 200, 500, 1000]'
  double_array:
    value: '[10.0, 25.0, 50.0, 100.0, 3.14159]'
  str_array:
    value: '["un", "deux", "trois", "quatre"]'

  int_array_size:
    value: int_array.size
  int_array_first:
    value: int_array.first
  int_array_mid:
    value: int_array[1]
  int_array_last:
    value: int_array.last
  int_array_min:
    value: int_array.min
  int_array_max:
    value: int_array.max

  double_array_size:
    value: double_array.size
  double_array_first:
    value: double_array.first
  double_array_mid:
    value: double_array[1]
  double_array_last:
    value: double_array.last
  double_array_min:
    value: double_array.min
  double_array_max:
    value: double_array.max

  str_array_size:
    value: str_array.size
  str_array_first:
    value: str_array.first
  str_array_mid:
    value: str_array[1]
  str_array_last:
    value: str_array.last
  str_array_min:
    value: str_array.min
  str_array_max:
    value: str_array.max
