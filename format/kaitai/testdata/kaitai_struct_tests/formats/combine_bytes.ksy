meta:
  id: combine_bytes
seq:
  - id: bytes_term
    terminator: 0x7c
  - id: bytes_limit
    size: 4
  - id: bytes_eos
    size-eos: true
instances:
  bytes_calc:
    value: '[0x52, 0x6e, 0x44]'

  term_or_limit:
    value: 'true ? bytes_term : bytes_limit'
  term_or_eos:
    value: 'false ? bytes_term : bytes_eos'
  term_or_calc:
    value: 'true ? bytes_term : bytes_calc'

  limit_or_eos:
    value: 'true ? bytes_limit : bytes_eos'
  limit_or_calc:
    value: 'false ? bytes_limit : bytes_calc'

  eos_or_calc:
    value: 'true ? bytes_eos : bytes_calc'
