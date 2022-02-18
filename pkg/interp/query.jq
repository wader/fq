# a() -> b()
def _query_func_rename(name):
  .term.func.name = name;

# . | r
def _query_pipe(r):
  { op: "|",
    left: .,
    right: r
  };

# . -> .[]
def _query_iter:
  .term.suffix_list = [{iter: true}];

def _query_ident:
  {term: {type: "TermTypeIdentity"}};

def _query_try(f):
  { term: {
      try: {
        body: f,
      },
      type: "TermTypeTry"
    }
  };

def _query_func($name; $args):
  { term: {
      func: {
        args: $args,
        name: $name
      },
      type: "TermTypeFunc"
    }
  };

def _query_func($name):
  _query_func($name; null);

def _query_is_func(name):
  .term.func.name == name;

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
          ($q | _query_transform_last(
            del(.term.suffix_list[-1])
          )),
        type: "index",
        prefix: .index.name
      }
    elif .term.index.name then
      { query:
          ($q | _query_transform_last(
            _query_ident
          )),
        type: "index",
        prefix: .term.index.name
      }
    elif .term.func then
      { query:
          ($q | _query_transform_last(
            _query_ident
          )),
        type:
          ( .term.func.name
          | if startswith("$") then "var"
            else "func"
            end
          ),
        prefix: .term.func.name
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
                  _query_pipe(
                    _query_try(
                      _query_func($c | f)
                    )
                  )
                ])
              | _query_pipe(
                  _query_func("add")
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

# <filter...> | <slurp_func> ->
# map(<filter...> | .) | (<slurp_func> | f)
def _query_slurp_wrap(f):
  # save and move directives to new root query
  ( . as {$meta, $imports}
  | del(.meta)
  | del(.imports)
  | _query_pipe_last as $lq
  | _query_transform_pipe_last(_query_ident) as $pipe
  | _query_func("map"; [$pipe])
  | _query_pipe($lq | f)
  | .meta = $meta
  | .imports = $imports
  );

# filter -> .[] | filter
def _query_iter_wrap:
  ( . as $q
  | _query_ident
  | _query_iter
  | _query_pipe($q)
  );
