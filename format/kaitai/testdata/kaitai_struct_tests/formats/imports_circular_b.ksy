meta:
  id: imports_circular_b
  imports:
    - imports_circular_a
seq:
  - id: initial
    type: u1
  - id: back_ref
    type: imports_circular_a
    if: initial == 0x41
