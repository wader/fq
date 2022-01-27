# Test proper propagation of default bit endianness in hierarchies
meta:
  id: default_bit_endian_mod
seq:
  - id: main
    type: main_obj
types:
  main_obj:
    meta:
      bit-endian: le
    seq:
      - id: one
        type: b9
      - id: two
        type: b15
      - id: nest
        type: subnest
      - id: nest_be
        type: subnest_be
    types:
      subnest:
        seq:
          - id: two
            type: b16
      subnest_be:
        meta:
          bit-endian: be
        seq:
          - id: two
            type: b16
