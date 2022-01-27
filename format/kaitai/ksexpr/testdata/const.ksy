meta:
  id: test
instances:
  true: { value: "true" }
  false: { value: "false" }
  int0: { value: "0" }
  int1: { value: "123" }
  int2: { value: "-123" }
  bin0: { value: "0b0" }
  bin1: { value: "0b101" }
  bin2: { value: "-0b101" }
  oct0: { value: "0o0" }
  oct1: { value: "0o77" }
  oct2: { value: "-0o77" }
  hex0: { value: "0x0" }
  hex1: { value: "0xff" }
  hex2: { value: "-0xff" }
  bigint: { value: "0xff_ff_ff_ff_ff_ff_ff_ff" }
  float0: { value: "123.0" }
  float1: { value: "123.456" }
  float2: { value: "1.234e2" }
  # float2:  {value: '1.23e2'} # TODO: ends up as 123 integer != 123 float
  string0: { value: '""' }
  string1: { value: '"abc"' }
  # string2: {value: '"åäö"'} # TODO: ksdump says invalid utf-8
  string3: { value: '"a\\uf00fb"' }
  string4: { value: '"a\\u0020b"' }
  string5: { value: '"a\\24b"' }
  string6: { value: '"abc\\a\\b\\t\\n\\v\\f\\r\\e\\\"\\\\"' } # TODO: more escape, \', \uffff \123
  string7: { value: "''" }
  string8: { value: "'abc'" }
  # string8: {value: "'åäö'"}
  # note one level of \\ here is because of yaml double quote
  string9:
    { value: "'abc\\\\a\\\\b\\\\t\\\\n\\\\v\\\\f\\\\r\\\\e\\\\\"\\\\\\\\'" }
  array0: { value: "[1]" } # TODO: test empty?
  array1: { value: "[1,2,3]" } # byte array
  array3: { value: "[1000,1002,1003]" } # array
