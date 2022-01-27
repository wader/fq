meta:
  id: cast_to_imported
  imports:
    - hello_world
seq:
  - id: one
    type: hello_world
instances:
  one_casted:
    value: one.as<hello_world>
