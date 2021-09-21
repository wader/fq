include "internal";
include "funcs";
include "args";
include "query";
# generated decode functions per format and format helpers
include "formats";
# optional user init
include "@config/init?";

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


# . will have additional array of options taking priority
# NOTE: is called from go *interp.Interp Options()
def options($opts):
  [_default_options] + _options_stack + $opts | add;
def options: options([{}]);

def _obj_to_csv_kv:
  [to_entries[] | [.key, .value] | join("=")] | join(",");

def _build_default_options:
  ( (null | stdout) as $stdout
  | {
      addrbase:       16,
      arg:            [],
      argjson:        [],
      arraytruncate:  50,
      bitsformat:     "snippet",
      bytecolors:     "0-0xff=brightwhite,0=brightblack,32-126:9-13=white",
      color:          ($stdout.is_terminal and (env.NO_COLOR | . == null or . == "")),
      colors: (
        {
          null: "brightblack",
          false: "yellow",
          true: "yellow",
          number: "cyan",
          string: "green",
          objectkey: "brightblue",
          array: "white",
          object: "white",
          index: "white",
          value: "white",
          error: "brightred",
          dumpheader: "yellow+underline",
          dumpaddr: "yellow"
        } | _obj_to_csv_kv
      ),
      compact:         false,
      decode_file:      [],
      decode_format:   "probe",
      decode_progress: (env.NO_DECODE_PROGRESS == null),
      depth:           0,
      # TODO: intdiv 2 * 2 to get even number, nice or maybe not needed?
      displaybytes:    (if $stdout.is_terminal then [intdiv(intdiv($stdout.width; 8); 2) * 2, 4] | max else 16 end),
      expr:            ".",
      expr_file:       null,
      expr_eval_path:  "arg",
      filenames:       ["-"],
      include_path:    null,
      join_string:     "\n",
      linebytes:       (if $stdout.is_terminal then [intdiv(intdiv($stdout.width; 8); 2) * 2, 4] | max else 16 end),
      null_input:      false,
      rawfile:         [],
      raw_output:      ($stdout.is_terminal | not),
      raw_string:      false,
      repl:            false,
      sizebase:        10,
      show_formats:    false,
      show_help:       false,
      show_options:    false,
      slurp:           false,
      string_input:    false,
      unicode:         ($stdout.is_terminal and env.CLIUNICODE != null),
      verbose:         false,
    }
  );

def _toboolean:
  try
    if . == "true" then true
    elif . == "false" then false
    else tonumber != 0
    end
  catch
    null;

def _tonumber:
  try tonumber catch null;

def _tostring:
  if . != null then
    ( "\"\(.)\""
    | try
        ( fromjson
        | if type != "string" then error end
        )
      catch null
    )
  end;

def _toarray(f):
  try
    ( fromjson
    | if type == "array" and (all(f) | not) then null end
    )
  catch null;

def _is_string_pair:
  type == "array" and length == 2 and all(type == "string");

def _to_options:
  ( {
      addrbase:        (.addrbase | _tonumber),
      arg:             (.arg | _toarray(_is_string_pair)),
      argjson:         (.argjson | _toarray(_is_string_pair)),
      arraytruncate:   (.arraytruncate | _tonumber),
      bitsformat:      (.bitsformat | _tostring),
      bytecolors:      (.bytecolors | _tostring),
      color:           (.color | _toboolean),
      colors:          (.colors | _tostring),
      compact:         (.compact | _toboolean),
      decode_file:     (.decode_file | _toarray(type == "string")),
      decode_format:   (.decode_format | _tostring),
      decode_progress: (.decode_progress | _toboolean),
      depth:           (.depth | _tonumber),
      displaybytes:    (.displaybytes | _tonumber),
      expr:            (.expr | _tostring),
      expr_file:       (.expr_file | _tostring),
      filename:        (.filenames | _toarray(type == "string")),
      include_path:    (.include_path | _tostring),
      join_string:     (.join_string | _tostring),
      linebytes:       (.linebytes | _tonumber),
      null_input:      (.null_input | _toboolean),
      rawfile:         (.rawfile| _toarray(_is_string_pair)),
      raw_output:      (.raw_output | _toboolean),
      raw_string:      (.raw_string | _toboolean),
      repl:            (.repl | _toboolean),
      sizebase:        (.sizebase | _tonumber),
      show_formats:    (.show_formats | _toboolean),
      show_help:       (.show_help | _toboolean),
      show_options:    (.show_options | _toboolean),
      slurp:           (.slurp | _toboolean),
      string_input:    (.string_input | _toboolean),
      unicode:         (.unicode | _toboolean),
      verbose:         (.verbose | _toboolean),
    }
  | with_entries(select(.value != null))
  );

