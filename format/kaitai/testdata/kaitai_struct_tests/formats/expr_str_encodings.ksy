meta:
  id: expr_str_encodings
  endian: le
seq:
  - id: len_of_1
    type: u2
  - id: str1
    type: str
    size: len_of_1
    encoding: ASCII
  - id: len_of_2
    type: u2
  - id: str2
    type: str
    size: len_of_2
    encoding: UTF-8
  - id: len_of_3
    type: u2
  - id: str3
    type: str
    size: len_of_3
    encoding: SJIS
  - id: len_of_4
    type: u2
  - id: str4
    type: str
    size: len_of_4
    encoding: CP437
instances:
  str1_eq:
    value: str1 == "Some ASCII"
  str2_eq:
    value: str2 == "こんにちは"
  str3_eq:
    value: str3 == "こんにちは"
  str3_eq_str2:
    value: str3 == str2
  str4_eq:
    value: str4 == "░▒▓"
  str4_gt_str_calc:
    value: str4 > "┤" # in UTF-8 "░" (U+2591) > "┤" (U+2524),
                      # in CP437 "░" (0xB0)   < "┤" (0xB4)
  str4_gt_str_from_bytes:
    value: 'str4 > [0xb4].to_s("CP437")'
