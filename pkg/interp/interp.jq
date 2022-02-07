include "internal";
include "options";
include "buffer";
include "decode";
include "match";
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

def d($opts): display($opts);
def d: display({});
def da($opts): display({array_truncate: 0} + $opts);
def da: da({});
def dd($opts): display({array_truncate: 0, display_bytes: 0} + $opts);
def dd: dd({});
def dv($opts): display({array_truncate: 0, verbose: true} + $opts);
def dv: dv({});
def ddv($opts): display({array_truncate: 0, display_bytes: 0, verbose: true} + $opts);
def ddv: ddv({});

# next valid input
def input:
  def _input($opts; f):
    ( _input_filenames
    | if length == 0 then error("break") end
    | [.[0], .[1:]] as [$h, $t]
    | _input_filenames($t)
    | _input_filename(null) as $_
    | ($h // "<stdin>") as $name
    | $h
    | try
        # null input here means stdin
        ( open
        | _input_filename($name) as $_
        | .
        )
      catch
        ( . as $err
        | _input_io_errors(. += {($name): $err}) as $_
        | $err
        | (_error_str | printerrln)
        , _input($opts; f)
        )
    | try f
      catch
        ( . as $err
        | _input_decode_errors(. += {($name): $err}) as $_
        | [ "\($name): \($opts.decode_format)"
          , if $err | type == "string" then ": \($err)"
            # TODO: if not string assume decode itself failed for now
            else ": failed to decode (try -d FORMAT)"
            end
          ] | join("")
        | (_error_str | printerrln)
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
  | (_error_str | printerrln)
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
  def _banner:
    ( "fq - jq for binary formats"
    , "Tool, language and decoders for inspecting binary data."
    , "For more information see https://github.com/wader/fq"
    );
  def _usage($arg0):
    "Usage: \($arg0) [OPTIONS] [--] [EXPR] [FILE...]";
  ( . as {$version, $os, $arch, $args, args: [$arg0]}
  # make sure we don't unintentionally use . to make things clearer
  | null
  | ( try _args_parse($args[1:]; _opt_cli_opts)
      catch halt_error(_exit_code_args_error)
    ) as {parsed: $parsed_args, $rest}
  | _opt_build_default_fixed as $default_fixed_opts
  # combine default fixed opt, --args opts and -o key=value opts
  | ( $default_fixed_opts
    + $parsed_args
    + ($parsed_args.option | _opt_cli_arg_options)
    ) as $combined_opts
  # "eval" options
  | _options_stack(
      [ $combined_opts
      + ( {
            argjson: (
              ( $combined_opts.argjson
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
              if $combined_opts.monochrome_output == true then false
              elif $combined_opts.color_output == true then true
              end
            ),
            decode_file: (
              ( $combined_opts.decode_file
              | if . then
                  # [[name, path], ...] pairs
                  map(
                    ( . as $a
                    | .[1] |=
                      try (open | decode($combined_opts.decode_format))
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
              ( $combined_opts.expr_file
              | if . then
                  try (open | tobytes | tostring)
                  catch halt_error(_exit_code_args_error)
                else $rest[0] // null
                end
              )
            ),
            expr_eval_path: $combined_opts.expr_file,
            filenames: (
              ( if $combined_opts.filenames then $combined_opts.filenames
                elif $combined_opts.expr_file then $rest
                else $rest[1:]
                end
              # null means stdin
              | if . == [] then [null] end
              )
            ),
            join_string: (
              if $combined_opts.join_output then ""
              elif $combined_opts.null_output then "\u0000"
              else null
              end
            ),
            null_input: (
              ( if $combined_opts.expr_file then $rest
                else $rest[1:]
                end
              | if . == [] and $combined_opts.repl then true
                else null
                end
              )
            ),
            raw_file: (
              ( $combined_opts.raw_file
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
              if $combined_opts.raw_string
                or $combined_opts.join_output
                or $combined_opts.null_output
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
      , "Example usages:"
      , "  fq . file"
      , "  fq d file"
      , "  fq tovalue file"
      , "  cat file.cbor | fq -d cbor torepr"
      , "  fq 'grep(\"^main$\") | parent' /bin/ls"
      , "  fq 'grep_by(format == \"exif\") | d' *.png *.jpeg"
      , ""
      , args_help_text(_opt_cli_opts)
      ) | println
    elif $opts.show_version then
      "\($version) (\($os) \($arch))" | println
    elif $opts.show_formats then
      _formats_list | println
    elif
      ( $opts.filenames == [null] and
        $opts.null_input == false and
        ($opts.repl | not) and
        ($opts.expr_file | not) and
        stdin_tty.is_terminal and
        stdout_tty.is_terminal
      ) then
      ( (_usage($arg0) | printerrln)
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
