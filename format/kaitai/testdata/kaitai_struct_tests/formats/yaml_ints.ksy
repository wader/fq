# Tests "forgiving" YAML KS expression parsing
meta:
  id: yaml_ints
instances:
  test_u4_dec:
    value: 4294967295
  test_u4_hex:
    value: 0xffffffff
  test_u8_dec:
    value: 18446744073709551615
  test_u8_hex:
    value: 0xffffffffffffffff
