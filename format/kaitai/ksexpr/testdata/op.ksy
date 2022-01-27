meta:
  id: test
instances:
  inv_int0: { value: "~10" }
  inv_int1: { value: "~0xffff_ffff" }
  inv_int2: { value: "~(-10)" }
  neg_int0: { value: "-10" }
  neg_int1: { value: "-0xffff_ffff" }
  neg_int2: { value: "-(-10)" }
  not_bool0: { value: "not true" }
  not_bool1: { value: "not not true" }

  add_int: { value: "4 + 2" }
  add_float: { value: "4.0 + 2" }
  add_float2: { value: "4.0 + 3.0" }
  add_string0: { value: '"" + "a"' }
  add_string1: { value: '"ab" + "cd"' }

  sub_int: { value: "4 - 2" }
  sub_float: { value: "4.0 - 2" }
  sub_float2: { value: "4.0 - 3.0" }
  sub_float3: { value: "4.0 - 3.2" }

  mul_int: { value: "4 * 2" }
  mul_float: { value: "4.0 * 2" }
  mul_float2: { value: "4.0 * 3.0" }
  mul_float3: { value: "4.0 * 3.2" }

  div_int: { value: "4 / 2" }
  div_int2: { value: "4 / 3" }
  div_float: { value: "4.0 / 2" }
  div_float2: { value: "4.0 / 3.0" }
  div_float3: { value: "4.0 / 3.2" }

  mod_int: { value: "4 % 2" }
  mod_float: { value: "4.0 % 2" }
  mod_float2: { value: "4.0 % 3.0" }
  mod_float3: { value: "4.0 % 3.2" }

  lt_int: { value: "4 < 2" }
  lt_float: { value: "4.0 < 2" }
  lt_float2: { value: "4.0 < 3.0" }
  lt_float3: { value: "4.0 < 3.2" }
  lt_string0: { value: '"a" < "a"' }
  lt_string1: { value: '"a" < "b"' }
  lt_string2: { value: '"b" < "a"' }

  lteq_int: { value: "4 <= 2" }
  lteq_float: { value: "4.0 <= 2" }
  lteq_float2: { value: "4.0 <= 3.0" }
  lteq_float3: { value: "4.0 <= 3.2" }
  lteq_string0: { value: '"a" <= "a"' }
  lteq_string1: { value: '"a" <= "b"' }
  lteq_string2: { value: '"b" <= "a"' }

  gt_int: { value: "4 > 2" }
  gt_float: { value: "4.0 > 2" }
  gt_float2: { value: "4.0 > 3.0" }
  gt_float3: { value: "4.0 > 3.2" }
  gt_string0: { value: '"a" > "a"' }
  gt_string1: { value: '"a" > "b"' }
  gt_string2: { value: '"b" > "a"' }

  gteq_int: { value: "4 >= 2" }
  gteq_float: { value: "4.0 >= 2" }
  gteq_float2: { value: "4.0 >= 3.0" }
  gteq_float3: { value: "4.0 >= 3.2" }
  gteq_string0: { value: '"a" >= "a"' }
  gteq_string1: { value: '"a" >= "b"' }
  gteq_string2: { value: '"b" >= "a"' }

  eq_bool0: { value: "true == true" }
  eq_bool1: { value: "true == false" }
  eq_int0: { value: "2 == 2" }
  eq_int1: { value: "4 == 2" }
  eq_float0: { value: "2.0 == 2" }
  eq_float1: { value: "4.0 == 2.0" }
  eq_float2: { value: "2.0 == 2.0" }
  eq_float3: { value: "3.2 == 3.2" }
  eq_string0: { value: '"a" == "a"' }
  eq_string1: { value: '"a" == "b"' }
  eq_array0: { value: "[1,2,3] == [1]" }
  eq_array1: { value: "[1,2,3] == [1,2,3]" }

  noteq_bool0: { value: "true != true" }
  noteq_bool1: { value: "true != false" }
  noteq_int0: { value: "2 != 2" }
  noteq_int1: { value: "4 != 2" }
  noteq_float0: { value: "2.0 != 2" }
  noteq_float1: { value: "4.0 != 2.0" }
  noteq_float2: { value: "2.0 != 2.0" }
  noteq_float3: { value: "3.2 != 3.2" }
  noteq_string0: { value: '"a" != "a"' }
  noteq_string1: { value: '"a" != "b"' }
  noteq_array0: { value: "[1,2,3] != [1]" }
  noteq_array1: { value: "[1,2,3] != [1,2,3]" }

  bsl_int: { value: "8 << 2" }
  #bsl_float0:  {value: '8.0 << 2'}
  #bsl_float1:  {value: '8.0 << 2.0'}

  bsr_int: { value: "8 >> 2" }
  #bsr_float0:  {value: '8.0 >> 2'}
  #bsr_float1:  {value: '8.0 >> 2.0'}

  band_int: { value: "15 & 255" }
  #band_float0: {value: '15.0 & 255'}
  #band_float1: {value: '15.0 & 255.0'}

  bor_int: { value: "15 | 240" }
  #bor_float0: {value: '15.0 | 240'}
  #bor_float1: {value: '15.0 | 240.0'}

  bxor_int: { value: "31 ^ 248" }
  #bxor_float0: {value: '31.0 ^ 248'}
  #bxor_float1: {value: '31.0 ^ 248.0'}

  band0: { value: "false and false" }
  band1: { value: "false and true" }
  band2: { value: "true and false" }
  band3: { value: "true and true" }

  bor0: { value: "false or false" }
  bor1: { value: "false or true" }
  bor2: { value: "true or false" }
  bor3: { value: "true or true" }

  # associativity
  bor_assoc0: { value: '"a" == "b" or false' }

  band_assoc0: { value: "123 & 0b111_00000 == 0b101_00000" }
