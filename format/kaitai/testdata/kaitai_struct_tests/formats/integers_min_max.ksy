meta:
  id: integers_min_max
seq:
  - id: unsigned_min
    type: unsigned
  - id: unsigned_max
    type: unsigned
  - id: signed_min
    type: signed
  - id: signed_max
    type: signed
types:
  unsigned:
    seq:
      - id: u1
        type: u1
      - id: u2le
        type: u2le
      - id: u4le
        type: u4le
      - id: u8le
        type: u8le
      - id: u2be
        type: u2be
      - id: u4be
        type: u4be
      - id: u8be
        type: u8be
  signed:
    seq:
      - id: s1
        type: s1
      - id: s2le
        type: s2le
      - id: s4le
        type: s4le
      - id: s8le
        type: s8le
      - id: s2be
        type: s2be
      - id: s4be
        type: s4be
      - id: s8be
        type: s8be
