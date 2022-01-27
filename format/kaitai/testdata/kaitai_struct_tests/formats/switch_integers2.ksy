meta:
  id: switch_integers2
  endian: le
seq:
  - id: code
    type: u1
  - id: len
    type:
      switch-on: code
      cases:
        1: u1
        2: u2
        4: u4
        8: u8
  - id: ham
    size: len
  - id: padding
    type: u1
    if: len > 3
instances:
  len_mod_str:
    value: (len * 2 - 1).to_s
