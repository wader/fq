meta:
  id: valid_optional_id
seq:
  - contents: 'PACK-1'
  - type: u1
    valid:
      eq: 0xff
  - type: s1
    valid:
      expr: _ == -1
