meta:
  id: docstrings_docref
doc: Another one-liner
doc-ref: http://www.example.com/some/path/?even_with=query&more=2
seq:
  - id: one
    type: u1
    doc-ref: Plain text description of doc ref, page 42
  - id: two
    type: u1
    doc: Both doc and doc-ref are defined
    doc-ref: http://www.example.com/with/url/again
  - id: three
    type: u1
    doc-ref: http://www.example.com/three Documentation name
instances:
  foo:
    doc-ref: Doc ref for instance, a plain one
    value: true
  parse_inst:
    pos: 0
    type: u1
    doc-ref: |
      Now this is a really
      long document ref that
      spans multiple lines.
