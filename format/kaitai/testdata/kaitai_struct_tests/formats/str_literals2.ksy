# Tests double-quoted interpolation variables used in various
# languages. None of this should be actually processed as interpolated
# variables, all should be escaped properly in relevant translators.
meta:
  id: str_literals2
instances:
  # PHP, Perl
  dollar1:
    value: '"$foo"'
  # PHP, Perl
  dollar2:
    value: '"${foo}"'
  # Ruby
  hash:
    value: '"#{foo}"'
  # Perl
  at_sign:
    value: '"@foo"'
