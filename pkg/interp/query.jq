# []
def _query_array:
  {
    term: {
      type: "TermTypeArray",
      array: {
        query: .
      }
    }
  };

# a() -> b()
def _query_func_rename(name): .term.func.name = name;

# . | r
def _query_pipe(r):
  { op: "|",
    left: .,
    right: r
  };

def _query_ident: {term: {type: "TermTypeIdentity"}};

# .[]
def _query_iter:
  { "term": {
      "suffix_list": [{
        "iter": true
      }],
      "type": "TermTypeIdentity"
    }
  };

# last query in pipeline
def _query_pipe_last:
  if .term then
    ( . as $t
    | .term
    | if .suffix_list then
        ( .suffix_list[-1]
        | if .bind.body then (.bind.body | _query_pipe_last)
          else .
          end
        )
      else $t
      end
    )
  elif .op == "|" then (.right | _query_pipe_last)
  else .
  end;

def _query_is_func(name): .term.func.name == name;

def _query_replace_last(f):
  # TODO: hack TCO bug
  def _f:
    if .term.suffix_list then
      .term.suffix_list[-1] |=
        if .bind.body then (.bind.body |= _f)
        else f
        end
    elif .term then f
    elif .op == "|" then (.right |= _f)
    else f
    end;
  _f;

def _query_find(f):
  ( if f then . else empty end
  , if .op == "|" or .op == ","  then
      ( (.left | _query_find(f))
      , (.right | _query_find(f))
      )
    elif .term.suffix_list then
      ( .term.suffix_list
      | map(.bind.body | _query_find(f))
      )
    else empty
    end
  );

# <filter...> | <slurp_func> -> [.[] | <filter...> | .] | (<slurp_func> | f)
def _query_slurp_wrap(f):
  ( _query_pipe_last as $lq
  | _query_replace_last(_query_ident) as $pipe
  | _query_iter
  | _query_pipe($pipe)
  | _query_array
  | _query_pipe($lq | f)
  );
