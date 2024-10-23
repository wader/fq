# null
def _query_null:
  {term: {type: "TermTypeNull"}};

# . -> (.)
def _query_query:
  { term:
      { type: "TermTypeQuery"
      , query: .
      }
  };

# string
def _query_string($str):
  { term:
      { type: "TermTypeString"
      , str: {str: $str}
      }
  };

# .
def _query_ident:
  {term: {type: "TermTypeIdentity"}};
def _query_is_ident:
  .term.type == "TermTypeIdentity";

# a($args...) -> b($args...)
def _query_func_rename(name):
  .term.func.name = name;
# $name($args)
def _query_func($name; $args):
  { term:
      { type: "TermTypeFunc"
      , func:
          { args: $args
          , name: $name
          }
      }
  };
def _query_func($name):
  _query_func($name; null);

def _query_func_name:
  .term.func.name;
def _query_func_args:
  .term.func.args;
def _query_is_func:
  .term.type == "TermTypeFunc";
def _query_is_func($name):
  _query_is_func and _query_func_name == $name;

def _query_is_string:
  .term.type == "TermTypeString";
def _query_string_str:
  .term.str.str;

def _query_empty:
  _query_func("empty");

# l | r
def _query_pipe(l; r):
  { op: "|"
  , left: l
  , right: r
  };

# . -> [.]
def _query_array:
  ( . as $q
  | { term:
        { type: "TermTypeArray"
        , array: {}
        }
    }
  | if $q then .term.array.query = $q end
  );

# {} -> {}
def _query_object:
  { term:
      { object:
          { key_vals:
            ( to_entries
            | map(
                { key: .key
                , val: .value
                }
              )
            )
          }
      , type: "TermTypeObject"
      }
  };

# l,r
def _query_comma(l; r):
  { left: l
  , op: ","
  , right: r
  };

# [1,2,3] -> 1,2,3
# output each query in array
def _query_commas:
  if length == 0 then _query_empty
  else
    reduce .[1:][] as $q (
      .[0];
      _query_comma(.; $q)
    )
  end;

# . -> .[]
def _query_iter:
  .term.suffix_list = [{iter: true}];

# try b catch c
def _query_try(b; c):
  { term:
      { type: "TermTypeTry"
      , try:
          { body: b
          , catch: c
          }
      }
  };
def _query_try(b):
  _query_try(b; null);

# last query in pipeline
def _query_pipe_last:
  if .term.suffix_list then
    ( .term.suffix_list[-1]
    | if .bind.body then
        ( .bind.body
        | _query_pipe_last
        )
      end
    )
  elif .op == "|" then
    ( .right
    | _query_pipe_last
    )
  end;

def _query_transform_pipe_last(f):
  def _f:
    if .term.suffix_list then
      .term.suffix_list[-1] |=
        if .bind.body then
          .bind.body |= _f
        else f
        end
    elif .op == "|" then
      .right |= _f
    else f
    end;
  _f;

# last term, the right most
def _query_last:
  if .term.suffix_list then
    ( .term.suffix_list[-1]
    | if .bind.body then
        ( .bind.body
        | _query_last
        )
      end
    )
  elif .op then
    ( .right
    | _query_last
    )
  end;

# TODO: rename? what to call when need to transform suffix_list
def _query_transform_last(f):
  def _f:
    if .term.suffix_list then
      ( .
      | if .term.suffix_list[-1].bind.body then
          .term.suffix_list[-1].bind.body |= _f
        else f
        end
      )
    elif .op then
      .right |= _f
    else f
    end;
  _f;

def _query_completion_type:
  ( . as $q
  | _query_last
  | if .index.name then
      { query:
          ( $q
          | _query_transform_last(
              del(.term.suffix_list[-1])
            )
          )
      , type: "index"
      , prefix: .index.name
      }
    elif .term.index.name then
      { query:
          ( $q
          | _query_transform_last(
              _query_ident
            )
          )
      , type: "index"
      , prefix: .term.index.name
      }
    elif .term.func then
      { query:
          ( $q
          | _query_transform_last(
              _query_ident
            )
          )
      , type:
          ( .term.func.name
          | if startswith("$") then "var"
            else "func"
            end
          )
      , prefix: .term.func.name
      }
    else
      null
    end
  );

# TODO: simplify
def _query_completion(f):
  ( . as $expr
  # HACK: if ends with . or $, add a dummy prefix to make the query
  # valid and then trim it later
  | ( if (.[-1] | . == "." or . == "$") then "a"
      else ""
      end
    ) as $probesuffix
  | ($expr + $probesuffix)
  | try
      ( try _query_fromstring
        catch .error
      # move directives to new root query
      | . as {$meta, $imports}
      | del(.meta)
      | del(.imports)
      | _query_completion_type
      | . as $c
      | if . then
          ( .query |=
              ( _query_func("map"; [
                  _query_pipe(.; _query_try(_query_func($c | f)))
                ])
              | _query_pipe(.; _query_func("add")
                )
              | .meta = $meta
              | .imports = $imports
              | _query_tostring
              )
          | .prefix |= rtrimstr($probesuffix)

          )
        else
          {type: "none"}
        end
      )
    catch {type: "error", name: "", error: .}
  );

# query ast to ast of quey itself, used by query rewrite/slurp
def _query_toquery:
  ( tojson
  | _query_fromstring
  );

# query rewrite helper, takes care of from/to and directives
def _query_fromtostring(f):
  ( _query_fromstring
  # save and move directives to possible new root query
  | . as {$meta, $imports}
  | del(.meta)
  | del(.imports)
  | f
  | .meta = $meta
  | .imports = $imports
  | _query_tostring
  );
