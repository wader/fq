meta:
  id: eof_exception_bytes
seq:
  - id: buf
    # only 12 bytes available, should fail with EOF exception
    size: 13
