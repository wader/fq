meta:
  id: valid_not_parsed_if
seq:
  - id: not_parsed
    if: false
    type: u1
    valid: 42
  - id: parsed
    if: true
    type: u1
    valid: 0x50
