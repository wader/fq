meta:
  id: valid_fail_min_int
seq:
  - id: foo
    type: u1
    valid:
      min: 123 # there is actually 0x50 in the file
