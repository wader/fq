# https://github.com/kaitai-io/kaitai_struct_compiler/issues/44
meta:
  id: opaque_external_type_02_child
  endian: le
seq:
  - id: s1
    type: str
    encoding: UTF-8
    terminator: 0x7C
  - id: s2
    type: str
    encoding: UTF-8
    terminator: 0x7C
    consume: false
  - id: s3
    type: opaque_external_type_02_child_child
types:
  opaque_external_type_02_child_child:
    seq:
      - id: s3
        type: str
        encoding: UTF-8
        terminator: 0x40
        include: true
        if: _root.some_method
instances:
  some_method:
    value: "true"
