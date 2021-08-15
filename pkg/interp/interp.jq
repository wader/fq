include "@builtin/internal";
include "@builtin/funcs";
include "@builtin/args";

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

# TODO: completionMode
# TODO: return escaped identifier, not sure current readline implementation supports
# completions that needs to change previous input, ex: .a\t -> ."a \" b" etc
def _complete($e):
  ( ( $e | _complete_query) as {$type, $query, $prefix}
  | {
      prefix: $prefix,
      names: (
        ( if $type == "function" or $type == "variable" then
            [.[] | eval($query) | scope] | add
          elif $type == "index" then
            [.[] | eval($query) | keys?, _extkeys?] | add
          else
            []
          end
        | map(select(strings and _is_ident and startswith($prefix)))
        | unique
        | sort
        )
      )
    }
  );

def _obj_to_csv_kv:
  [to_entries[] | [.key, .value] | join("=")] | join(",");

def _build_default_options:
  {
    depth:          0,
    verbose:        false,
    bitsformat:     "snippet",
    color:          (tty.is_terminal and env.CLICOLOR != null),
    unicode:        (tty.is_terminal and env.CLIUNICODE != null),
    raw:            (tty.is_terminal | not),
    # TODO: div 2 * 2 to get even number, nice or maybe not needed?
    linebytes:      (if tty.is_terminal then [((tty.width div 8) div 2) * 2, 4] | max else 16 end),
    displaybytes:   (if tty.is_terminal then [((tty.width div 8) div 2) * 2, 4] | max else 16 end),
    addrbase:       16,
    sizebase:       10,
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
    bytecolors:     "0-0xff=brightwhite,0=brightblack,32-126:9-13=white",
    decodeprogress: (env.NODECODEPROGRESS == null),
  };

def _eval_options:
  ( {
      depth:          (.depth | if . then eval(.) else null end),
      bitsformat:     .bitsformat,
      verbose:        (.verbose | if . then eval(.) else null end),
      decodeprogress: (.decodeprogress | if . then eval(.) else null end),
      color:          (.color | if . then eval(.) else null end),
      unicode:        (.unicode | if . then eval(.) else null end),
      raw:            (.raw | if . then eval(.) else null end),
      linebytes:      (.linebytes | if . then eval(.) else null end),
      displaybytes:   (.displaybytes | if . then eval(.) else null end),
      addrbase:       (.addrbase | if . then eval(.) else null end),
      sizebase:       (.sizebase | if . then eval(.) else null end),
    }
  | with_entries(select(.value != null))
  );

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
    if (._format? // false) or type != "array" then
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

def _eval($e; f; on_error; on_compile_error):
  ( _default_options(_build_default_options) as $_
  | try eval($e) | f
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
def _repl_eval($e): _eval($e; _repl_display; _repl_on_error; _repl_on_compile_error);

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
def _cli_expr_eval($e; f): _eval($e; f; _cli_expr_on_error; _cli_expr_on_compile_error);
def _cli_expr_eval($e): _eval($e; .; _cli_expr_on_error; _cli_expr_on_compile_error);

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
def var($k):
  ( . as $c
  | if ($k | _is_ident | not) then error("invalid variable name \"\($k)\"") end
  | _variables(
      .[$k] |= ($c | if . == null then empty end)
    )
  | empty
  );

# . will have additional array of options taking priority
# NOTE: is used by go Options()
def options($opts):
  [_default_options] + _options_stack + $opts | add;
def options: options([{}]);

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
      "file": {
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
        description: "Null input (can still use input/0 and inputs/0)",
        bool: true
      },
      "options": {
        short: "-o",
        long: "--option",
        description: "Set option, eg: color=true",
        object: "KEY=VALUE,...",
        default: {},
        help_default: _build_default_options
      },
      "raw_output": {
        short: "-r",
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
  def _info:
    ( "fq - jq for binary"
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
      [ ($parsed_args.options | _eval_options)
      + {
          repl: ($parsed_args.repl == true),
          rawstring: (
            $parsed_args.raw_output == true
            or $parsed_args.join_output == true
            or $parsed_args.null_output == true
          ),
          joinstring: (
            if $parsed_args.join_output == true then ""
            elif $parsed_args.null_output == true then "\u0000"
            else "\n"
            end
          ),
          compact: ($parsed_args.compact == true),
        }
      ]
    ) as $_
  | if $parsed_args.help then
      ( _info
      , ""
      , _usage($arg0; $version)
      , args_help_text(_opts($version))
      ) | println
    elif $parsed_args.version then
      $version | println
    elif $parsed_args.formats then
      _formats_list | println
    elif ($rest | length) == 0 and (($parsed_args.repl | not) and ($parsed_args.file | not)) then
      ( (( _usage($arg0; $version), "\n") | stderr)
      , null | halt_error(_exit_code_args_error)
      )
    else
      # use finally as display etc prints and results in empty
      finally(
        ( { null_input: ($parsed_args.null_input == true) }
        | if $parsed_args.file then
            ( .expr = ($parsed_args.file | open | string)
            | .filenames = $rest
            )
          else
            ( .expr = ($rest[0] // ".")
            | .filenames = $rest[1:]
            )
          end
        | if $parsed_args.repl and .filenames == [] then
            .null_input = true
          elif .filenames == [] then
            .filenames = ["-"]
          end
        | . as {$expr, $filenames, $null_input}
        | _include_paths([
            $parsed_args.include_path // empty
          ]) as $_
        | _input_filenames($filenames) as $_ # store inputs
        | if $null_input then null
          elif $parsed_args.slurp then [inputs]
          else inputs
          end
          # will iterate zero or more inputs
          | if $parsed_args.repl then [_cli_expr_eval($expr)] | repl({}; .[])
            else
              ( _cli_last_expr_error(null) as $_
              | _cli_expr_eval($expr; _repl_display)
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
