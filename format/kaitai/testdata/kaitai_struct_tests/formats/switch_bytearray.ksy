meta:
  id: switch_bytearray
seq:
  - id: opcodes
    type: opcode
    repeat: eos
types:
  opcode:
    seq:
      - id: code
        size: 1
      - id: body
        type:
          switch-on: code
          cases:
            '[73]': intval
            '[83]': strval
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
