# --debug (or actually --no-auto-read) with arrays of user types requires
# special handling to avoid spoiling whole object due to exception handler.
meta:
  id: debug_array_user
  ks-debug: true
seq:
  - id: one_cat
    type: cat
  - id: array_of_cats
    type: cat
    repeat: expr
    repeat-expr: 3
types:
  cat:
    seq:
      - id: meow
        type: u1
