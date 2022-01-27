meta:
  id: non_standard
  endian: le
  -vendor-key: some value
  -vendor-key2:
    some: hash
    with_more: values
  -vendor-key3:
    - foo
    - bar
-vendor-key: value
seq:
  - id: foo
    type: u1
    -vendor-key: value
  - id: bar
    type:
      switch-on: foo
      -vendor-key: value
      cases:
        42: u2
        43: u4
instances:
  vi:
    value: foo
    -vendor-key: value
  pi:
    pos: 0
    type: u1
    -vendor-key: value
