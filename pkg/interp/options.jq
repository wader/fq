include "internal";
include "binary";


def _opt_build_default_fixed:
  ( stdout_tty as $stdout
  | { addrbase:       16
    , arg:            []
    , argdecode:      []
    , argjson:        []
    , array_truncate: 50
    , bits_format:    "string"
    # 0-0xff=brightwhite,0=brightblack,32-126:9-13=default
    , byte_colors:
        [ { ranges: [[0,255]]
          , value: "default+bold"
          }
        , { ranges: [[0]]
          , value: "brightblack"
          }
        , { ranges: [[32,126],[9,13]]
          , value: "default"
          }
        ]
    , color: ($stdout.is_terminal and (env.NO_COLOR | . == null or . == ""))
    , colors:
        { null:              "brightblack"
        , false:             "yellow"
        , true:              "yellow"
        , number:            "cyan"
        , string:            "green"
        , objectkey:         "brightblue"
        , array:             "default"
        , object:            "default"
        , index:             "default"
        , value:             "default"
        , error:             "brightred"
        , dumpheader:        "yellow+underline"
        , dumpaddr:          "yellow"
        , prompt_repl_level: "brightblack"
        , prompt_value:      "default"
        }
    , compact:            false
    , completion_timeout: (env.COMPLETION_TIMEOUT | if . != null then tonumber else 1 end)
    , decode_group:       "probe"
    , decode_progress:    (env.NO_DECODE_PROGRESS == null)
    , depth:              0
    , expr_eval_path:     "arg"
    , expr_file:          null
    , expr_given:         false
    , expr:               "."
    , filenames:          null
    , force:              false
    , include_path:       null
    , join_string:        "\n"
    , null_input:         false
    , raw_file:           []
    , raw_output:         ($stdout.is_terminal | not)
    , raw_string:         false
    , repl:               false
    , show_formats:       false
    , show_help:          false
    , sizebase:           10
    , skip_gaps:          false
    , slurp:              false
    , string_input:       false
    , string_truncate:    50
    , unicode:            ($stdout.is_terminal and env.CLIUNICODE != null)
    , value_output:       false
    , verbose:            false
    }
  );

def _opt_options:
  { addrbase:           "number"
  , arg:                "array_string_pair"
  , argdecode:          "array_string_pair"
  , argjson:            "array_string_pair"
  , array_truncate:     "number"
  , bits_format:        "string"
  , byte_colors:        "csv_ranges_array"
  , color:              "boolean"
  , colors:             "csv_kv_obj"
  , compact:            "boolean"
  , completion_timeout: "number"
  , decode_group:       "string"
  , decode_progress:    "boolean"
  , depth:              "number"
  , display_bytes:      "number"
  , expr_eval_path:     "string"
  , expr_file:          "string"
  , expr_given:         "boolean"
  , expr:               "string"
  , filenames:          "array_string"
  , force:              "boolean"
  , include_path:       "string"
  , join_string:        "string"
  , line_bytes:         "number"
  , null_input:         "boolean"
  , raw_file:           "array_string_pair"
  , raw_output:         "boolean"
  , raw_string:         "boolean"
  , repl:               "boolean"
  , show_formats:       "boolean"
  , show_help:          "boolean"
  , sizebase:           "number"
  , skip_gaps:          "boolean"
  , slurp:              "boolean"
  , string_input:       "boolean"
  , string_truncate:    "number"
  , unicode:            "boolean"
  , value_output:       "boolean"
  , verbose:            "boolean"
  , width:              "number"
  };

