meta:
  id: eof_exception_u4
seq:
  # only 12 bytes available, should fail with EOF exception
  - id: prebuf
    size: 9
  - id: fail_int
    type: u4le
