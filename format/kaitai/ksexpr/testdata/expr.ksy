meta:
  id: test
seq:
  - id: zeros
    encoding: utf-8
    type: str
    size: 3
instances:
  ident_zeros: { value: "zeros" }

  ident_underscore: { value: "_root.zeros" }

  chain_methods: { value: "[1,2,3].size.to_s.to_i" }

  # TODO: does not work with kaitai ruby:
  # ((1+1+1) == (1+(1+1))) == true
  #
  # syntax error, unexpected == (SyntaxError)
  # ...1 + 1) + 1) == (1 + (1 + 1)) == true

  subexpr0: { value: "((1+1+1) + (1+(1+1))) > 2" }

  array_size: { value: "[1,2,3].size" }
  array_length: { value: "[1,2,3].length" } # TODO: only byte array
  # TODO: byte array to_s(encoding)
  array_max: { value: "[2,3,1].max" }
  array_min: { value: "[2,1,3].min" }
  array_first: { value: "[1,2,3].first" }
  array_last: { value: "[1,2,3].last" }

  ternary0: { value: "true ? 1 : 2" }
  ternary1: { value: "false ? 1 : 2" }
  ternary2: { value: "true ? true ? 1 : 2 : 3" }
  ternary3: { value: "false ? 1 : false ? 2 : 3" }

  boolean_to_i0: { value: "(true).to_i" }
  boolean_to_i1: { value: "(false).to_i" }

  integer_to_s: { value: "(123).to_s" }

  float_to_i: { value: "(123.45).to_i" }

  string_length0: { value: '"abc".length' }
  string_length1: { value: '"abc\u1234".length' }

  string_reverse0: { value: '"".reverse' }
  string_reverse1: { value: '"abc".reverse' }
  string_reverse2: { value: '"abcd".reverse' }
  # string_reverse3: {value: '"\u00e5\u00e4".reverse'} # TODO: use "åä", kaitai ruby issue? also see const.ksy

  string_substring0: { value: '"abcd".substring(0,4)' }
  string_substring1: { value: '"abcd".substring(1,3)' }
  string_substring2: { value: '"abcd".substring(2,2)' }
  string_substring3: { value: '"abcd".substring(3,1)' }

  string_to_i0: { value: '"123".to_i' }
  string_to_i1: { value: '"101".to_i(2)' }
  string_to_i2: { value: '"77".to_i(8)' }
  string_to_i3: { value: '"ff".to_i(16)' }
  string_to_i4: { value: '"abz".to_i(36)' }

  # TODO: enum to_i
