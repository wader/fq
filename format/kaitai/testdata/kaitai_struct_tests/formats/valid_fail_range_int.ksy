meta:
  id: valid_fail_range_int
seq:
  - id: foo
    type: u1
    valid:
      min: 5
      max: 10 # there is actually 0x50 in the file
