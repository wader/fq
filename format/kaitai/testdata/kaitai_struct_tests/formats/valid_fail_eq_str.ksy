meta:
  id: valid_fail_eq_str
  encoding: ASCII
seq:
  - id: foo
    size: 4
    type: str
    valid: '"BACK"' # there is actually "PACK" in the file
