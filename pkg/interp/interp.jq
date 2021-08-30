include "internal";
include "funcs";
include "args";

# will include all per format specific function etc
include "@format/all";

# optional user init
include "@config/init?";

# def readline: #:: [a]| => string
# Read a line.

# def readline($prompt): #:: [a]|(string) => string
# Read a line with prompt.

# def readline($prompt; $completion): #:: [a]|(string;string) => string
# $prompt is prompt to show.
# $completion name of completion function [a](string) => [string],
# it will be called with same input as readline and a string argument being the
# current line from start to current cursor position. Should return possible completions.

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
  {
    addrbase:       16,
    arraytruncate:  50,
    bitsformat:     "snippet",
    bytecolors:     "0-0xff=brightwhite,0=brightblack,32-126:9-13=white",
    color:          (tty.is_terminal and env.CLICOLOR != null),
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
    decode_progress: (env.NODECODEPROGRESS == null),
    depth:           0,
    # TODO: intdiv 2 * 2 to get even number, nice or maybe not needed?
    displaybytes:    (if tty.is_terminal then [intdiv(intdiv(tty.width; 8); 2) * 2, 4] | max else 16 end),
    expr_file:       null,
    include_path:    null,
    join_string:     "\n",
    linebytes:       (if tty.is_terminal then [intdiv(intdiv(tty.width; 8); 2) * 2, 4] | max else 16 end),
    null_input:      false,
    raw_output:      (tty.is_terminal | not),
    raw_string:      false,
    repl:            false,
    sizebase:        10,
    slurp:           false,
    unicode:         (tty.is_terminal and env.CLIUNICODE != null),
    verbose:         false,
  };

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
  if . != null then "\"\(.)\"" | fromjson end;

def _to_options:
  ( {
      addrbase:        (.addrbase | _tonumber),
      arraytruncate:   (.arraytruncate | _tonumber),
      bitsformat:      (.bitsformat | _tostring),
      bytecolors:      (.bytecolors | _tostring),
      color:           (.color | _toboolean),
      colors:          (.colors | _tostring),
      compact:         (.compact | _toboolean),
      decode_progress: (.decode_progress | _toboolean),
      depth:           (.depth | _tonumber),
      displaybytes:    (.displaybytes | _tonumber),
      expr_file:       (.expr_file | _tostring),
      include_path:    (.include_path | _tostring),
      join_string:     (.join_string | _tostring),
      linebytes:       (.linebytes | _tonumber),
      null_input:      (.null_input | _toboolean),
      raw_output:      (.raw_output | _toboolean),
      raw_string:      (.raw_string | _toboolean),
      repl:            (.repl | _toboolean),
      sizebase:        (.sizebase | _tonumber),
      slurp:           (.slurp | _toboolean),
      unicode:         (.unicode | _toboolean),
      verbose:         (.verbose | _toboolean),
    }
  | with_entries(select(.value != null))
  );


