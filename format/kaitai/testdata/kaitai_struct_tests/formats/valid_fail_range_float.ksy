meta:
  id: valid_fail_range_float
seq:
  - id: foo
    type: f4le
    valid:
      min: 0.2
      max: 0.4 # there is actually 0.5 in the file
