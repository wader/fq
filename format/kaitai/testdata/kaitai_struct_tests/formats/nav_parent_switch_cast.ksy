meta:
  id: nav_parent_switch_cast
seq:
  - id: main
    type: foo
types:
  foo:
    seq:
      - id: buf_type
        type: u1
      - id: flag
        type: u1
      - id: buf
        size: 4
        type:
          switch-on: buf_type
          cases:
            0: zero
            1: one
    types:
      zero:
        seq:
        - id: branch
          type: common
      one:
        seq:
        - id: branch
          type: common
      common:
        instances:
          flag:
            value: '_parent._parent.as<foo>.flag'
