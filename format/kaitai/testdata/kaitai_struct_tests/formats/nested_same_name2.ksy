meta:
  id: nested_same_name2
  endian: le
seq:
  - id: version
    type: u4
  - id: main_data
    type: main
  - id: dummy
    type: dummy_obj
types:
  main:
    seq:
      - id: main_size
        type: s4
      - id: foo
        type: foo_obj
    types:
      foo_obj:
        seq:
          - id: data1
            size: _parent.main_size * 2
  dummy_obj:
    seq:
      - id: dummy_size
        type: s4
      - id: foo
        type: foo_obj
    types:
      foo_obj:
        seq:
          - id: data2
            size: _parent.dummy_size * 2
