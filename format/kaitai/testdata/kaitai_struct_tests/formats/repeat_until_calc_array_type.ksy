meta:
  id: repeat_until_calc_array_type
  endian: le
seq:
  - id: records
    size: 5
    type: record
    repeat: until
    repeat-until: _.marker == 0xaa
types:
  record:
    seq:
      - id: marker
        type: u1
      - id: body
        type: u4
instances:
  recs_accessor:
    value: records

  first_rec:
    value: recs_accessor.first
