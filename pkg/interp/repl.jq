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

def _prompt:
  def _type_name_error:
    ( . as $c
    | try
        ( _display_name
        , if ._error then "!" else empty end
        )
      catch ($c | type)
    );
  def _path_prefix:
    (._path? // []) | if . == [] then "" else path_to_expr end;
  def _preview:
    if format != null or type != "array" then
      _type_name_error
    else
      ( "["
      , if length > 0 then (.[0] | _type_name_error) else empty end
      , if length > 1 then ", ..." else empty end
      , "]"
      , if length > 1 then "[\(length)]" else empty end
      )
    end;
  ( [ (_options_stack | length | if . > 2 then ((.-2) * ">") + " " else empty end)
    , if length == 0 then
        "empty"
      else
        ( .[0]
        | _path_prefix
        , _preview
        )
      end
    , if length > 1 then ", [\(length)]" else empty end
    , "> "
    ]
  ) | join("");

def _repl_display: _display({depth: 1});
def _repl_on_error:
  ( if _eval_is_compile_error then _eval_compile_error_tostring end
  | (_error_str | println)
  );
def _repl_on_compile_error: _repl_on_error;
def _repl_eval($expr): _eval($expr; "repl"; _repl_display; _repl_on_error; _repl_on_compile_error);

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

# just gives error, call appearing last will be renamed to _repl_slurp
def repl($_):
  if options.repl then error("repl must be last")
  else error("repl can only be used from interactive repl")
  end;
def repl: repl(null);
