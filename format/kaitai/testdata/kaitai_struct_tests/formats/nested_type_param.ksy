meta:
  id: nested_type_param
seq:
  - id: main_seq
    type: nested::my_type(5)
types:
  nested:
    types:
      my_type:
        params:
          - id: my_len
            type: u4
        seq:
          - id: body
            type: str
            size: my_len
            encoding: ASCII
