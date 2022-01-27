meta:
  id: instance_user_array
  endian: le
seq:
  - id: ofs
    type: u4
  - id: entry_size
    type: u4
  - id: qty_entries
    type: u4
types:
  entry:
    seq:
      - id: word1
        type: u2
      - id: word2
        type: u2
instances:
  user_entries:
    pos: ofs
    repeat: expr
    repeat-expr: qty_entries
    size: entry_size
    type: entry
    if: ofs > 0
