meta:
  id: params_call_extra_parens
seq:
  - id: buf1
    type: my_str1((5))
types:
  my_str1:
    params:
      - id: len
        type: u4
    seq:
      - id: body
        type: str
        size: len
        encoding: UTF-8
