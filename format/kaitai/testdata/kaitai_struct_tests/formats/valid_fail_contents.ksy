meta:
  id: valid_fail_contents
seq:
  - id: foo
    contents: [0x51, 0x41] # there is actually [0x50, 0x41] in the file
