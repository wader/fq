meta:
  id: params_def
  endian: le
params:
  - id: len
    type: u4
  - id: has_trailer
    type: bool
seq:
  - id: buf
    type: str
    size: len
    encoding: UTF-8
  - id: trailer
    type: u1
    if: has_trailer
