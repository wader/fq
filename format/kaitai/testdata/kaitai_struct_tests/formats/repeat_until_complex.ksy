meta:
  id: repeat_until_complex
  endian: le
seq:
  - id: first
    type: type_u1
    repeat: until
    repeat-until: _.count == 0
  - id: second
    type: type_u2
    repeat: until
    repeat-until: _.count == 0
  - id: third
    type: u1
    repeat: until
    repeat-until: _ == 0
types:
  type_u1:
   seq:
     - id: count
       type: u1
     - id: values
       type: u1
       repeat: expr
       repeat-expr: count
  type_u2:
   seq:
     - id: count
       type: u2
     - id: values
       type: u2
       repeat: expr
       repeat-expr: count
