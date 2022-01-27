meta:
  id: combine_str
  encoding: ASCII
seq:
  - id: str_term
    type: str
    terminator: 0x7c
  - id: str_limit
    type: str
    size: 4
  - id: str_eos
    type: str
    size-eos: true
instances:
  str_calc:
    value: '"bar"'
  calc_bytes:
    value: '[0x62, 0x61, 0x7a]'
  str_calc_bytes:
    value: 'calc_bytes.to_s("ASCII")'

  term_or_limit:
    value: 'true ? str_term : str_limit'
  term_or_eos:
    value: 'false ? str_term : str_eos'
  term_or_calc:
    value: 'true ? str_term : str_calc'
  term_or_calc_bytes:
    value: 'false ? str_term : str_calc_bytes'

  limit_or_eos:
    value: 'true ? str_limit : str_eos'
  limit_or_calc:
    value: 'false ? str_limit : str_calc'
  limit_or_calc_bytes:
    value: 'true ? str_limit : str_calc_bytes'

  eos_or_calc:
    value: 'false ? str_eos : str_calc'
  eos_or_calc_bytes:
    value: 'true ? str_eos : str_calc_bytes'

  calc_or_calc_bytes:
    value: 'false ? str_calc : str_calc_bytes'
