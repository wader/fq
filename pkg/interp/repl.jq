include "internal";
include "options";
include "eval";
include "query";
include "decode";
include "interp";
include "funcs";
include "ansi";

# TODO: currently only make sense to allow keywords starting a term or directive
def _complete_keywords:
  [ "and"
  #"as"
  #"break"
  #"catch"
  , "def"
  #"elif"
  #"else"
  #"end"
  , "false"
  , "foreach"
  , "if"
  , "import"
  , "include"
  , "label"
  , "module"
  , "null"
  , "or"
  , "reduce"
  #"then"
  , "true"
  , "try"
  ];

def _complete_scope:
  [ (scope | map(split("/")[0]) | unique)
  , _complete_keywords
  ] | add;
def _complete_keys:
  # uses try as []? will not catch errors
  [try keys[] catch empty, try _extkeys[] catch empty];

# TODO: "def zzz: null; abc" scope won't see zzz after completion rewrite as def will be inside map(def zzz: null; abc | ...)
# TODO: handle variables via ast walk?
# TODO: refactor this
# TODO: completionMode
# TODO: return escaped identifier, not sure current readline implementation supports
# modifying "previous" characters if quoting is needed
# completions that needs to change previous input, ex: .a\t -> ."a \" b" etc
def _complete($line; $cursor_pos):
  # TODO: reverse this? word or non-ident char?
  def _is_separator: . as $c | " .;[]()|=" | contains($c);
  def _is_internal: startswith("_") or startswith("$_");
  def _query_index_or_key($q):
    ( ([.[] | _eval($q; {}) | type]) as $n
    | if ($n | all(. == "object")) then "."
      elif ($n | all(. == "array")) then "[]"
      else null
      end
    );
  # only complete if at end or there is a whitespace for now
  if ($line[$cursor_pos] | . == null or _is_separator) then
    ( . as $c
    | $line[0:$cursor_pos]
    | . as $line_query
    # expr -> map(partial-expr | . | f?) | add
    # TODO: move map/add logic to here?
    | _query_completion(
        if .type | . == "func" or . == "var" then "_complete_scope"
        elif .type == "index" then "_complete_keys"
        else error("unreachable")
        end
      ) as {$type, $query, $prefix}
    | { prefix: $prefix
      , names: (
          if $type == "none" then
            ( $c
            | _query_index_or_key($line_query)
            | if . then [.] else [] end
            )
          else
            ( $c
            | _eval($query; {})
            | ($prefix | _is_internal) as $prefix_is_internal
            | map(
                select(
                  strings and
                  # TODO: var type really needed? just func?
                  (_is_ident or $type == "var") and
                  ((_is_internal | not) or $prefix_is_internal) and
                  startswith($prefix)
                )
              )
            | unique
            | sort
            | if length == 1 and .[0] == $prefix then
                ( $c
                | _query_index_or_key($line_query)
                | if . then [$prefix+.] else [$prefix] end
                )
              end
            )
          end
        )
      }
    )
  else
    {prefix: "", names: []}
  end;
def _complete($line): _complete($line; $line | length);

