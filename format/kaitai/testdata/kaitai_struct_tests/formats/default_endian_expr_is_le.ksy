meta:
  id: default_endian_expr_is_le
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
              '[0x49, 0x49]': le
              _: be
        seq:
          - id: some_int
            type: u4
          - id: some_int_be
            type: u2be
          - id: some_int_le
            type: u2le
