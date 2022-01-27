meta:
  id: nav_parent_override
seq:
  - id: child_size
    type: u1
  - id: child_1
    type: child
  - id: mediator_2
    type: mediator
types:
  mediator:
    seq:
      - id: child_2
        type: child
        parent: _parent
  child:
    seq:
      - id: data
        size: _parent.child_size
