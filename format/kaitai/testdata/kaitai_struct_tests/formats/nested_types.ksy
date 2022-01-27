meta:
  id: nested_types
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
      - id: typed_here
        type: subtype_c
    types:
      subtype_c:
        seq:
          - id: value_c
            type: s1
  subtype_b:
    seq:
      - id: value_b
        type: s1