# empty input []
# >* empty>
# single input [v]
# >* VALUE_PATH VALUE_PREVIEW>
# multiple inputs [v,...]
# >* VALUE_PATH VALUE_PREVIEW, ...[#]>
# single/multi inputs where first input is array [[v,...], ...]
# >* [VALUE_PATH VALUE_PREVIEW, ...][#], ...[#]>
def _prompt($opts):
  def _repl_level:
    (_options_stack | length | if . > 2 then ((.-2) * ">") else empty end);
  def _value_path:
    (._path? // []) | if . == [] then empty else _path_to_expr($opts) end;
  def _value_preview($depth):
    if $depth == 0 and format == null and _is_array then
      [ "["
      , if length == 0 then empty
        else
          ( (.[0] | _value_preview(1))
          , if length > 1  then ", ..." else empty end
          )
        end
      , "]"
      , if length > 1 then
          ( ("[" | _ansi_if($opts; "array"))
          , ("0" | _ansi_if($opts; "number"))
          , ":"
          , (length | tostring | _ansi_if($opts; "number"))
          , ("]" | _ansi_if($opts; "array"))
          )
        else empty
        end
      ] | join("")
    else
      ( . as $c
      | format
      | if . != null then
          ( .
          + if $c._error then "!" else "" end
          )
        else
          ( $c
          | if _is_decode_value then type
            else (_exttype // type)
            end
          )
        end
      ) | _ansi_if($opts; "prompt_value")
    end;
  def _value:
    [ _value_path
    , _value_preview(0)
    ] | join(" ");
  def _values:
    if length == 0 then "empty"
    else
      [ (.[0] | _value)
      , if length > 1 then
          ( ", ..."
          , ("[" | _ansi_if($opts; "array"))
          , ("0" |  _ansi_if($opts; "number"))
          , ":"
          , (length | tostring | _ansi_if($opts; "number"))
          , ("]" | _ansi_if($opts; "array"))
          , "[]"
          )
        else empty
        end
      ] | join("")
    end;
  [ (_repl_level | _ansi_if($opts; "prompt_repl_level"))  , _values
  ] | join(" ") + "> ";
def _prompt: _prompt(null);

# user expr error
def _repl_on_expr_error:
  ( if _eval_is_compile_error then _eval_compile_error_tostring
    else tostring
    end
  | _error_str
  | println
  );
# other expr error, interrupted or something unexpected happened
def _repl_on_error:
  # was interrupted by user, just ignore
  if .error | _is_context_canceled_error then empty
  else _fatal_error(_exit_code_expr_error)
  end;
# compile error
def _repl_on_compile_error:
  ( if .error | _eval_is_compile_error then
      ( # TODO: move, redo as: def _symbols: if unicode then {...} else {...} end?
        def _arrow_up: if options.unicode then "⬆" else "^" end;
        if .error.column != 0 then
          ( ((.input | _prompt | length) + .error.column-1) as $pos
          | " " * $pos + "\(_arrow_up) \(.error.error)"
          )
        else
          ( .error
          | _eval_compile_error_tostring
          | _error_str
          )
        end
      )
    else .error | _error_str
    end
  | println
  );
def _repl_display:
  display(_display_default_opts);
def _repl_eval($expr; on_error; on_compile_error):
  eval(
    $expr;
    { slurps:
        { repl: "_repl_slurp"
        , help: "_help_slurp"
        , slurp: "_slurp"
        }
      # input to repl is always array of values to iterate
    , input_query: (_query_ident | _query_iter) # .[]
      # each input should be evaluated separately like cli file args, so catch and just print errors
    , catch_query: _query_func("_repl_on_expr_error")
      # run display in sub eval so it can be interrupted
    , output_query: _query_func("_repl_display")
    };
    on_error;
    on_compile_error
  );

# run read-eval-print-loop
# input is array of inputs to iterate
def _repl($opts):
  def _read_expr:
    _repeat_break(
      # both _prompt and _complete want input arrays
      ( _readline(
          { prompt: _prompt(options($opts))
          , complete: "_complete"
          , timeout: options.completion_timeout
          }
        )
      | if trim == "" then empty
        else (., error("break"))
        end
      )
    );
  def _repl_loop:
    try
      _repl_eval(
        _read_expr;
        _repl_on_error;
        _repl_on_compile_error
      )
    catch
      if . == "interrupt" then empty
      elif . == "eof" then error("break")
      elif _eval_is_compile_error then _repl_on_error
      else error
      end;
  if $opts | type != "object" then
    error("options must be an object")
  elif _is_completing | not then
    ( _options_stack(. + [$opts]) as $_
    | _finally(
        _repeat_break(_repl_loop);
        _options_stack(.[:-1])
      )
    )
  else empty
  end;

def _repl_slurp_eval($query):
  try
    [ eval(
        $query | _query_tostring;
        {};
        _repl_on_expr_error;
        error
      )
    ]
  catch
    error(.error);

def _repl_slurp($query):
  if ($query.slurp_args | length) > 1 then
    _eval_error("compile"; "repl requires none or one options argument. ex: ... | repl or ... | repl({compact: true})")
  else
    # only allow one output for args, multiple would be confusing i think (would start multiples repl:s)
    ( ( if ($query.slurp_args | length) > 0 then
          first(_repl_slurp_eval($query.slurp_args[0])[])
        else {}
        end
      ) as $opts
    | if $opts | type != "object" then
        _eval_error("compile"; "options must be an object")
      end
    | _repl_slurp_eval($query.rewrite)
    | _repl($opts)
    )
  end;

# just gives error, call appearing last will be renamed to _repl_slurp
def repl($_): error("repl must be last in pipeline. ex: ... | repl");
def repl: repl(null);

def _slurp($query):
  if ($query.slurp_args | length != 1) then
    _eval_error("compile"; "slurp requires one string argument. ex: ... | slurp(\"name\")")
  else
    # TODO: allow only one output?
    ( _repl_slurp_eval($query.slurp_args[0])[] as $name
    | if ($name | _is_ident | not) then
        _eval_error("compile"; "invalid slurp name \"\($name)\", must be a valid identifier. ex: ... | slurp(\"name\")")
      else
        ( _repl_slurp_eval($query.rewrite) as $v
        | _slurps(.[$name] |= $v)
        | empty
        )
      end
    )
  end;

def slurp($_): error("slurp must be last in pipeline. ex: ... | slurp(\"name\")");
def slurp: slurp(null);

def spew($name):
  ( _slurps[$name]
  | if . then .[]
    else error("no such slurp: \($name)")
    end
  );
def spew:
  _slurps;
