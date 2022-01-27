meta:
  id: default_endian_expr_inherited
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
          - id: insides
            type: sub_obj
        types:
          sub_obj:
            seq:
              - id: some_int
                type: u4
              - id: more
                type: subsub_obj
            types:
              subsub_obj:
                seq:
                  - id: some_int1
                    type: u2
                  - id: some_int2
                    type: u2
                instances:
                  some_inst:
                    pos: 2
                    type: u4
