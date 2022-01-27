meta:
  id: instance_io_user
  endian: le
seq:
  - id: qty_entries
    type: u4
  - id: entries
    type: entry
    repeat: expr
    repeat-expr: qty_entries
  - id: strings
    type: strings_obj
    size-eos: true
types:
  entry:
    seq:
      - id: name_ofs
        type: u4
      - id: value
        type: u4
    instances:
      name:
        io: _root.strings._io
        pos: name_ofs
        type: strz
        encoding: UTF-8
  strings_obj:
    seq:
      - id: str
        type: strz
        encoding: UTF-8
        repeat: eos
