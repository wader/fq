meta:
  id: valid_fail_range_str
  encoding: ASCII
seq:
  - id: foo
    size: 2
    type: str
    # there is actually [80, 65] ("PA") in the file
    valid:
      min: '"P"' # [80]
      max: '"P1"' # [80, 49]
