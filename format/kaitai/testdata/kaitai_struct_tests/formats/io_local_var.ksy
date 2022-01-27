# https://github.com/kaitai-io/kaitai_struct/issues/536
meta:
  id: io_local_var
seq:
  # [0..19]
  - id: skip
    size: 20
  # Invoke `mess_up` that can potentially mess up _root._io pointer
  - id: always_null
    type: u1
    if: mess_up.as<dummy>._io.pos < 0
  # Check where we are at, should consume [20]
  - id: followup
    type: u1
instances:
  mess_up:
    io: _root._io # required to trigger `io` assignment
    pos: 8
    size: 2
    type:
      switch-on: 2
      cases:
        1: dummy
        2: dummy
types:
  dummy: {}
