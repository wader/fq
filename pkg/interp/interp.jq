include "internal";
include "options";
include "binary";
include "decode";
include "match";
include "funcs";
include "grep";
include "args";
include "eval";
include "query";
include "repl";
include "help";
# generate torepr, format decode helpers and include format specific functions
include "formats";
# optional user init
include "@config/init?";


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
        | (_error_str([$name]) | printerrln)
        , _input($opts; f)
        )
    | try f
      catch
        ( . as $err
        | _input_decode_errors(. += {($name): $err}) as $_
        | [ $opts.decode_format
          , if $err | type == "string" then ": \($err)"
            # TODO: if not string assume decode itself failed for now
            else ": failed to decode (try -d FORMAT)"
            end
          ] | join("")
        | (_error_str([$name]) | printerrln)
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

# user expr error, report and continue
def _cli_eval_on_expr_error:
  ( if type == "object" then
      if .error | _eval_is_compile_error then .error | _eval_compile_error_tostring
      elif .error then .error
      end
    else tostring
    end
  | . as $err
  | _cli_last_expr_error($err) as $_
  | (_error_str([input_filename // empty]) | printerrln)
  );
# other expr error, should not happen, report and halt
def _cli_eval_on_error:
  halt_error(_exit_code_expr_error);
# could not compile expr, report and halt
def _cli_eval_on_compile_error:
  ( .error
  | _eval_compile_error_tostring
  | halt_error(_exit_code_compile_error)
  );
def _cli_repl_error($_):
  _eval_error("compile"; "repl can only be used from interactive repl");
def _cli_slurp_error(_):
  _eval_error("compile"; "slurp can only be used from interactive repl");
# _cli_eval halts on compile errors
def _cli_eval($expr; $opts):
  eval(
    $expr;
    $opts + {
      slurps: {
        help: "_help_slurp",
        repl: "_cli_repl_error",
        slurp: "_cli_slurp_error"
      },
      catch_query: _query_func("_cli_eval_on_expr_error")
    };
    _cli_eval_on_error;
    _cli_eval_on_compile_error
  );


def _main:
  def _banner:
    ( "fq - jq for binary formats"
    , "Tool, language and decoders for inspecting binary data."
    , "For more information see https://github.com/wader/fq"
    );
  def _usage($arg0):
    "Usage: \($arg0) [OPTIONS] [--] [EXPR] [FILE...]";
  def _help($arg0):
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
    );
  def _formats_list:
    ( [ formats
      | to_entries[]
      | [(.key+"  "), .value.description]
      ]
    | table(
        .;
        map(
          ( . as $rc
          # right pad format name to align description
          | if .column == 0 then .string | rpad(" "; $rc.maxwidth)
            else $rc.string
            end
          )
        ) | join("")
      )
    );
  def _map_decode_file:
    map(
      ( . as $a
      | .[1] |=
        try (open | decode)
        catch
          ( "--decode-file \($a[0]): \(.)"
          | halt_error(_exit_code_args_error)
          )
      )
    );
  ( . as {$version, $os, $arch, $args, args: [$arg0]}
  # make sure we don't unintentionally use . to make things clearer
  | null
  | ( try _args_parse($args[1:]; _opt_cli_opts)
      catch halt_error(_exit_code_args_error)
    ) as {parsed: $parsed_args, $rest}
  # combine default fixed opt, parsed args and -o key=value opts
  | _options_stack([
      ( ( _opt_build_default_fixed
        + $parsed_args
        + ($parsed_args.option | _opt_cli_arg_options)
        )
      | . + _opt_eval($rest)
      )
    ]) as $_
  | _opt_build_default_fixed as $default_fixed_opts
  # combine default fixed opt, --args opts and -o key=value opts
  | ( $default_fixed_opts
    + $parsed_args
    + ($parsed_args.option | _opt_cli_arg_options)
    ) as $combined_opts
  | options as $opts
  | if $opts.show_help then _help($arg0) | println
    elif $opts.show_version then "\($version) (\($os) \($arch))" | println
    elif $opts.show_formats then _formats_list | println
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
      ( # store some global state
        ( _include_paths($opts.include_path) as $_
        | _input_filenames($opts.filenames) as $_
        | _slurps(
            ( $opts.arg +
              $opts.argjson +
              $opts.raw_file +
              ($opts.decode_file | if . then _map_decode_file end)
            | map({key: .[0], value: .[1]})
            | from_entries
            )
          )
        ) as $_
      | { filename: $opts.expr_eval_path
        } as $eval_opts
      # use _finally as display etc prints and outputs empty
      | _finally(
        if $opts.repl then
          # TODO: share input_query but first have to figure out how to handle
          # context/interrupts better as open will happen in a sub repl which
          # context will be cancelled.
          ( def _inputs:
              if $opts.null_input then null
              elif $opts.string_input then inputs
              elif $opts.slurp then [inputs]
              else inputs
              end;
            [_inputs]
          | map(_cli_eval($opts.expr; $eval_opts))
          | _repl({})
          )
        else
          ( _cli_last_expr_error(null) as $_
          | _display_default_opts as $default_opts
          | _cli_eval(
              $opts.expr;
              ( $eval_opts
              | .input_query =
                  ( if $opts.null_input then _query_null
                    # note that jq --slurp --raw-input (string_input) is special, will concat
                    # all files into one string instead of iterating lines
                    elif $opts.string_input then _query_func("inputs")
                    elif $opts.slurp then _query_func("inputs") | _query_array
                    else _query_func("inputs")
                    end
                  )
              )
            )
          | display($default_opts)
          )
        end;
        # finally
        ( if _input_io_errors then null | halt_error(_exit_code_input_io_error) end
        | if _input_decode_errors then null | halt_error(_exit_code_input_decode_error) end
        | if _cli_last_expr_error then null | halt_error(_exit_code_expr_error) end
        )
      )
    )
    end
  );
