# Checks coercion of a switch type (with and without processing)
meta:
  id: process_coerce_switch
seq:
  - id: buf_type
    type: u1
  - id: flag
    type: u1
  - id: buf_unproc
    size: 4
    if: flag == 0
    type:
      switch-on: buf_type
      cases:
        0: foo
  - id: buf_proc
    size: 4
    process: xor(0xaa)
    if: flag != 0
    type:
      switch-on: buf_type
      cases:
        0: foo
instances:
  buf:
    value: 'flag == 0 ? buf_unproc : buf_proc'
types:
  foo:
    seq:
      - id: bar
        size: 4
