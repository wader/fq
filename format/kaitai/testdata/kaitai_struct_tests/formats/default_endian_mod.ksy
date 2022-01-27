# Test proper propagation of default endianness in hierarchies
meta:
  id: default_endian_mod
seq:
  - id: main
    type: main_obj
types:
  main_obj:
    meta:
      endian: le
    seq:
      - id: one
        type: s4
      - id: nest
        type: subnest
      - id: nest_be
        type: subnest_be
    types:
      subnest:
        seq:
          - id: two
            type: s4
      subnest_be:
        meta:
          endian: be
        seq:
          - id: two
            type: s4
