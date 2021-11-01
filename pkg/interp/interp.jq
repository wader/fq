include "internal";
include "options";
include "funcs";
include "grep";
include "args";
include "repl";
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
  def _input_string($opts):
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
    );
  # TODO: don't rebuild options each time
  ( options as $opts
  # this is a bit strange as jq for --raw-input can return one string
  # instead of iterating lines
  | if $opts.string_input then _input_string($opts)
    else _input($opts; decode)
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
        description: "Set option, eg: color=true (use options/0 to see all options)",
        object: "KEY=VALUE",
      },
      "string_input": {
        short: "-R",
        long: "--raw-input",
        description: "Read raw input strings (don't decode)",
        bool: true
      },
      "raw_file": {
        long: "--raw-file",
        # for jq compatibility
        aliases: ["--raw-file"],
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
    ( "fq - jq for binary formats"
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
  | _build_default_fixed_options as $default_fixed_opts
  # combine --args and -o key=value args
  | ( $default_fixed_opts
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
            raw_file: (
              ( $args_opts.raw_file
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
      , ""
      , args_help_text(_opts)
      ) | println
    elif $opts.show_version then
      $version | println
    elif $opts.show_formats then
      _formats_list | println
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
      # use _finally as display etc prints and outputs empty
      _finally(
        # store some globals
        ( _include_paths($opts.include_path) as $_
        | _input_filenames($opts.filenames) as $_
        | _variables(
            ( $opts.arg +
              $opts.argjson +
              $opts.raw_file +
              $opts.decode_file
            | map({key: .[0], value: .[1]})
            | from_entries
            )
          )
        # for inputs a, b, c:
        # repl:       [a,b,c] | repl
        # repl slurp: [[a, b, c]] | repl
        # cli         a, b, c | expr
        # cli slurp   [a ,b c] | expr
        | ( def _inputs:
              ( if $opts.null_input then null
                # note that jq --slurp --raw-input (string_input) is special, will concat
                # all files into one string instead of iterating lines
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
