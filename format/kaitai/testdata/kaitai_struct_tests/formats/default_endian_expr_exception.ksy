# Parses doc[0] and doc[1], then raises an exception on doc[2] due to
# unknown endianness
meta:
  id: default_endian_expr_exception
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
              '[0x4d, 0x4d]': be
        seq:
          - id: some_int
            type: u4
          - id: some_int_be
            type: u2be
          - id: some_int_le
            type: u2le
