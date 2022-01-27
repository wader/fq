meta:
  id: docstrings
doc: One-liner description of a type.
seq:
  - id: one
    type: u1
    doc: A pretty verbose description for sequence attribute "one"
types:
  complex_subtype:
    doc: |
      This subtype is never used, yet has a very long description
      that spans multiple lines. It should be formatted accordingly,
      even in languages that have single-line comments for
      docstrings. Actually, there's even a MarkDown-style list here
      with several bullets:

      * one
      * two
      * three

      And the text continues after that. Here's a MarkDown-style link:
      [woohoo](http://example.com) - one day it will be supported as
      well.
instances:
  two:
    pos: 0
    type: u1
    doc: Another description for parse instance "two"
  three:
    value: 0x42
    doc: And yet another one for value instance "three"
