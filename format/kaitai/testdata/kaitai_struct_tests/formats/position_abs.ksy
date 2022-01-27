meta:
  id: position_abs
  endian: le
seq:
  - id: index_offset
    type: u4
types:
  index_obj:
    seq:
     - id: entry
       type: strz
       encoding: UTF-8
instances:
  index:
    pos: index_offset
    type: index_obj
