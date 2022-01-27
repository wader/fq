meta:
  id: str_encodings_default
  endian: le
  encoding: UTF-8
seq:
  - id: len_of_1
    type: u2
  - id: str1
    type: str
    size: len_of_1
  - id: rest
    type: subtype
types:
  subtype:
    seq:
      - id: len_of_2
        type: u2
      - id: str2
        type: str
        size: len_of_2
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
