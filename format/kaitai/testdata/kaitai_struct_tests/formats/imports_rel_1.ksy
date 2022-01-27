# Test chain of relative imports
#
# this -> for_rel_imports/imported_1 -> for_rel_imports/imported_2
meta:
  id: imports_rel_1
  imports:
    - for_rel_imports/imported_1
seq:
  - id: one
    type: u1
  - id: two
    type: imported_1
