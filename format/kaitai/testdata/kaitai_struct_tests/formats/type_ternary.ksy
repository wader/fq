meta:
  id: type_ternary
seq:
  - id: dif_wo_hack
    size: 1
    type: dummy
    if: not is_hack
  - id: dif_with_hack
    size: 1
    type: dummy
    process: xor(0b00000011)
types:
  dummy:
    seq:
      - id: value
        type: u1
instances:
  is_hack:
    value: "true"
  dif:
    value: "not is_hack ? dif_wo_hack : dif_with_hack"
  dif_value:
    value: dif.value
