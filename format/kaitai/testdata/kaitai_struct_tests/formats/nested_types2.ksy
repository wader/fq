meta:
  id: nested_types2
seq:
  - id: one
    type: subtype_a
  - id: two
    type: subtype_b
types:
  subtype_a:
    seq:
      - id: typed_at_root
        type: subtype_b
      - id: typed_here1
        type: subtype_c
      - id: typed_here2
        type: subtype_cc
    types:
      subtype_c:
        seq:
          - id: value_c
            type: s1
          - id: typed_here
            type: subtype_d
          - id: typed_parent
            type: subtype_cc
          - id: typed_root
            type: subtype_b
        types:
          subtype_d:
            seq:
              - id: value_d
                type: s1
      subtype_cc:
        seq:
          - id: value_cc
            type: s1
  subtype_b:
    seq:
      - id: value_b
        type: s1
