# Check usage of a::b::c syntax
meta:
  id: nested_types3
seq:
  - id: a_cc
    type: subtype_a::subtype_cc
  - id: a_c_d
    type: subtype_a::subtype_c::subtype_d
  - id: b
    type: subtype_b
types:
  subtype_a:
    types:
      subtype_c:
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
      - id: a_cc
        type: subtype_a::subtype_cc
      - id: a_c_d
        type: subtype_a::subtype_c::subtype_d