# TODO: currently only make sense to allow keywords start  start a term or directive
def _complete_keywords:
  [
    "and",
    #"as",
    #"break",
    #"catch",
    "def",
    #"elif",
    #"else",
    #"end",
    "false",
    "foreach",
    "if",
    "import",
    "include",
    "label",
    "module",
    "null",
    "or",
    "reduce",
    #"then",
    "true",
    "try"
  ];

def _complete_scope:
  [scope[], _complete_keywords[]];

# TODO: handle variables via ast walk?
# TODO: refactor this
# TODO: completionMode
# TODO: return escaped identifier, not sure current readline implementation supports
# modifying "previous" characters if quoting is needed
# completions that needs to change previous input, ex: .a\t -> ."a \" b" etc
def _complete($line; $cursor_pos):
  # TODO: reverse this? word or non-ident char?
  def _is_separator: . as $c | " .;[]()|=" | contains($c);
  def _is_internal: startswith("_") or startswith("$_");
  def _query_index_or_key($q):
    ( ([.[] | eval($q) | type]) as $n
    | if ($n | all(. == "object")) then "."
      elif ($n | all(. == "array")) then "[]"
      else null
      end
    );
  # only complete if at end or there is a whitespace for now
  if ($line[$cursor_pos] | . == "" or _is_separator) then
    ( . as $c
    | $line[0:$cursor_pos]
    | . as $line_query
    # expr -> map(partial-expr | . | f?) | add
    # TODO: move map/add logic to here?
    | _query_completion(
        if .type | . == "func" or . == "var" then "_complete_scope"
        elif .type == "index" then
          if (.prefix | startswith("_")) then "_extkeys"
          else "keys"
          end
        else error("unreachable")
        end
      ) as {$type, $query, $prefix}
    | {
        prefix: $prefix,
        names: (
          if $type == "none" then
            ( $c
            | _query_index_or_key($line_query)
            | if . then [.] else [] end
            )
          else
            ( $c
            | eval($query)
            | ($prefix | _is_internal) as $prefix_is_internal
            | map(
                select(
                  strings and
                  # TODO: var type really needed? just func?
                  (_is_ident or $type == "var") and
                  ((_is_internal | not) or $prefix_is_internal or $type == "index") and
                  startswith($prefix)
                )
              )
            | unique
            | sort
            | if length == 1 and .[0] == $prefix then
                ( $c
                | _query_index_or_key($line_query)
                | if . then [$prefix+.] else [$prefix] end
                )
              end
            )
          end
        )
      }
    )
  else
    {prefix: "", names: []}
  end;
def _complete($line): _complete($line; $line | length);


