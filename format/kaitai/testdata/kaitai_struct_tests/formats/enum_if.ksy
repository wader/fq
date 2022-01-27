meta:
  id: enum_if
  endian: le
seq:
  - id: op1
    type: operation
  - id: op2
    type: operation
  - id: op3
    type: operation
types:
  operation:
    seq:
      - id: opcode
        type: u1
        enum: opcodes
      - id: arg_tuple
        type: arg_tuple
        if: opcode == opcodes::a_tuple
      - id: arg_str
        type: arg_str
        if: opcode == opcodes::a_string
  arg_tuple:
    seq:
      - id: num1
        type: u1
      - id: num2
        type: u1
  arg_str:
    seq:
      - id: len
        type: u1
      - id: str
        type: str
        size: len
        encoding: UTF-8
enums:
  opcodes:
    0x53: a_string
    0x54: a_tuple
