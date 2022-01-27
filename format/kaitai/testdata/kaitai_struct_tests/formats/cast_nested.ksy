meta:
  id: cast_nested
seq:
  - id: opcodes
    type: opcode
    repeat: eos
instances:
  opcodes_0_str:
    value: opcodes[0].body.as<opcode::strval>
  opcodes_0_str_value:
    value: opcodes[0].body.as<opcode::strval>.value
  opcodes_1_int:
    value: opcodes[1].body.as<opcode::intval>
  opcodes_1_int_value:
    value: opcodes[1].body.as<opcode::intval>.value
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
