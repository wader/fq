meta:
  id: expr_2
  endian: le
seq:
  - id: str1
    type: mod_str
  - id: str2
    type: mod_str
types:
  mod_str:
    seq:
      - id: len_orig
        type: u2
      - id: str
        type: str
        size: len_mod
        encoding: UTF-8
      - id: rest
        type: tuple
        size: 3
    instances:
      len_mod:
        value: len_orig - 3
      char5:
        pos: 5
        type: str
        size: 1
        encoding: ASCII
      tuple5:
        pos: 5
        type: tuple
  tuple:
    seq:
      - id: byte0
        type: u1
      - id: byte1
        type: u1
      - id: byte2
        type: u1
    instances:
      avg:
        value: (byte1 + byte2) / 2
instances:
  str1_len:
    value: str1.str.length
  str1_len_mod:
    value: str1.len_mod
  str1_byte1:
    value: str1.rest.byte1
  str1_avg:
    value: str1.rest.avg
  str1_char5:
    value: str1.char5
  str1_tuple5:
    value: str1.tuple5
  str2_tuple5:
    value: str2.tuple5
