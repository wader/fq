meta:
  id: switch_manual_enum
seq:
  - id: opcodes
    type: opcode
    repeat: eos
types:
  opcode:
    seq:
      - id: code
        type: u1
        enum: code_enum
      - id: body
        type:
          switch-on: code
          cases:
            code_enum::intval: intval
            code_enum::strval: strval
    enums:
      code_enum:
        73: intval # 'I'
        83: strval # 'S'
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
