meta:
  id: valid_fail_range_bytes
seq:
  - id: foo
    size: 2
    valid:
      min: '[80]'
      max: '[80, 49]' # there is actually [80, 65] in the file
