meta:
  id: position_to_end
  endian: le
instances:
  index:
    pos: _io.size - 8
    type: index_obj
types:
  index_obj:
    seq:
     - id: foo
       type: u4
     - id: bar
       type: u4
