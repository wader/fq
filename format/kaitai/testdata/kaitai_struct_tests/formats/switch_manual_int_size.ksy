meta:
  id: switch_manual_int_size
  endian: le
seq:
  - id: chunks
    type: chunk
    repeat: eos
types:
  chunk:
    seq:
      - id: code
        type: u1
      - id: size
        type: u4
      - id: body
        size: size
        type:
          switch-on: code
          cases:
            0x11: chunk_meta
            0x22: chunk_dir
    types:
      chunk_meta:
        seq:
          - id: title
            type: strz
            encoding: UTF-8
          - id: author
            type: strz
            encoding: UTF-8
      chunk_dir:
        seq:
          - id: entries
            type: str
            size: 4
            repeat: eos
            encoding: UTF-8
