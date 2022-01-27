meta:
  id: valid_fail_eq_bytes
seq:
  - id: foo
    size: 2
    valid: '[0x51, 0x41]' # there is actually [0x50, 0x41] in the file
