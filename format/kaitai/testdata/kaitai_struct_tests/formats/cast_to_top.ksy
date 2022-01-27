# Simplest way to check that casting to top-level type works
meta:
  id: cast_to_top
seq:
  - id: code
    type: u1
instances:
  header:
    pos: 1
    type: cast_to_top
  # This is silly and does nothing, but it checks that casting can find
  # top-level type
  header_casted:
    value: header.as<cast_to_top>
