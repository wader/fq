# is here to be defined as early as possible to allow debugging
# TODO: move some _* to builtin.jq etc?

def stdin_tty: null | _stdin;
def stdout_tty: null | _stdout;

def print: tostring | _stdout;
def println: ., "\n" | print;
def printerr: tostring | _stderr;
def printerrln: ., "\n" | printerr;

# jq compat
def debug:
  ( ((["DEBUG", .] | tojson) | printerrln)
  , .
  );
def debug(f): . as $c | f | debug | $c;
# jq compat, output to compact json to stderr and let input thru
def stderr:
  ( (tojson | printerr)
  , .
  );

# try to be same exit codes as jq
# TODO: jq seems to halt processing inputs on JSON decode error but not IO errors,
# seems strange.
# jq '(' <(echo 1) <(echo 2) ; echo $? => 3 and no inputs processed
# jq '.' missing <(echo 2) ; echo $? => 2 and continues process inputs
# jq '.' <(echo 'a') <(echo 123) ; echo $? => 4 and stops process inputs
# jq '.' missing <(echo 'a') <(echo 123) ; echo $? => 2 ???
# jq '"a"+.' <(echo '"a"') <(echo 1) ; echo $? => 5
# jq '"a"+.' <(echo 1) <(echo '"a"') ; echo $? => 0
def _exit_code_args_error: 2;
def _exit_code_input_io_error: 2;
def _exit_code_compile_error: 3;
def _exit_code_input_decode_error: 4;
def _exit_code_expr_error: 5;

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

# eval f and finally eval fin even if empty or error.
# _finally(1; debug)
# _finally(null; debug)
# _finally(error("a"); debug)
# _finally(empty; debug)
# _finally(1,2,3; debug)
# _finally({a:1}; debug)
def _finally(f; fin):
  try
    ( f
    , (fin | empty)
    )
  catch
    ( fin as $_
    | error
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
  [ (.filename | if . == "" then "expr" end)
  , if .line != 1 or .column != 0 then "\(.line):\(.column)" else empty end
  , " \(.error)"
  ] | join(":");
def _eval($expr; $filename; f; on_error; on_compile_error):
  try
    eval($expr; $filename) | f
  catch
    if _eval_is_compile_error then on_compile_error
    else on_error
    end;

def _is_scalar:
  type |. != "array" and . != "object";

def _is_context_canceled_error: . == "context canceled";

def _error_str: "error: \(.)";
