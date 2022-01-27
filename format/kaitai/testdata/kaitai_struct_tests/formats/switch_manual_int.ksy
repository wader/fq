meta:
  id: switch_manual_int
seq:
  - id: opcodes
    type: opcode
    repeat: eos
types:
  opcode:
    seq:
      - id: code
        type: u1
      - id: body
        type:
          switch-on: code
          cases:
            73: intval
            83: strval
    types:
      intval:
        seq:
          - id: value
            type: u1
      strval:
        seq:
          - id: value
            type: strz
            encoding: ASCII
