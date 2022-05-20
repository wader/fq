include "internal";
include "binary";
include "decode";

def grep_by(f):
  ( ..
  | select(f)?
  );

def _value_grep_string_cond($v; $flags):
  if _is_string then test($v; $flags)
  else false
  end;

def _value_grep_other_cond($v; $flags):
  . == $v;

def vgrep($v; $flags):
  if $v | _is_string then
    grep_by(_is_scalar and _value_grep_string_cond($v; $flags))
  else
    grep_by(_is_scalar and _value_grep_other_cond($v; $flags))
  end;
def vgrep($v): vgrep($v; "");

def _buf_grep_any_cond($v; $flags):
  (isempty(tobytesrange | match($v; $flags)) | not)? // false;
def bgrep($v; $flags):
  if $v | _is_string then
    grep_by(_is_scalar and _buf_grep_any_cond($v; $flags))
  else
    grep_by(_is_scalar and _buf_grep_any_cond($v; $flags))
  end;
def bgrep($v): bgrep($v; "");

def grep($v; $flags):
  if $v | _is_string then
    grep_by(_is_scalar and _buf_grep_any_cond($v; $flags) or _value_grep_string_cond($v; $flags))
  else
    grep_by(_is_scalar and _buf_grep_any_cond($v; $flags) or _value_grep_other_cond($v; $flags))
  end;
def grep($v): grep($v; "");

def fgrep($v; $flags):
  grep_by(_is_decode_value and (._name | test($v; $flags))? // false);
def fgrep($v): fgrep($v; "");
