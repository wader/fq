meta:
  id: valid_fail_inst
seq:
  - id: a
    type: u1
    if: inst >= 0 # invoke instance
instances:
  inst:
    pos: 5
    type: u1
    valid: 0x50 # there is actually 0x31 in the file
