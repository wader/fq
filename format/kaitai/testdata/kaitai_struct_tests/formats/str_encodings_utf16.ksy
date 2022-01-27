meta:
  id: str_encodings_utf16
  endian: le
seq:
  - id: len_be
    type: u4
  - id: be_bom_removed
    size: len_be
    type: str_be_bom_removed

  - id: len_le
    type: u4
  - id: le_bom_removed
    size: len_le
    type: str_le_bom_removed

# instances:
#   be_with_bom:
#     pos: 4
#     size: len_be
#     type: str
#     encoding: UTF-16

#   le_with_bom:
#     pos: 4 + len_be + 4
#     size: len_le
#     type: str
#     encoding: UTF-16

types:
  str_be_bom_removed:
    seq:
      - id: bom
        type: u2be
      - id: str
        size-eos: true
        type: str
        encoding: UTF-16BE

  str_le_bom_removed:
    seq:
      - id: bom
        type: u2le
      - id: str
        size-eos: true
        type: str
        encoding: UTF-16LE
