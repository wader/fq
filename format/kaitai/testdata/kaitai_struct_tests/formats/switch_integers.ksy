meta:
  id: switch_integers
  endian: le
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
            1: u1
            2: u2
            4: u4
            8: u8
