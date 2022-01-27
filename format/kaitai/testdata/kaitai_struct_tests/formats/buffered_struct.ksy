meta:
  id: buffered_struct
  endian: le
seq:
  - id: len1
    type: u4
  - id: block1
    type: block
    size: len1
  - id: len2
    type: u4
  - id: block2
    type: block
    size: len2
  - id: finisher
    type: u4
types:
  block:
    seq:
      - id: number1
        type: u4
      - id: number2
        type: u4
