include "@builtin/internal";
include "@builtin/funcs";
include "@builtin/args";

# will include all per format specific function etc
include "@format/all";

# optional user init
include "@config/init?";

# def read: #:: [a]| => string
# read with no prompt or completion

# def read($promp): #:: [a]|(string) => string
# read with prompt and no completion

# def read($promp; $completion): #:: [a]|(string;string) => string
# First argument is prompt to use.
# Second argument is name of completion function [a](string) => [string],
# it will be called with same input as read and a string argument being the
# current line from start to current cursor position. Should return possible completions.

# TODO: completionMode
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
        | map(select(strings and startswith($prefix)))
        | unique
        | sort
        )
      )
    }
  );

def _obj_to_csv_kv:
  [to_entries[] | [.key, .value] | join("=")] | join(",");

def _color_themes:
  {
    default: {
      colors: ({
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
      } | _obj_to_csv_kv),
      bytecolors: "0-0xff=brightwhite,0=brightblack,32-126:9-13=white",
    },
    # TODO: more configurable? colors=neon?
    neon: {
      colors: ({
        null: "brightblack",
        false: "brightyellow",
        true: "brightyellow",
        number: "brightcyan",
        string: "brightgreen",
        objectkey: "brightblue",
        array: "brightwhite",
        object: "brightwhite",
        index: "brightwhite",
        value: "brightwhite",
        error: "brightred",
        dumpheader: "brightyellow+underline",
        dumpaddr: "brightyellow"
      } | _obj_to_csv_kv),
      bytecolors: "0-0xff=brightwhite,0=brightblack,32-126:9-13=brightgreen",
    }
  };

def _build_default_options:
  {
    depth:        0,
    verbose:      false,
    color:        (tty.is_terminal and env.CLICOLOR!=null),
    unicode:      (tty.is_terminal and env.CLIUNICODE!=null),
    raw:          (tty.is_terminal | not),
    linebytes:    (if tty.is_terminal then [((tty.width div 8) div 2) * 2, 4] | max else 16 end),
    displaybytes: (if tty.is_terminal then [((tty.width div 8) div 2) * 2, 4] | max else 16 end),
    addrbase:     16,
    sizebase:     10,
    colors:       _color_themes.default.colors,
    bytecolors:   _color_themes.default.bytecolors,
  };

def _parse_options:
  ( {
      depth:        (.depth | if . then eval(.) else null end),
      verbose:      (.verbose | if . then eval(.) else null end),
      color:        (.color | if . then eval(.) else null end),
      unicode:      (.unicode | if . then eval(.) else null end),
      raw:          (.raw | if . then eval(.) else null end),
      linebytes:    (.linebytes | if . then eval(.) else null end),
      displaybytes: (.displaybytes | if . then eval(.) else null end),
      addrbase:     (.addrbase | if . then eval(.) else null end),
      sizebase:     (.sizebase | if . then eval(.) else null end),
      colors:       .colors,
      bytecolors:   .bytecolors,
    }
  | with_entries(select(.value != null))
  );

def _prompt(iter):
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
    if type != "array" then
      _type_name_error
    else
      ( "["
      , if length > 0 then (.[0] | _type_name_error) else empty end
      , if length > 1 then ", ..." else empty end
      , "]"
      , "[\(length)]"
      )
    end;
  ( [iter]
  | [ (options.repllevel | if . > 1 then ((.-1) * ">") + " " else empty end)
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

def _eval_debug:
  (["DEBUG", .] | tojson, "\n") | stderr;

def _eval_f($e; f):
  ( _default_options(_build_default_options) as $_
  | try eval($e; "eval_debug") | f
    catch (. as $err | ("error: " + ($err | tostring)) | println)
  );

def _default_display: display({depth: 1});

def _eval_print($e):
  _eval_f($e; _default_display);

# run read-eval-print-loop
def repl($opts; iter): #:: a|(Opts) => @
  def _read_expr: read(_prompt(iter); "_complete") | trim;
  def _repl:
    ( . as $c
    | try
        ( _read_expr as $e
        | if $e != "" then
            (iter | _eval_print($e))
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
  _with_options($opts | .repllevel = options.repllevel+1; _repl);
# same as repl({})
def repl($opts): repl($opts; .);
def repl: repl({}; .); #:: a| => @

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
        [.[] as $rc | $rc.string | rpad(" "; $rc.maxwidth)] | join("")
      );
  def _opts($version):
    {
      "version": {
        short: "-v",
        long: "--version",
        description: "Show version (\($version))",
        bool: true
      },
      "help": {
        short: "-h",
        long: "--help",
        description: "Show help",
        bool: true
      },
      "formats": {
        long: "--formats",
        description: "Show formats",
        bool: true
      },
      "nullinput": {
        short: "-n",
        description: "Null input",
        bool: true
      },
      "decode_format": {
        short: "-d",
        long: "--decode",
        description: "Decode format",
        default: "probe",
        string: true
      },
      "repl": {
        short: "-i",
        long: "--repl",
        description: "Interactive REPL",
        bool: true
      },
      "file": {
        short: "-f",
        long: "--file",
        description: "Read script from file",
        string: true
      },
      "rawstring": {
        short: "-r",
        description: "Raw strings",
        bool: true
      },
      "compact": {
        short: "-c",
        long: "--compact",
        description: "Compact output",
        bool: true
      },
      "options": {
        short: "-o",
        long: "--option",
        description: "Set option, eg: color=true",
        object: true,
        default: {},
        help_default: _build_default_options
      },
    };
  def _usage($arg0; $version):
    "Usage: \($arg0) [OPTIONS] [EXPR] [FILE...]";
  ( .version as $version
  | .args[0] as $arg0
  | args_parse(.args[1:]; _opts($version)) as {$parsed, $rest}
  # store parsed arguments, .format is used by input
  | _parsed_args($parsed) as $_
  | _default_options(_build_default_options) as $_
  # TODO: hack, pass opts some other way?
  | _push_options(
      ( ($parsed.options | _parse_options)
      + {
          repl: ($parsed.repl == true),
          rawstring: ($parsed.rawstring == true),
          compact: ($parsed.compact == true),
          repllevel: 0,
        }
      )
    )
  | if $parsed.help then
      ( _usage($arg0; $version)
      , args_help_text(_opts($version))
      ) | println
    elif $parsed.version then
      $version | println
    elif $parsed.formats then
      _formats_list | println
    elif ($rest | length) == 0 and (($parsed.repl | not) and ($parsed.file | not)) then
      _usage($arg0; $version) | println
    else
      try
        ( { nullinput: ($parsed.nullinput == true) }
        | if $parsed.file then
            ( .expr = ($parsed.file | open | string)
            | .filenames = $rest
            )
          else
            ( .expr = ($rest[0] // ".")
            | .filenames = $rest[1:]
            )
          end
        | if $parsed.repl and .filenames == [] then
            .nullinput = true
          elif .filenames == [] then
            .filenames = ["-"]
          end
        | . as {$expr, $filenames, $nullinput}
        | inputs($filenames) as $_ # store inputs
        | if $nullinput then null
          # TODO: exit codes on input error and expr error
          else inputs # will iterate inputs
          end
        | if $parsed.repl then [eval($expr)] | repl({}; .[])
          else
            try (eval($expr) | _default_display)
            catch (. as $err | ("error: ", ($err | tostring), "\n") | stderr)
          end
        )
      catch tostring | halt_error(1)
    end
  );
