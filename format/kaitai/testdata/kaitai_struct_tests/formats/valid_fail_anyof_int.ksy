meta:
  id: valid_fail_anyof_int
seq:
  - id: foo
    type: u1
    valid: # there is actually 0x50 in the file
      any-of:
        - 5
        - 6
        - 7
        - 8
        - 10
        - 11
        - 12
        - 47
