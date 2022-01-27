meta:
  id: default_endian_expr_is_be
seq:
  - id: docs
    repeat: eos
    type: doc
types:
  doc:
    seq:
      - id: indicator
        size: 2
      - id: main
        type: main_obj
    types:
      main_obj:
        meta:
          endian:
            switch-on: _parent.indicator
            cases:
              '[0x4d, 0x4d]': be
              _: le
        seq:
          - id: some_int
            type: u4
          - id: some_int_be
            type: u2be
          - id: some_int_le
            type: u2le
        instances:
          inst_int:
            pos: 2
            type: u4
          inst_sub:
            pos: 2
            type: sub_main_obj
        types:
          sub_main_obj:
            seq:
              - id: foo
                type: u4
