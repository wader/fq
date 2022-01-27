meta:
  id: instance_std_array
  endian: le
seq:
  - id: ofs
    type: u4
  - id: entry_size
    type: u4
  - id: qty_entries
    type: u4
instances:
  entries:
    pos: ofs
    repeat: expr
    repeat-expr: qty_entries
    size: entry_size
