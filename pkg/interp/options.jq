def _obj_to_csv_kv:
  [to_entries[] | [.key, .value] | join("=")] | join(",");

def _build_default_fixed_options:
  ( (null | stdout) as $stdout
  | {
      addrbase:       16,
      arg:            [],
      argjson:        [],
      array_truncate: 50,
      bits_format:    "snippet",
      byte_colors:    "0-0xff=brightwhite,0=brightblack,32-126:9-13=white",
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
      expr:            ".",
      expr_file:       null,
      expr_eval_path:  "arg",
      filenames:       ["-"],
      include_path:    null,
      join_string:     "\n",
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

def _build_default_dynamic_options:
  ( (null | stdout) as $stdout
  | {
      # TODO: intdiv 2 * 2 to get even number, nice or maybe not needed?
      display_bytes:   (if $stdout.is_terminal then [intdiv(intdiv($stdout.width; 8); 2) * 2, 4] | max else 16 end),
      line_bytes:      (if $stdout.is_terminal then [intdiv(intdiv($stdout.width; 8); 2) * 2, 4] | max else 16 end),
    }
  );

# these _to* function do a bit for fuzzy string to type conversions
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
      array_truncate:  (.array_truncate | _tonumber),
      bits_format:     (.bits_format | _tostring),
      byte_colors:     (.byte_colors | _tostring),
      color:           (.color | _toboolean),
      colors:          (.colors | _tostring),
      compact:         (.compact | _toboolean),
      decode_file:     (.decode_file | _toarray(type == "string")),
      decode_format:   (.decode_format | _tostring),
      decode_progress: (.decode_progress | _toboolean),
      depth:           (.depth | _tonumber),
      display_bytes:   (.display_bytes | _tonumber),
      expr:            (.expr | _tostring),
      expr_file:       (.expr_file | _tostring),
      filename:        (.filenames | _toarray(type == "string")),
      include_path:    (.include_path | _tostring),
      join_string:     (.join_string | _tostring),
      line_bytes:      (.line_bytes | _tonumber),
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

# . will have additional array of options taking priority
# NOTE: is called from go *interp.Interp Options()
def options($opts):
  [_build_default_dynamic_options] + _options_stack + $opts | add;
def options: options([{}]);