def _opt_eval($rest):
  ( with_entries(
      ( select(.value | _is_string and startswith("@"))
      | .key as $opt
      | .value |=
          ( . as $v
          | try
              ( .[1:]
              | open
              | tobytes
              | tostring
              )
            catch
              ( "-o \($opt)=@\($v[1:]): \(.)"
              | _fatal_error(_exit_code_args_error)
              )
          )
      )
    )
  + { argjson: (
        ( .argjson
        | if . then
            map(
              ( . as $a
              | .[1] |=
                try fromjson
                catch
                  ( "--argjson \($a[0]): \(.)"
                  | _fatal_error(_exit_code_args_error)
                  )
              )
            )
          end
        )
      )
    , color: (
        if .monochrome_output == true then false
        elif .color_output == true then true
        else null
        end
      )
    , expr: (
        # if -f was used, all rest non-args are filenames
        # otherwise first is expr rest is filenames
        ( .expr_file
        | . as $expr_file
        | if . then
            try (open | tobytes | tostring)
            catch ("\($expr_file): \(.)" | _fatal_error(_exit_code_args_error))
          else $rest[0] // null
          end
        )
      )
    , expr_given: (
        # was a expr arg given
        $rest[0] != null
      )
    , expr_eval_path: .expr_file
    , filenames: (
        ( if .filenames then .filenames
          elif .expr_file then $rest
          else $rest[1:]
          end
        # null means stdin
        | if . == [] then [null] end
        )
      )
    , join_string: (
        if .join_output then ""
        elif .null_output then "\u0000"
        else null
        end
      )
    , null_input: (
        ( ( if .expr_file then $rest
            else $rest[1:]
            end
          ) as $files
        | if $files == [] and .repl then true
          else null
          end
        )
      )
    , raw_file: (
        ( .raw_file
        | if . then
            ( map(.[1] |=
                ( . as $f
                | try (open | tobytes | tostring)
                  catch ("\($f): \(.)" | _fatal_error(_exit_code_args_error))
                )
              )
            )
          end
        )
      )
    , raw_string: (
        if .raw_string
          or .join_output
          or .null_output
        then true
        else null
        end
      )
    , unicode: (
        if .unicode_output == true then true
        else null
        end
      )
    , value_output: (
        if .value_output == true then true
        else null
        end
      ),
    }
  | with_entries(select(.value != null))
  );

# these _to* function do a bit for fuzzy string to type conversions
def _opt_to_boolean:
  try
    if . == "true" then true
    elif . == "false" then false
    else tonumber != 0
    end
  catch
    null;

def _opt_from_boolean: tostring;

def _opt_to_number:
  try tonumber catch null;

def _opt_from_number: tostring;

def _opt_to_string:
  if . != null then
    ( "\"\(.)\""
    | try
        ( fromjson
        | if type != "string" then error end
        )
      catch null
    )
  end;

def _opt_from_string: if . then tojson[1:-1] else "" end;

def _opt_is_string_pair:
  _is_array and length == 2 and all(_is_string);

def _opt_to_array(f):
  try
    ( fromjson
    | if _is_array and (all(f) | not) then null end
    )
  catch null;

def _opt_to_array_string_pair: _opt_to_array(_opt_is_string_pair);
def _opt_to_array_string: _opt_to_array(_is_string);

def _opt_from_array: tojson;

# TODO: cleanup
def _trim: capture("^\\s*(?<str>.*?)\\s*$"; "").str;

# "0-255=brightwhite,0=brightblack,32-126:9-13=default" -> [{"ranges": [[0-255]], value: "brightwhite"}, ...]
def _csv_ranges_to_array:
  ( split(",")
  | map(
    ( _trim
    | split("=")
    | { ranges:
          ( .[0]
          | split(":")
          | map(split("-") | map(tonumber))
          )
      , value: .[1]
      }
    ))
  );

def _opt_to_csv_ranges_array:
  try _csv_ranges_to_array
  catch null;

def _opt_from_csv_ranges_array:
  ( map(
      ( (.ranges | map(join("-")) | join(":"))
      + "="
      + .value
      )
    )
  | join(",")
  );

# "key=value,a=b,..." -> {"key": "value", "a": "b", ...}
def _csv_kv_to_obj:
  ( split(",")
  | map(_trim | split("=") | {key: .[0], value: .[1]})
  | from_entries
  );

def _opt_to_csv_kv_obj:
  try _csv_kv_to_obj
  catch null;

def _opt_from_csv_kv_obj:
  ( to_entries
  | map("\(.key)=\(.value)")
  | join(",")
  );

def _opt_to_fuzzy:
  ( . as $s
  | try fromjson
    catch
      ( $s
      | _opt_to_string
      // $s
      )
  );

def _opt_to($type):
  if $type == "array_string" then _opt_to_array_string
  elif $type == "array_string_pair" then _opt_to_array_string_pair
  elif $type == "boolean" then _opt_to_boolean
  elif $type == "csv_kv_obj" then _opt_to_csv_kv_obj
  elif $type == "csv_ranges_array" then _opt_to_csv_ranges_array
  elif $type == "number" then _opt_to_number
  elif $type == "string" then _opt_to_string
  elif $type == "fuzzy" then _opt_to_fuzzy
  else error("unknown type \($type)")
  end;

def _opt_from($type):
  if $type == "array_string" then _opt_from_array
  elif $type == "array_string_pair" then _opt_from_array
  elif $type == "boolean" then _opt_from_boolean
  elif $type == "csv_kv_obj" then _opt_from_csv_kv_obj
  elif $type == "csv_ranges_array" then _opt_from_csv_ranges_array
  elif $type == "number" then _opt_from_number
  elif $type == "string" then _opt_from_string
  else error("unknown type \($type)")
  end;

