meta:
  id: switch_cast
seq:
  - id: opcodes
    type: opcode
    repeat: eos
instances:
  first_obj:
    value: opcodes[0].body.as<strval>
  second_val:
    value: opcodes[1].body.as<intval>.value
  err_cast:
    value: opcodes[2].body.as<strval>
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
  intval:
    seq:
      - id: value
        type: u1
  strval:
    seq:
      - id: value
        type: strz
        encoding: ASCII