def _prompt:
  def _type_name_error:
    ( . as $c
    | try
        ( _display_name
        , if ._error then "!" else empty end
        )
      catch ($c | type)
    );
  def _path_prefix:
    (._path? // []) | if . == [] then "" else path_to_expr + " " end;
  def _preview:
    if format != null or type != "array" then
      _type_name_error
    else
      ( "["
      , if length > 0 then (.[0] | _type_name_error) else empty end
      , if length > 1 then ", ..." else empty end
      , "]"
      , if length > 1 then "[\(length)]" else empty end
      )
    end;
  ( [ (_options_stack | length | if . > 2 then ((.-2) * ">") + " " else empty end)
    , if length == 0 then
        "empty"
      else
        ( .[0]
        | _path_prefix
        , _preview
        )
      end
    , if length > 1 then ", [\(length)]" else empty end
    , "> "
    ]
  ) | join("");


# TODO: better way? what about nested eval errors?
def _eval_is_compile_error: type == "object" and .error != null and .what != null;
def _eval_compile_error_tostring:
  "\(.filename // "src"):\(.line):\(.column): \(.error)";

def _eval($expr; $filename; f; on_error; on_compile_error):
  ( _default_options(_build_default_options) as $_
  | try eval($expr; $filename) | f
    catch
      if _eval_is_compile_error then on_compile_error
      else on_error
      end
  );

def _repl_display: _display({depth: 1});
def _repl_on_error:
  ( if _eval_is_compile_error then _eval_compile_error_tostring end
  | (_error_str | println)
  );
def _repl_on_compile_error: _repl_on_error;
def _repl_eval($expr): _eval($expr; "repl"; _repl_display; _repl_on_error; _repl_on_compile_error);

# run read-eval-print-loop
def _repl($opts): #:: a|(Opts) => @
  def _read_expr:
    # both _prompt and _complete want arrays
    ( . as $c
    | _readline(_prompt; "_complete")
    | if trim == "" then
        $c | _read_expr
      end
    );

  def _repl_loop:
    ( . as $c
    | try
        ( _read_expr
        | . as $expr
        | try _query_fromstring
          # TODO: nicer way to set filename for error message
          catch (. | .filename = "repl")
        | if _query_pipe_last | _query_is_func("repl") then
            ( _query_slurp_wrap(_query_func_rename("_repl_slurp"))
            | _query_tostring as $wrap_expr
            | $c
            | _repl_eval($wrap_expr)
            )
          else
            ( $c
            | .[]
            | _repl_eval($expr)
            )
          end
        )
      catch
        if . == "interrupt" then empty
        elif . == "eof" then error("break")
        elif _eval_is_compile_error then _repl_on_error
        else error
        end
    );
  ( _options_stack(. + [$opts]) as $_
  | _finally(
      _repeat_break(_repl_loop);
      _options_stack(.[:-1])
    )
  );

def _repl_slurp($opts): _repl($opts);
def _repl_slurp: _repl({});

# just gives error, call appearing last will be renamed to _repl_slurp
def repl($_):
  if options.repl then error("repl must be last")
  else error("repl can only be used from interactive repl")
  end;
def repl: repl(null);


def _cli_expr_on_error:
  ( . as $err
  | _cli_last_expr_error($err) as $_
  | (_error_str | _errorln)
  );
def _cli_expr_on_compile_error:
  ( _eval_compile_error_tostring
  | halt_error(_exit_code_compile_error)
  );
# _cli_expr_eval halts on compile errors
def _cli_expr_eval($expr; $filename; f):
  _eval($expr; $filename; f; _cli_expr_on_error; _cli_expr_on_compile_error);
def _cli_expr_eval($expr; $filename):
  _eval($expr; $filename; .; _cli_expr_on_error; _cli_expr_on_compile_error);


# TODO: introspect and show doc, reflection somehow?
def help:
  ( "Type jq expression to evaluate"
  , "\\t          Auto completion"
  , "Up/Down     History"
  , "^C          Interrupt execution"
  , "... | repl  Start a new REPL"
  , "^D          Exit REPL"
  ) | println;

def display($opts): _display($opts);
def display: _display({});
def d($opts): _display($opts);
def d: _display({});
def full($opts): _display({arraytruncate: 0} + $opts);
def full: full({});
def f($opts): full($opts);
def f: full;
def verbose($opts): _display({verbose: true, arraytruncate: 0} + $opts);
def verbose: verbose({});
def v($opts): verbose($opts);
def v: verbose;

def formats:
  _registry.formats;

def _esc: "\u001b";
def _ansi:
  {
    clear_line: "\(_esc)[2K",
  };

# null input means done, otherwise {approx_read_bytes: 123, total_size: 123}
# TODO: decode provide even more detailed progress, post-process sort etc?
def _decode_progress:
  # _input_filenames is remaning files to read
  ( (_input_filenames | length) as $inputs_len
  | ( options.filenames | length) as $filenames_len
  | _ansi.clear_line
  , "\r"
  , if . != null then
      ( if $filenames_len > 1 then
          "\($filenames_len - $inputs_len)/\($filenames_len) \(_input_filename) "
        else empty
        end
      , "\((.approx_read_bytes / .total_size * 100 | _numbertostring(1)))%"
      )
    else empty
    end
  | stderr
  );

def decode($name; $opts):
  ( options as $opts
  | (null | stdout) as $stdout
  | _decode(
      $name;
      $opts + {
        _progress: (
          if $opts.decode_progress and $opts.repl and $stdout.is_terminal then
            "_decode_progress"
          else null
          end
        )
      }
    )
  );
def decode($name): decode($name; {});
def decode: decode(options.decode_format; {});

# next valid input
def input:
  def _input($opts; f):
    ( _input_filenames
    | if length == 0 then error("break") end
    | [.[0], .[1:]] as [$h, $t]
    | _input_filenames($t)
    | _input_filename(null) as $_
    | $h
    | try
        ( open
        | _input_filename($h) as $_
        | .
        )
      catch
        ( . as $err
        | _input_io_errors(. += {($h): $err}) as $_
        | $err
        | (_error_str | _errorln)
        , _input($opts; f)
        )
    | try f
      catch
        ( . as $err
        | _input_decode_errors(. += {($h): $err}) as $_
        | "\($h): failed to decode (\($opts.decode_format)), try -d FORMAT to force"
        | (_error_str | _errorln)
        , _input($opts; f)
        )
    );
  # TODO: don't rebuild options each time
  ( options as $opts
  # TODO: refactor into def
  # this is a bit strange as jq for --raw-string can return string instead
  # with data from multiple inputs
  | if $opts.string_input then
      ( _input_strings_lines
      | if . then
          # we're already iterating lines
          if length == 0 then error("break")
          else
            ( [.[0], .[1:]] as [$h, $t]
            | _input_strings_lines($t)
            | $h
            )
          end
        else
          ( [_repeat_break(_input($opts; tobytes | tostring))]
          | . as $chunks
          | if $opts.slurp then
              # jq --raw-input combined with --slurp reads all inputs into a string
              # make next input break
              ( _input_strings_lines([]) as $_
              | $chunks
              | join("")
              )
            else
              # TODO: different line endings?
              # jq strips last newline, "a\nb" and "a\nb\n" behaves the same
              # also jq -R . <(echo -ne 'a\nb') <(echo c) produces "a" and "bc"
              if ($chunks | length) > 0 then
                ( _input_strings_lines(
                    ( $chunks
                    | join("")
                    | rtrimstr("\n")
                    | split("\n")
                    )
                  ) as $_
                | input
                )
              else error("break")
              end
            end
          )
        end
      )
    else _input($opts; decode($opts.decode_format))
    end
  );

# iterate all valid inputs
def inputs: _repeat_break(input);

def input_filename: _input_filename;

def var: _variables;
def var($k; f):
  ( . as $c
  | if ($k | _is_ident | not) then error("invalid variable name: \($k)") end
  | _variables(.[$k] |= f)
  | empty
  );
def var($k): . as $c | var($k; $c);


def _main:
  def _formats_list:
    [ ( formats
      | to_entries[]
      | [(.key+"  "), .value.description]
      )
    ]
    | table(
        .;
        map(
          ( . as $rc
          | .string
          | if $rc.column != 1 then rpad(" "; $rc.maxwidth) end
          )
        ) | join("")
      );
  def _opts:
    {
      "arg": {
        long: "--arg",
        description: "Set variable $NAME to string VALUE",
        pairs: "NAME VALUE"
      },
      "argjson": {
        long: "--argjson",
        description: "Set variable $NAME to JSON",
        pairs: "NAME JSON"
      },
      "compact": {
        short: "-c",
        long: "--compact-output",
        description: "Compact output",
        bool: true
      },
      "color_output": {
        short: "-C",
        long: "--color-output",
        description: "Force color output",
        bool: true
      },
      "decode_format": {
        short: "-d",
        long: "--decode",
        description: "Decode format (probe)",
        string: "NAME"
      },
      "decode_file": {
        long: "--decode-file",
        description: "Set variable $NAME to decode of file",
        pairs: "NAME PATH"
      },
      "expr_file": {
        short: "-f",
        long: "--from-file",
        description: "Read EXPR from file",
        string: "PATH"
      },
      "show_formats": {
        long: "--formats",
        description: "Show supported formats",
        bool: true
      },
      "show_help": {
        short: "-h",
        long: "--help",
        description: "Show help",
        bool: true
      },
      "join_output": {
        short: "-j",
        long: "--join-output",
        description: "No newline between outputs",
        bool: true
      },
      "include_path": {
        short: "-L",
        long: "--include-path",
        description: "Include search path",
        array: "PATH"
      },
      "null_output": {
        short: "-0",
        long: "--null-output",
        # for jq compatibility
        aliases: ["--nul-output"],
        description: "Null byte between outputs",
        bool: true
      },
      "null_input": {
        short: "-n",
        long: "--null-input",
        description: "Null input (use input/0 and inputs/0 to read input)",
        bool: true
      },
      "monochrome_output": {
        short: "-M",
        long: "--monochrome-output",
        description: "Force monochrome output",
        bool: true
      },
      "option": {
        short: "-o",
        long: "--option",
        description: "Set option, eg: color=true",
        object: "KEY=VALUE",
      },
      "show_options": {
        long: "--options",
        description: "Show all options",
        bool: true
      },
      "string_input": {
        short: "-R",
        long: "--raw-input",
        description: "Read raw input strings (don't decode)",
        bool: true
      },
      "rawfile": {
        long: "--rawfile",
        description: "Set variable $NAME to string content of file",
        pairs: "NAME PATH"
      },
      "raw_string": {
        short: "-r",
        # for jq compat, is called raw string internally, "raw output" is if
        # we can output raw bytes or not
        long: "--raw-output",
        description: "Raw string output (without quotes)",
        bool: true
      },
      "repl": {
        short: "-i",
        long: "--repl",
        description: "Interactive REPL",
        bool: true
      },
      "slurp": {
        short: "-s",
        long: "--slurp",
        description: "Read (slurp) all inputs into an array",
        bool: true
      },
      "show_version": {
        short: "-v",
        long: "--version",
        description: "Show version",
        bool: true
      },
    };
  def _banner:
    ( "fq - jq for files"
    , "Tool, language and format decoders for exploring binary data."
    , "For more information see https://github.com/wader/fq"
    );
  def _usage($arg0):
    "Usage: \($arg0) [OPTIONS] [--] [EXPR] [FILE...]";
  ( . as {$version, $args, args: [$arg0]}
  | (null | [stdin, stdout]) as [$stdin, $stdout]
  # make sure we don't unintentionally use . to make things clearer
  | null
  | ( try _args_parse($args[1:]; _opts)
      catch halt_error(_exit_code_args_error)
    ) as {parsed: $parsed_args, $rest}
  | _build_default_options as $default_opts
  | _default_options($default_opts) as $_
  # combine --args and -o key=value args
  | ( $default_opts
    + $parsed_args
    + ($parsed_args.option | _to_options)
    ) as $args_opts
  | _options_stack(
      [ $args_opts
      + ( {
            argjson: (
              ( $args_opts.argjson
              | if . then
                  map(
                    ( . as $a
                    | .[1] |=
                      try fromjson
                      catch
                        ( "--argjson \($a[0]): \(.)"
                        | halt_error(_exit_code_args_error)
                        )
                    )
                  )
                end
              )
            ),
            color: (
              if $args_opts.monochrome_output == true then false
              elif $args_opts.color_output == true then true
              end
            ),
            decode_file: (
              ( $args_opts.decode_file
              | if . then
                  map(
                    ( . as $a
                    | .[1] |=
                      try (open | decode($args_opts.decode_format))
                      catch
                        ( "--decode-file \($a[0]): \(.)"
                        | halt_error(_exit_code_args_error)
                        )
                    )
                  )
                end
              )
            ),
            expr: (
              # if -f was used, all rest non-args are filenames
              # otherwise first is expr rest is filesnames
              ( $args_opts.expr_file
              | if . then
                  try (open | tobytes | tostring)
                  catch halt_error(_exit_code_args_error)
                else $rest[0] // null
                end
              )
            ),
            expr_eval_path: $args_opts.expr_file,
            filenames: (
              ( if $args_opts.expr_file then $rest
                else $rest[1:]
                end
              | if . == [] then null end
              )
            ),
            join_string: (
              if $args_opts.join_output then ""
              elif $args_opts.null_output then "\u0000"
              else null
              end
            ),
            null_input: (
              ( if $args_opts.expr_file then $rest
                else $rest[1:]
                end
              | if . == [] and $args_opts.repl then true
                else null
                end
              )
            ),
            rawfile: (
              ( $args_opts.rawfile
              | if . then
                  ( map(.[1] |=
                      try (open | tobytes | tostring)
                      catch halt_error(_exit_code_args_error)
                    )
                  )
                end
              )
            ),
            raw_string: (
              if $args_opts.raw_string
                or $args_opts.join_output
                or $args_opts.null_output
              then true
              else null
              end
            )
          }
        | with_entries(select(.value != null))
        )
      ]
    ) as $_
  | options as $opts
  | if $opts.show_help then
      ( _banner
      , ""
      , _usage($arg0)
      , args_help_text(_opts)
      ) | println
    elif $opts.show_version then
      $version | println
    elif $opts.show_formats then
      _formats_list | println
    elif $opts.show_options then
      $opts | display
    elif
      ( ($rest | length) == 0 and
        ($opts.repl | not) and
        ($opts.expr_file | not) and
        $stdin.is_terminal and $stdout.is_terminal
      ) then
      ( (( _usage($arg0), "\n") | stderr)
      , null | halt_error(_exit_code_args_error)
      )
    else
      # use _finally as display etc prints and results in empty
      _finally(
        # store some globals
        ( _include_paths($opts.include_path) as $_
        | _input_filenames($opts.filenames) as $_
        | _variables(
            ( $opts.arg +
              $opts.argjson +
              $opts.rawfile +
              $opts.decode_file
            | map({key: .[0], value: .[1]})
            | from_entries
            )
          )
        | ( def _inputs:
              ( if $opts.null_input then null
                # note jq --slurp --raw-string is special, will be just
                # a string not an array
                elif $opts.string_input then inputs
                elif $opts.slurp then [inputs]
                else inputs
                end
              );
            if $opts.repl then
              ( [_inputs]
              | map(_cli_expr_eval($opts.expr; $opts.expr_eval_path))
              | _repl({})
              )
            else
              ( _inputs
              # iterate all inputs
              | _cli_last_expr_error(null) as $_
              | _cli_expr_eval($opts.expr; $opts.expr_eval_path; _repl_display)
              )
            end
          )
        )
        ; # finally
        ( if _input_io_errors then
            null | halt_error(_exit_code_input_io_error)
          end
        | if _input_decode_errors then
            null | halt_error(_exit_code_input_decode_error)
          end
        | if _cli_last_expr_error then
            null | halt_error(_exit_code_expr_error)
          end
        )
      )
    end
  );
