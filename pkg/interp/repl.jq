include "query";

# TODO: currently only make sense to allow keywords start  start a term or directive
def _complete_keywords:
  [
    "and",
    #"as",
    #"break",
    #"catch",
    "def",
    #"elif",
    #"else",
    #"end",
    "false",
    "foreach",
    "if",
    "import",
    "include",
    "label",
    "module",
    "null",
    "or",
    "reduce",
    #"then",
    "true",
    "try"
  ];

def _complete_scope:
  [scope[], _complete_keywords[]];

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
    ( ([.[] | eval($q) | type]) as $n
    | if ($n | all(. == "object")) then "."
      elif ($n | all(. == "array")) then "[]"
      else null
      end
    );
  # only complete if at end or there is a whitespace for now
  if ($line[$cursor_pos] | . == "" or _is_separator) then
    ( . as $c
    | $line[0:$cursor_pos]
    | . as $line_query
    # expr -> map(partial-expr | . | f?) | add
    # TODO: move map/add logic to here?
    | _query_completion(
        if .type | . == "func" or . == "var" then "_complete_scope"
        elif .type == "index" then
          if (.prefix | startswith("_")) then "_extkeys"
          else "keys"
          end
        else error("unreachable")
        end
      ) as {$type, $query, $prefix}
    | {
        prefix: $prefix,
        names: (
          if $type == "none" then
            ( $c
            | _query_index_or_key($line_query)
            | if . then [.] else [] end
            )
          else
            ( $c
            | eval($query)
            | ($prefix | _is_internal) as $prefix_is_internal
            | map(
                select(
                  strings and
                  # TODO: var type really needed? just func?
                  (_is_ident or $type == "var") and
                  ((_is_internal | not) or $prefix_is_internal or $type == "index") and
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
def _prompt:
  def _repl_level:
    (_options_stack | length | if . > 2 then ((.-2) * ">") else empty end);
  def _value_path:
    (._path? // []) | if . == [] then empty else path_to_expr end;
  def _value_preview($depth):
    if $depth == 0 and format == null and type == "array" then
      [ "["
      , if length == 0 then empty
        else
          ( (.[0] | _value_preview(1))
          , if length > 1  then ", ..." else empty end
          )
        end
      , "]"
      , if length > 1 then "[\(length)]" else empty end
      ] | join("")
    else
      ( . as $c
      | format
      | if . != null then
          ( .
          + if $c._error then "!" else "" end
          )
        else
          ($c | type)
        end
      )
    end;
  def _value:
    [ _value_path
    , _value_preview(0)
    ] | join(" ");
  def _values:
    if length == 0 then "empty"
    else
      [ (.[0] | _value)
      , if length > 1 then ", ...[\(length)]" else empty end
      ] | join("")
    end;
  [ _repl_level
  , _values
  ] | join(" ") + "> ";

# _repl_display takes a opts arg to make it possible for repl_eval to
# just call options/0 once per eval even if it was multiple outputs
def _repl_display_opts: options({depth: 1});
def _repl_display($opts): _display($opts);
def _repl_display: _display(_repl_display_opts);
def _repl_on_error:
  ( if _eval_is_compile_error then _eval_compile_error_tostring
    # was interrupted by user, just ignore
    elif _is_context_canceled_error then empty
    end
  | (_error_str | println)
  );
def _repl_on_compile_error: _repl_on_error;
def _repl_eval($expr):
  ( _repl_display_opts as $opts
  | _eval(
      $expr;
      "repl";
      _repl_display($opts);
      _repl_on_error;
      _repl_on_compile_error
    )
  );

# run read-eval-print-loop
def _repl($opts): #:: a|(Opts) => @
  def _read_expr:
    _repeat_break(
      # both _prompt and _complete want input arrays
      ( _readline(_prompt; {complete: "_complete", timeout: 0.5})
      | if trim == "" then empty
        else (., error("break"))
        end
      )
    );
  def _repl_loop:
    ( . as $c
    | try
        ( _read_expr
        | . as $expr
        | try _query_fromstring
          # TODO: nicer way to set filename for error message
          catch (. | .filename = "repl")
        | if _query_pipe_last | _query_is_func("repl") then
            ( _query_slurp_wrap(_query_func_rename("_repl_slurp"))
            | _query_tostring as $wrap_expr
            | $c
            | _repl_eval($wrap_expr)
            )
          else
            ( $c
            | .[]
            | _repl_eval($expr)
            )
          end
        )
      catch
        if . == "interrupt" then empty
        elif . == "eof" then error("break")
        elif _eval_is_compile_error then _repl_on_error
        else error
        end
    );
  ( _options_stack(. + [$opts]) as $_
  | _finally(
      _repeat_break(_repl_loop);
      _options_stack(.[:-1])
    )
  );

def _repl_slurp($opts): _repl($opts);
def _repl_slurp: _repl({});

# TODO: introspect and show doc, reflection somehow?
def help:
  ( "Type jq expression to evaluate"
  , "\\t          Auto completion"
  , "Up/Down     History"
  , "^C          Interrupt execution"
  , "... | repl  Start a new REPL"
  , "^D          Exit REPL"
  ) | println;

# just gives error, call appearing last will be renamed to _repl_slurp
def repl($_):
  if options.repl then error("repl must be last")
  else error("repl can only be used from interactive repl")
  end;
def repl: repl(null);
