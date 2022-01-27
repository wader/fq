# Same as nav_parent2, but "tags" is instance and is actually positioned at ofs_tags
meta:
  id: nav_parent3
  endian: le
seq:
  - id: ofs_tags
    type: u4
  - id: num_tags
    type: u4
types:
  tag:
    seq:
      - id: name
        type: str
        size: 4
        encoding: ASCII
      - id: ofs
        type: u4
      - id: num_items
        type: u4
    types:
      tag_char:
        seq:
          - id: content
            type: str
            size: _parent.num_items
            encoding: ASCII
    instances:
      tag_content:
        pos: ofs
        type:
          switch-on: name
          cases:
            '"RAHC"': tag_char
        io: _root._io
instances:
  tags:
    pos: ofs_tags
    type: tag
    repeat: expr
    repeat-expr: num_tags
