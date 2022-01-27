meta:
  id: eos_exception_u4
seq:
  - id: envelope
    type: data
    size: 6
types:
  data:
    seq:
      - id: prebuf
        size: 3
      - id: fail_int
        type: u4le
