include "internal";
include "options";
include "binary";
include "decode";
include "registry_include";
include "format_decode";
include "format_func";
include "grep";
include "args";
include "eval";
include "query";
include "interp";
include "repl";
include "help";
include "funcs";
# optional user init
include "@config/init?";

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
        | [ $opts.decode_group
          , if $err | _is_string then ": \($err)"
            # TODO: if not string assume decode itself failed for now
            else ": failed to decode: try fq -d FORMAT to force format, see fq -h formats for list"
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
  ( if _is_object then
      if .error | _eval_is_compile_error then .error | _eval_compile_error_tostring
      elif .error then .error
      end
    else tostring
    end
  | . as $err
  | _cli_last_expr_error($err) as $_
  | (_error_str([input_filename // empty]) | printerrln)
  );
# other expr error, other errors then cancel should not happen, report and halt
def _cli_eval_on_error:
  if .error | _is_context_canceled_error then (null | halt_error(_exit_code_expr_error))
  else _fatal_error(_exit_code_expr_error)
  end;
# could not compile expr, report and halt
def _cli_eval_on_compile_error:
  ( .error
  | _eval_compile_error_tostring
  | _fatal_error(_exit_code_compile_error)
  );
def _cli_repl_error($_):
  _eval_error("compile"; "repl can only be used from interactive repl");
def _cli_slurp_error(_):
  _eval_error("compile"; "slurp can only be used from interactive repl");
# TODO: rewrite query to reuse _display_default_opts value? also _repl_display
def _cli_display:
  display_implicit(_display_default_opts);
# _cli_eval halts on compile errors
def _cli_eval($expr; $opts):
  eval(
    $expr;
    ( $opts
    + { slurps:
          { help: "_help_slurp"
          , repl: "_cli_repl_error"
          , slurp: "_cli_slurp_error"
          }
      , catch_query: _query_func("_cli_eval_on_expr_error"),
      }
    );
    _cli_eval_on_error;
    _cli_eval_on_compile_error
  );


def _main:
  def _map_argdecode:
    map(
      ( . as $a
      | .[1] |=
        try (open | decode)
        catch
          ( "--argdecode \($a[0]): \(.)"
          | _fatal_error(_exit_code_args_error)
          )
      )
    );
  ( . as {$version, $os, $arch, $go_version, $args, args: [$arg0]}
  # make sure we don't unintentionally use . to make things clearer
  | null
  | ( try _args_parse($args[1:]; _opt_cli_opts)
      catch _fatal_error(_exit_code_args_error)
    ) as {parsed: $parsed_args, $rest}
  # combine default fixed opt, parsed args and -o key=value opts
  | _options_stack([
      ( ( _opt_build_default_fixed
        + $parsed_args
        + ($parsed_args.option | if . then _opt_cli_arg_to_options end)
        )
      | . + _opt_eval($rest)
      )
    ]) as $_
  | options as $opts
  | if $opts.show_help then
      ( # if show_help is a string -h <topic> was used
        if ($opts.show_help | type) == "boolean" then
          ( # "" to print separators
            ( "banner"
            , ""
            , "usage"
            , ""
            , "example_usage"
            , ""
            , "args"
            )
          | if . != "" then _help($arg0; .) end
          )
        else _help($arg0; $opts.show_help)
        end
      | println
      )
    elif $opts.show_version then
      "\($version) (\($os) \($arch) \($go_version))" | println
    elif
      ( $opts.filenames == [null] and
        $opts.null_input == false and
        ($opts.repl | not) and
        ($opts.expr_file | not) and
        ($opts.expr_given | not) and
        stdin_tty.is_terminal and
        stdout_tty.is_terminal
      ) then
      ( (_help($arg0; "usage") | printerrln)
      , (null | halt_error(_exit_code_args_error))
      )
    else
      ( # store some global state
        ( _include_paths($opts.include_path) as $_
        | _input_filenames($opts.filenames) as $_
        | _slurps(
            ( $opts.arg +
              $opts.argjson +
              $opts.raw_file +
              ($opts.argdecode | if . then _map_argdecode end)
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
              # call display in sub eval so it can be interrupted
              # for repl case value will used as input to _repl instead
              | .output_query = _query_func("_cli_display")
              )
            )
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
