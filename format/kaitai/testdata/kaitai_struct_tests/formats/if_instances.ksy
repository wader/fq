# Checks that `if` for instance protects us from doing any activity
# for instance completely, including seeks.
meta:
  id: if_instances
instances:
  never_happens:
    pos: 100500 # does not exist in the stream
    type: u1
    if: 'false' # should never happen
