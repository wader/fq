def print: stdout;
def println: ., "\n" | stdout;
def debug:
  ( ((["DEBUG", .] | tojson), "\n" | stderr)
  , .
  );

# eval f and finally eval fin even on empty or error
def finally(f; fin):
  ( try f // (fin | empty)
    catch (fin as $_ | error)
  | fin as $_
  | .
  );

def _error_str: "error: \(.)";
def _errorln: ., "\n" | stderr;

def _global_var($k): _global_state[$k];
def _global_var($k; f): _global_state(_global_state | .[$k] |= f);

def _include_paths: _global_var("include_paths");
def _include_paths(f): _global_var("include_paths"; f);

def _default_options: _global_var("default_options");
def _default_options(f): _global_var("default_options"; f);

def _options_stack: _global_var("options_stack");
def _options_stack(f): _global_var("options_stack"; f);

def _parsed_args: _global_var("parsed_args");
def _parsed_args(f): _global_var("parsed_args"; f);

def _cli_last_expr_error: _global_var("cli_last_expr_error");
def _cli_last_expr_error(f): _global_var("cli_last_expr_error"; f);

def _input_filename: _global_var("input_filename");
def _input_filename(f): _global_var("input_filename"; f);

def _input_filenames: _global_var("input_filenames");
def _input_filenames(f): _global_var("input_filenames"; f);

def _input_io_errors: _global_var("input_io_errors");
def _input_io_errors(f): _global_var("input_io_errors"; f);

def _input_decode_errors: _global_var("input_decode_errors");
def _input_decode_errors(f): _global_var("input_decode_errors"; f);

def _variables: _global_var("variables");
def _variables(f): _global_var("variables"; f);
