# Test chain of absolute-into-absolute imports
#
# this -> $KS_PATH/for_abs_imports/imported_and_rel -> $KS_PATH/more_abs/imported_root
meta:
  id: imports_abs_rel
  imports:
    - /for_abs_imports/imported_and_rel
seq:
  - id: one
    type: u1
  - id: two
    type: imported_and_rel
