meta:
  id: type_ternary_opaque
  ks-opaque-types: true
seq:
  - id: dif_wo_hack
    size: 12
    type: term_strz
    if: not is_hack
  - id: dif_with_hack
    size: 12
    type: term_strz
    process: xor(0b00000011)
    if: is_hack
instances:
  is_hack:
    value: "false"
  dif:
    value: "not is_hack ? dif_wo_hack : dif_with_hack"
