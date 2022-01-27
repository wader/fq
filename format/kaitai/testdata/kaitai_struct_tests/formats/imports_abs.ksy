meta:
  id: imports_abs
  imports:
    - /common/vlq_base128_le
seq:
  - id: len
    type: vlq_base128_le
  - id: body
    size: len.value
