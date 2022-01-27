meta:
  id: process_to_user
seq:
  - id: buf1
    size: 5
    process: rol(3)
    type: just_str
types:
  just_str:
    seq:
      - id: str
        type: str
        encoding: UTF-8
        size-eos: true
