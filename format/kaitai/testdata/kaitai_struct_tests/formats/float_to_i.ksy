meta:
  id: float_to_i
  endian: le
seq:
  - id: single_value
    type: f4
  - id: double_value
    type: f8
instances:
  calc_float1:
    value: 1.234
  calc_float2:
    value: 1.5
  calc_float3:
    value: 1.9
  calc_float4:
    value: -2.7
  single_i:
    value: single_value.to_i
  double_i:
    value: double_value.to_i
  float1_i:
    value: calc_float1.to_i
  float2_i:
    value: calc_float2.to_i
  float3_i:
    value: calc_float3.to_i
  float4_i:
    value: calc_float4.to_i
