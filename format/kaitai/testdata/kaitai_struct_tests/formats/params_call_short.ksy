meta:
  id: params_call_short
seq:
  - id: buf1
    type: my_str1(5)
  - id: buf2
    type: my_str2(2 + 3, true)
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
  my_str2:
    params:
      - id: len
        type: u4
      - id: has_trailer
        type: bool
    seq:
      - id: body
        type: str
        size: len
        encoding: UTF-8
      - id: trailer
        type: u1
        if: has_trailer
