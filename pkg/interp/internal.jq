# eval f and finally eval fin even on empty or error
def finally(f; fin):
  ( try f // (fin | empty)
    catch (fin as $_ | error(.))
  | fin as $_
  | .
  );

def _print_error: "error: \(.)" | println;
def _stderr_error: "error: \(.)\n" | stderr;

def _default_options: _eval_state("default_options");
def _default_options($opts): _eval_state("default_options"; $opts);

def _push_options($opts): _eval_state("options_stack"; [$opts] + (_eval_state("options_stack") // []));
def _pop_options: _eval_state("options_stack"; _eval_state("options_stack")[1:]);

def _with_options($opts; f):
  _push_options($opts) as $_ | finally(f; _pop_options);

def _parsed_args: _global_state("parsed_args");
def _parsed_args($v): _global_state("parsed_args"; $v);

def _cli_last_expr_error: _global_state("cli_last_expr_error");
def _cli_last_expr_error($v): _global_state("cli_last_expr_error"; $v);

# next valid input
def input:
  ( _global_state("inputs")
  | if length == 0 then error("break") end
  | [.[0], .[1:]] as [$h, $t]
  | _global_state("inputs"; $t)
  | _global_state("input_filename"; null) as $_
  | $h
  | try
      ( open
      | _global_state("input_filename"; $h) as $_
      | .
      )
    catch
      ( _global_state("input_io_errors";
          (_global_state("input_io_errors") // {}) + {($h): .}
        ) as $_
      | _stderr_error
      , input
      )
  | try
      decode(_parsed_args.decode_format)
    catch
      ( _global_state("input_decode_errors";
          (_global_state("input_decode_errors") // {}) + {($h): .}
        ) as $_
      | "\($h): failed to decode (\(_parsed_args.decode_format)), try -d FORMAT to force"
      | _stderr_error
      , input
      )
  );

# iterate all valid inputs
def inputs:
  try repeat(input)
  catch if . == "break" then empty else error end;
def inputs($v): _global_state("inputs"; $v);

def input_filename: _global_state("input_filename");

def _input_io_errors: _global_state("input_io_errors");
def _input_decode_errors: _global_state("input_decode_errors");
