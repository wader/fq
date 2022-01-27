# Checks coercion of two different byte arrays (with and without processing)
meta:
  id: process_coerce_bytes
seq:
  - id: records
    type: record
    repeat: expr
    repeat-expr: 2
types:
  record:
    seq:
      - id: flag
        type: u1
      - id: buf_unproc
        size: 4
        if: flag == 0
      - id: buf_proc
        size: 4
        process: xor(0xaa)
        if: flag != 0
    instances:
      buf:
        value: 'flag == 0 ? buf_unproc : buf_proc'
