meta:
  id: valid_switch
seq:
  - id: a
    type: u1
    valid: 0x50
  - id: b
    type:
      switch-on: a
      cases:
        0x50: u2le
        _: u2be
    valid: 0x4341
