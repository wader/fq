meta:
  id: switch_manual_str_else
  endian: le
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
            _: noneval
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
      noneval:
        seq:
          - id: filler
            type: u4
