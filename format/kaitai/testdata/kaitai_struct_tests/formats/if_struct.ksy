meta:
  id: if_struct
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
      - id: arg_tuple
        type: arg_tuple
        if: opcode == 0x54 # "T"
      - id: arg_str
        type: arg_str
        if: opcode == 0x53 # "S"
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