# TODO: refactor this
# TODO: completionMode
# TODO: return escaped identifier, not sure current readline implementation supports
# completions that needs to change previous input, ex: .a\t -> ."a \" b" etc
def _complete($e; $cursor_pos):
  def _is_internal: startswith("_") or startswith("$_");
  def _query_index_or_key($q):
    ( ([.[] | eval($q) | type]) as $n
    | if ($n | all(. == "object")) then "."
      elif ($n | all(. == "array")) then "[]"
      else null
      end
    );
  # only complete if at end of there is a whitespace for now
  if ($e[$cursor_pos] | . == "" or . == " ") then
    ( . as $c
    | ( $e[0:$cursor_pos] | _complete_query) as {$type, $query, $prefix}
    | {
        prefix: $prefix,
        names: (
          if $type == "none" then
            ( $c
            | _query_index_or_key($query)
            | if . then [.] else [] end
            )
          else
            ( $c
            | eval($query)
            | ($prefix | _is_internal) as  $prefix_is_internal
            | map(
                select(
                  strings and
                  (_is_ident or $type == "variable") and
                  ((_is_internal | not) or $prefix_is_internal or $type == "index") and
                  startswith($prefix)
                )
              )
            | unique
            | sort
            | if length == 1 and .[0] == $prefix then
                ( $c
                | _query_index_or_key($e)
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
def _complete($e): _complete($e; $e | length);


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
      , "[\(length)]"
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
    , if length > 1 then ", ..." else empty end
    , "> "
    ]
  ) | join("");


# TODO: better way? what about nested eval errors?
def _eval_is_compile_error: type == "object" and .error != null and .what != null;
def _eval_compile_error_tostring:
  "\(.filename // "src"):\(.line):\(.column): \(.error)";

def _eval($e; $filename; f; on_error; on_compile_error):
  ( _default_options(_build_default_options) as $_
  | try eval($e; $filename) | f
    catch
      if _eval_is_compile_error then on_compile_error
      else on_error
      end
  );

def _repl_display: display({depth: 1});
def _repl_on_error:
  ( if _eval_is_compile_error then _eval_compile_error_tostring end
  | (_error_str | println)
  );
def _repl_on_compile_error: _repl_on_error;
def _repl_eval($e): _eval($e; "repl"; _repl_display; _repl_on_error; _repl_on_compile_error);

# run read-eval-print-loop
def repl($opts; iter): #:: a|(Opts) => @
  def _read_expr:
    # both _prompt and _complete want arrays
    ( [iter]
    | readline(_prompt; "_complete")
    | trim
    );
  def _repl:
    ( . as $c
    | try
        ( _read_expr as $e
        | if $e != "" then
            (iter | _repl_eval($e))
          else
            empty
          end
        , _repl
        )
      catch
        if . == "interrupt" then $c | _repl
        elif . == "eof" then empty
        else error(.)
        end
    );
  ( _options_stack(. + [$opts]) as $_
  | finally(
      _repl;
      _options_stack(.[:-1])
    )
  );
# same as repl({})
def repl($opts): repl($opts; .);
def repl: repl({}; .); #:: a| => @

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
def _cli_expr_eval($e; $filename; f): _eval($e; $filename; f; _cli_expr_on_error; _cli_expr_on_compile_error);
def _cli_expr_eval($e; $filename): _eval($e; $filename; .; _cli_expr_on_error; _cli_expr_on_compile_error);


# next valid input
def input:
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
      , input
      )
  | try
      decode(_parsed_args.decode_format)
    catch
      ( . as $err
      | _input_decode_errors(. += {($h): $err}) as $_
      | "\($h): failed to decode (\(_parsed_args.decode_format)), try -d FORMAT to force"
      | (_error_str | _errorln)
      , input
      )
  );

# iterate all valid inputs
def inputs:
  try repeat(input)
  catch if . == "break" then empty else error end;

def input_filename: _input_filename;

def var: _variables;
def var($k; f):
  ( . as $c
  | if ($k | _is_ident | not) then error("invalid variable name: \($k)") end
  | _variables(.[$k] |= f)
  | empty
  );
def var($k): . as $c | var($k; $c);

# TODO: introspect and show doc, reflection somehow?
def help:
  ( builtins[]
  , "^C interrupt"
  , "^D exit REPL"
  ) | println;

def _main:
  def _formats_list:
    [ ["Name:", "Description:"]
    , ( formats
      | to_entries[]
      | [(.key+" "), .value.description]
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
  def _opts($version):
    {
      "compact": {
        short: "-c",
        long: "--compact",
        description: "Compact output",
        bool: true
      },
      "decode_format": {
        short: "-d",
        long: "--decode",
        description: "Decode format",
        default: "probe",
        string: "NAME"
      },
      "expr_file": {
        short: "-f",
        long: "--file",
        description: "Read EXPR from file",
        string: "PATH"
      },
      "formats": {
        long: "--formats",
        description: "Show supported formats",
        bool: true
      },
      "help": {
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
      # TODO: multiple paths
      "include_path": {
        short: "-L",
        long: "--include-path",
        description: "Include search path",
        string: "PATH"
      },
      "null_output": {
        short: "-0",
        long: "--null-output",
        description: "Null byte between outputs",
        bool: true
      },
      "null_input": {
        short: "-n",
        long: "--null-input",
        description: "Null input (use input/0 and inputs/0 to read input)",
        bool: true
      },
      "options": {
        short: "-o",
        long: "--option",
        description: "Set option, eg: color=true",
        object: "KEY=VALUE",
        default: {},
        help_default: _build_default_options
      },
      "raw_string": {
        short: "-r",
        # for jq compat, is called raw string internally, raw output is
        # if we can output raw bytes or not
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
      "version": {
        short: "-v",
        long: "--version",
        description: "Show version (\($version))",
        bool: true
      },
    };
  def _banner:
    ( "fq - jq for files"
    , "Tool, language and decoders for exploring binary data."
    , "For more information see https://github.com/wader/fq"
    );
  def _usage($arg0; $version):
    "Usage: \($arg0) [OPTIONS] [--] [EXPR] [FILE...]";
  ( . as {$version, $args, args: [$arg0]}
  # make sure we don't unintentionally use . to make things clearer
  | null
  | ( try args_parse($args[1:]; _opts($version))
      catch halt_error(_exit_code_args_error)
    ) as {parsed: $parsed_args, $rest}
  # store parsed arguments, .decode_format is used by input/0
  | _parsed_args($parsed_args) as $_
  | _default_options(_build_default_options) as $_
  | _options_stack(
      [ ($parsed_args.options | _to_options)
      + ({
          compact: $parsed_args.compact,
          expr_file: $parsed_args.expr_file,
          join_string: (
            if $parsed_args.join_output then ""
            elif $parsed_args.null_output then "\u0000"
            else null
            end
          ),
          include_path: $parsed_args.include_path,
          null_input: $parsed_args.null_input,
          repl: $parsed_args.repl,
          raw_string: (
            if $parsed_args.raw_string
              or $parsed_args.join_output
              or $parsed_args.null_output
            then true
            else null
            end
          ),
          slurp: $parsed_args.slurp
        } | with_entries(select(.value != null)))
      ]
    ) as $_
  | options as $opts
  | if $parsed_args.help then
      ( _banner
      , ""
      , _usage($arg0; $version)
      , args_help_text(_opts($version))
      ) | println
    elif $parsed_args.version then
      $version | println
    elif $parsed_args.formats then
      _formats_list | println
    elif ($rest | length) == 0 and (($opts.repl | not) and ($opts.expr_file | not)) then
      ( (( _usage($arg0; $version), "\n") | stderr)
      , null | halt_error(_exit_code_args_error)
      )
    else
      # use finally as display etc prints and results in empty
      finally(
        ( { null_input: ($opts.null_input == true) }
        | if $opts.expr_file then
            ( .expr = ($opts.expr_file | open | tobytes | tostring)
            | .expr_filename = $opts.expr_file
            | .filenames = $rest
            )
          else
            ( .expr = ($rest[0] // ".")
            | .expr_filename = "arg"
            | .filenames = $rest[1:]
            )
          end
        | if $opts.repl and .filenames == [] then
            .null_input = true
          elif .filenames == [] then
            .filenames = ["-"]
          end
        | . as {$expr, $expr_filename, $filenames, $null_input}
        | _include_paths([
            $opts.include_path // empty
          ]) as $_
        | _input_filenames($filenames) as $_ # store inputs
        | if $opts.repl then
            ( [ if $null_input then null
                elif $opts.slurp then [inputs]
                else inputs
                end
              ]
            | ( [.[] | _cli_expr_eval($expr; $expr_filename)]
              | repl({}; .[])
              )
            )
          else
            ( if $null_input then null
              elif $opts.slurp then [inputs]
              else inputs
              end
              # iterate inputs
              | ( _cli_last_expr_error(null) as $_
                | _cli_expr_eval($expr; $expr_filename; _repl_display)
                )
            )
          end
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
