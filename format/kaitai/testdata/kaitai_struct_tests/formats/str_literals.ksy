meta:
  id: str_literals
instances:
  complex_str:
    value: '"\0\1\2\a\b\n\r\t\v\f\e\75\7\u000a\u0024\u263b"'
  double_quotes:
    value: '"\"\u0022\42"'
  backslashes:
    value: '"\\\u005c\134"'
  octal_eatup:
    value: '"\0\62\62"'
  octal_eatup2:
    value: '"\2\62"'
