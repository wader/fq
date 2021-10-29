def _global_var($k): _global_state[$k];
def _global_var($k; f): _global_state(_global_state | .[$k] |= f) | .[$k];

def _include_paths: _global_var("include_paths");
def _include_paths(f): _global_var("include_paths"; f);

def _options_stack: _global_var("options_stack");
def _options_stack(f): _global_var("options_stack"; f);

def _options_cache: _global_var("options_cache");
def _options_cache(f): _global_var("options_cache"; f);

def _cli_last_expr_error: _global_var("cli_last_expr_error");
def _cli_last_expr_error(f): _global_var("cli_last_expr_error"; f);

def _input_filename: _global_var("input_filename");
def _input_filename(f): _global_var("input_filename"; f);

def _input_filenames: _global_var("input_filenames");
def _input_filenames(f): _global_var("input_filenames"; f);

def _input_strings: _global_var("input_strings");
def _input_strings(f): _global_var("input_strings"; f);

def _input_strings_lines: _global_var("input_strings_lines");
def _input_strings_lines(f): _global_var("input_strings_lines"; f);

def _input_io_errors: _global_var("input_io_errors");
def _input_io_errors(f): _global_var("input_io_errors"; f);

def _input_decode_errors: _global_var("input_decode_errors");
def _input_decode_errors(f): _global_var("input_decode_errors"; f);

def _variables: _global_var("variables");
def _variables(f): _global_var("variables"; f);

# eval f and finally eval fin even on empty or error
def _finally(f; fin):
  ( try f // (fin | empty)
    catch (fin as $_ | error)
  | fin as $_
  | .
  );

# TODO: figure out a saner way to force int
def _to_int: (. % (. + 1));

# integer division
# inspried by https://github.com/itchyny/gojq/issues/63#issuecomment-765066351
def _intdiv($a; $b):
  ( ($a | _to_int) as $a
  | ($b | _to_int) as $b
  | ($a - ($a % $b)) / $b
  );

def _esc: "\u001b";
def _ansi:
  {
    clear_line: "\(_esc)[2K",
  };

# valid jq identifier, start with alpha or underscore then zero or more alpha, num or underscore
def _is_ident: type == "string" and test("^[a-zA-Z_][a-zA-Z_0-9]*$");
# escape " and \
def _escape_ident: gsub("(?<g>[\\\\\"])"; "\\\(.g)");

# format number with fixed number of decimals
def _numbertostring($decimals):
  ( . as $n
  | [ (. % (. + 1)) # truncate to integer
    , "."
    , foreach range($decimals) as $_ (1; . * 10; ($n * .) % 10)
    ]
  | join("")
  );

def _repeat_break(f):
  try repeat(f)
  catch
    if . == "break" then empty
    else error
    end;

def _recurse_break(f):
  try recurse(f)
  catch
    if . == "break" then empty
    else error
    end;

# TODO: better way? what about nested eval errors?
def _eval_is_compile_error: type == "object" and .error != null and .what != null;
def _eval_compile_error_tostring:
  "\(.filename // "src"):\(.line):\(.column): \(.error)";
def _eval($expr; $filename; f; on_error; on_compile_error):
  try
      eval($expr; $filename) | f
  catch
    if _eval_is_compile_error then on_compile_error
    else on_error
    end;

# TODO: remove one symbolic value is done properly?
def _tovalue:
  ( if _is_decode_value then
      ( ._symbol as $s
      | if $s != "" then $s end
      )
    end
  );

def _is_scalar:
  type |. != "array" and . != "object";

# TODO: error value preview
def _expected_decode_value:
  error("expected a decode value but got: \(. | type) (\(. | tostring))");
# TODO: helper? _is_decode_value?
def _decode_value(f):
  if _is_decode_value then f
  else _expected_decode_value
  end;

def _error_str: "error: \(.)";
def _errorln: ., "\n" | stderr;
