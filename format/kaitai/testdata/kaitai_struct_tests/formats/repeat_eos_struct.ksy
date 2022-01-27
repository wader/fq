meta:
  id: repeat_eos_struct
  endian: le
seq:
  - id: chunks
    type: chunk
    repeat: eos
types:
  chunk:
    seq:
      - id: offset
        type: u4
      - id: len
        type: u4
