meta:
  id: nav_parent_false
seq:
  - id: child_size
    type: u1
  - id: element_a
    type: parent_a
  - id: element_b
    type: parent_b
types:
  parent_a:
    seq:
      - id: foo
        type: child
      - id: bar
        type: parent_b
  parent_b:
    seq:
      - id: foo
        type: child
        parent: false
  child:
    # should have only one parent of type `parent_a`
    seq:
      - id: code
        type: u1
      - id: more
        size: _parent._parent.child_size
        if: code == 73
