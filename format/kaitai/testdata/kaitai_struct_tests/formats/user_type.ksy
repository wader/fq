# Tests very basic user-defined type functionality. A single type is
# defined in top-level class and is invoked via a seq attribute.
meta:
  id: user_type
  endian: le
seq:
  - id: one
    type: header
types:
  header:
    seq:
      - id: width
        type: u4
      - id: height
        type: u4
