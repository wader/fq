meta:
  id: instance_repeat
  endian: be

instances:
  i0:
    pos: 1
    size: 3
    repeat: expr
    repeat-expr: 4
    type:
      switch-on: _root.i0.size
      cases:
        0: u1
        1: u2
        2: u4
