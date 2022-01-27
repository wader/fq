meta:
  id: switch_manual_str
seq:
  - id: opcodes
    type: opcode
    repeat: eos
types:
  opcode:
    seq:
      - id: code
        type: str
        size: 1
        encoding: ASCII
      - id: body
        type:
          switch-on: code
          cases:
            '"I"': intval
            '"S"': strval
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
