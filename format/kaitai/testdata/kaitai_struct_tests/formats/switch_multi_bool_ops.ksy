# https://github.com/kaitai-io/kaitai_struct_compiler/issues/39
meta:
  id: switch_multi_bool_ops
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
          switch-on: "(  (code >   0) and
                         (code <=  8) and
                        ((code != 10) ? true : false)) ? code : 0"
          cases:
            1: u1
            2: u2
            4: u4
            8: u8
