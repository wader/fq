meta:
  id: nested_same_name
  endian: le
seq:
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
          - id: data
            size: '_parent.main_size * 2'
  dummy_obj:
    types:
      foo: {}
