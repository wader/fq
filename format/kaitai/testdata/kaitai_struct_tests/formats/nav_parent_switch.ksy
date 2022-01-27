meta:
  id: nav_parent_switch
seq:
  - id: category
    type: u1
  - id: content
    type:
      switch-on: category
      cases:
        1: element_1
types:
  element_1:
    seq:
      - id: foo
        type: u1
      - id: subelement
        type: subelement_1
  subelement_1:
    seq:
      - id: bar
        type: u1
        if: _parent.foo == 0x42
