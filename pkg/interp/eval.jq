include "internal";
include "query";


def _eval_error($what; $error):
  error(
    { what: $what
    , error: $error
    , column: 0
    , line: 1
    , filename: ""
    }
  );

def _eval_error_function_not_defined($name; $args):
  _eval_error(
    "compile";
    "function not defined: \($name)/\($args | length)"
  );

# if catch_query . -> try (.) catch .catch_query
# if input_query . -> .input_query | .
# if ... | <.slurp.> -> .slurp({slurp: "<slurp>", slurp_args: [arg query ast], orig: orig query ast, rewrite: rewritten query})
# else if .output_query -> . | .output_query
#
# ex ... | slurp -> <slurp>({...})
# ex no slurp: . -> try (.input_query | . | .output_query) catch .catch_query
def _eval_query_rewrite($opts):
  _query_fromtostring(
    ( . as $orig_query
    | _query_pipe_last as $last
    | ( $last
      | if _query_is_func then [_query_func_name, _query_func_args]
        else ["", []]
        end
      ) as [$last_func_name, $last_func_args]
    | $opts.slurps[$last_func_name] as $slurp
    | if $slurp then
        _query_transform_pipe_last(_query_ident)
      end
    | if $opts.catch_query then
        # _query_query to get correct precedence and a valid query
        # try (1+1) catch vs try 1 + 1 catch
        _query_try(
            ( .
            # TODO: error instead or assuming ident?
            | if (.term or .op) | not then . + _query_ident end
            | _query_query
            );
            $opts.catch_query
          )
      end
    | if $opts.input_query then
        _query_pipe($opts.input_query; .)
      end
    | if $slurp then
        _query_func(
          $slurp;
          [ # pass original, rewritten and args queries as query ast trees
            ( { slurp: _query_string($last_func_name)
              , slurp_args:
                  ( $last_func_args
                  | if . then
                      ( map(_query_toquery)
                      | _query_commas
                      | _query_array
                      )
                    else (null | _query_array)
                    end
                  )
              , orig: ($orig_query | _query_toquery)
              , rewrite: _query_toquery
              }
            | _query_object
            )
          ]
        )
      elif $opts.output_query then
        _query_pipe(.; $opts.output_query)
      end
    )
  );

# TODO: better way? what about nested eval errors?
def _eval_is_compile_error:
  _is_object and .error != null and .what != null;
def _eval_compile_error_tostring:
  [ (.filename // "expr")
  , if .line != 1 or .column != 0 then "\(.line):\(.column)"
    else empty
    end
  , " \(.error)"
  ] | join(":");

def eval($expr; $opts; on_error; on_compile_error):
  ( . as $c
  | ($opts.filename // "expr") as $filename
  | try
      _eval(
        $expr | _eval_query_rewrite($opts);
        {filename: $filename}
      )
    catch
      if _eval_is_compile_error then
        # rewrite parse error will not have filename
        ( .filename = $filename
        | {error: ., input: $c}
        | on_compile_error
        )
      else
        ( {error: ., input: $c}
        | on_error
        )
      end
  );
def eval($expr): eval($expr; {}; .error | error; .error | error);