def _opt_cli_arg_to_options:
  ( _opt_options as $opts
  | with_entries(
      ( .key as $k
      | .value |= _opt_to($opts[$k] // "fuzzy")
      | select(.value != null)
      )
    )
  );

def _opt_cli_arg_from_options:
  ( _opt_options as $opts
  | with_entries(
      ( .key as $k
      | .value |= _opt_from($opts[$k] // "string")
      | select(.value != null)
      )
    )
  );

def _opt_cli_opts:
  { arg:
      { long: "--arg"
      , description: "Set variable $NAME to string VALUE"
      , pairs: "NAME VALUE"
      }
  , argdecode:
      { long: "--argdecode"
      # TODO: remove at some point
      , aliases: ["--decode-file"]
      , description: "Set variable $NAME to decode of PATH"
      , pairs: "NAME PATH"
      }
  , argjson:
      { long: "--argjson"
      , description: "Set variable $NAME to JSON"
      , pairs: "NAME JSON"
      }
  , compact:
      { short: "-c"
      , long: "--compact-output"
      , description: "Compact output"
      , bool: true
      }
  , color_output:
      { short: "-C"
      , long: "--color-output"
      , description: "Force color output"
      , bool: true
      }
  , decode_group:
      { short: "-d"
      , long: "--decode"
      , description: "Decode format or group (probe)"
      , string: "NAME"
      }
  , expr_file:
      { short: "-f"
      , long: "--from-file"
      , description: "Read EXPR from file"
      , string: "PATH"
      }
  , show_help:
      { short: "-h"
      , long: "--help"
      , description: "Show help for TOPIC (ex: -h formats, -h mp4)"
      , string: "[TOPIC]"
      , optional: true
      }
  , join_output:
      { short: "-j"
      , long: "--join-output"
      , description: "No newline after each output"
      , bool: true
      }
  , include_path:
      { short: "-L"
      , long: "--include-path"
      , description: "Include search path"
      , array: "PATH"
      }
  , null_output:
      { long: "--raw-output0"
      # for jq compatibility
      , aliases: ["--nul-output"]
      , description: "NUL (zero) byte after each output"
      , bool: true
      }
  , null_input:
      { short: "-n"
      , long: "--null-input"
      , description: "Null input (use input and inputs functions to read)"
      , bool: true
      }
  , monochrome_output:
      { short: "-M"
      , long: "--monochrome-output"
      , description: "Force monochrome output"
      , bool: true
      }
  , option:
      { short: "-o"
      , long: "--option"
      , description: "Set option (ex: -o color=true, see --help options)"
      , object: "KEY=VALUE/@PATH",
      }
  , string_input:
      { short: "-R"
      , long: "--raw-input"
      , description: "Read raw input strings (don't decode)"
      , bool: true
      }
  , raw_file:
      { long: "--raw-file"
      # for jq compatibility
      , aliases: ["--raw-file"]
      , description: "Set variable $NAME to string content of file"
      , pairs: "NAME PATH"
      }
  , raw_string:
      { short: "-r"
      # for jq compat, is called raw string internally, is different from "raw output" which
      # is if we can output raw bytes or not
      , long: "--raw-output"
      , description: "Raw string output (without quotes)"
      , bool: true
      }
  , repl:
      { short: "-i"
      , long: "--repl"
      , description: "Interactive REPL"
      , bool: true
      }
  , slurp:
      { short: "-s"
      ,  long: "--slurp"
      ,  description: "Slurp all inputs into an array or string (-Rs)"
      ,  bool: true
      }
  , unicode_output:
      { short: "-U"
      , long: "--unicode-output"
      , description: "Force unicode output"
      , bool: true
      }
  , value_output:
      { short: "-V"
      , long: "--value-output"
      , description: "Output JSON value (-Vr for raw string)"
      , bool: true
      }
  , show_version:
      { short: "-v"
      , long: "--version"
      , description: "Show version"
      , bool: true
      },
  };

def options($opts):
  ( stdout_tty as $stdout
  | ( [{width: $stdout.width}]
    + _options_stack
    + [$opts]
    )
  | add
  | ( if .width != 0 then [_intdiv(_intdiv(.width; 8); 2) * 2, 4] | max
      else 16
      end
    ) as $display_bytes
  # default if not set
  | .display_bytes |= (. // $display_bytes)
  | .line_bytes |= (. // $display_bytes)
  );
def options: options({});
