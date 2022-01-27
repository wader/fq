# Test if unsigned values in `terminator` work
meta:
  id: term_u1_val
seq:
  - id: foo
    terminator: 0xe3
    consume: false
  - id: bar
    type: str
    encoding: UTF-8
    terminator: 0xab
    include: true
