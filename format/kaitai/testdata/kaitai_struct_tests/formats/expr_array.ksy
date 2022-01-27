meta:
  id: expr_array
  endian: le
  encoding: UTF-8
seq:
  - id: aint
    type: u4
    repeat: expr
    repeat-expr: 4
  - id: afloat
    type: f8
    repeat: expr
    repeat-expr: 3
  - id: astr
    type: strz
    repeat: expr
    repeat-expr: 3
instances:
  aint_size:
    value: aint.size
  aint_first:
    value: aint.first
  aint_last:
    value: aint.last
  aint_min:
    value: aint.min
  aint_max:
    value: aint.max

  afloat_size:
    value: afloat.size
  afloat_first:
    value: afloat.first
  afloat_last:
    value: afloat.last
  afloat_min:
    value: afloat.min
  afloat_max:
    value: afloat.max

  astr_size:
    value: astr.size
  astr_first:
    value: astr.first
  astr_last:
    value: astr.last
  astr_min:
    value: astr.min
  astr_max:
    value: astr.max
