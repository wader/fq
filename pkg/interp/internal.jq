include "ansi";

# is here to be defined as early as possible to allow debugging
# TODO: move some _* to builtin.jq etc?

def _stdio($name):
  if . == null then _stdio_info($name)
  else _stdio_write($name)
  end;
def _stdio($name; $l):
  _stdio_read($name; $l);

def _stdin: _stdio("stdin");
def _stdin($l): _stdio("stdin"; $l);
def _stdout: _stdio("stdout");
def _stdout($l): _stdio("stdout"; $l);
def _stderr: _stdio("stderr");
def _stderr($l): _stdio("stderr"; $l);

def stdin_tty: null | _stdin;
def stdout_tty: null | _stdout;

def print: tostring | _stdout;
def println: ., "\n" | print;
def printerr: tostring | _stderr;
def printerrln: ., "\n" | printerr;

# jq compat
def debug: (["DEBUG:", .] | tojson | printerrln), .;
def debug(f): (f | debug | empty), .;
# output raw string or compact json to stderr and let input thru
def stderr: printerr, .;

def _fatal_error($code): "error: \(.)\n" | halt_error($code);

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

def _slurps: _global_var("slurps");
def _slurps(f): _global_var("slurps"; f);

# call f and finally eval fin even if empty or error.
# _finally(1; debug)
# _finally(null; debug)
# _finally(empty; debug)
# _finally(1,2,3; debug)
# _finally({a:1}; debug)
# _finally(error("a"); debug)
# _finally(error("a"); empty)
def _finally(f; fin):
  try
    ( f
    , (fin | empty)
    )
  catch
    ( (fin | empty)
    , error
    );

# TODO: figure out a saner way to force int
def _to_int: (. % (. + 1));

# integer division
# inspired by https://github.com/itchyny/gojq/issues/63#issuecomment-765066351
def _intdiv($a; $b):
  ( ($a | _to_int) as $a
  | ($b | _to_int) as $b
  | ($a - ($a % $b)) / $b
  );

# escape \ and "
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

def _is_null: type == "null";
def _is_string: type == "string";
def _is_number: type == "number";
def _is_boolean: type == "boolean";
def _is_array: type == "array";
def _is_object: type == "object";
def _is_scalar: (_is_array or _is_object) | not;

# valid jq identifier, start with alpha or underscore then zero or more alpha, num or underscore
def _is_ident: _is_string and test("^[a-zA-Z_][a-zA-Z_0-9]*$");

def _is_context_canceled_error: . == "context canceled";

def _error_str($contexts): (["error"] + $contexts + [.]) | join(": ");
def _error_str: _error_str([]);

# TODO: escape for safe key names
# path ["a", 1, "b"] -> "a[1].b"
def _path_to_expr($opts):
  ( if length == 0 or (.[0] | type) != "string" then
      [""] + .
    end
  | map(
      if _is_number then
        ( ("[" | _ansi_if($opts; "array"))
        , _ansi_if($opts; "number")
        , ("]" | _ansi_if($opts; "array"))
        )      else
        ( "."
        , # empty (special case for leading index or empty path) or key
          if . == "" or _is_ident then _ansi_if($opts; "objectkey")
          else
            "\"\(_escape_ident)\"" | _ansi_if($opts; "string")
          end
        )
      end
    )
  | join("")
  );
def _path_to_expr: _path_to_expr(null);

# TODO: don't use eval? should support '.a.b[1]."c.c"' and escapes?
def _expr_to_path:
  ( if type != "string" then error("require string argument") end
  | _eval("null | path(\(.))"; {})
  );

# helper to build path query/generate functions for tree structures with
# non-unique children, ex: mp4_path
def _tree_path(children; name; $v):
  def _lookup:
    # add implicit zeros to get first value
    # ["a", "b", 1] => ["a", 0, "b", 1]
    def _normalize_path:
      ( . as $np
      | if $np | last | _is_string then $np+[0] end
      # state is [path acc, possible pending zero index]
      | ( reduce .[] as $np ([[], []];
          if $np | _is_string then
            [(.[0]+.[1]+[$np]), [0]]
          else
            [.[0]+[$np], []]
          end
        ))
      )[0];
    ( . as $c
    | $v
    | _expr_to_path
    | _normalize_path
    | reduce .[] as $n (
        $c;
        if . then
          if $n | _is_string then
            children | map(select(name == $n))
          else
            .[$n]
          end
        else empty
        end
      )
    );
  def _path:
    [ . as $r
    | $v._path as $p
    | foreach range(($p | length)/2) as $i (
        null;
        null;
        ( ($r | getpath($p[0:($i+1)*2]) | name) as $name
        | [($r | getpath($p[0:($i+1)*2-1]))[] | name][0:$p[($i*2)+1]+1] as $before
        | [ $name
          , ($before | map(select(. == $name)) | length)-1
          ]
        )
      )
    | [ ".", .[0],
      (.[1] | if . == 0 then empty else "[", ., "]" end)
      ]
    ]
    | flatten
    | join("");
  if $v | _is_string then _lookup
  else _path
  end;
