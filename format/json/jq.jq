# to jq-flavoured json
def _to_jq($opts):
  def _is_ident: test("^[a-zA-Z_][a-zA-Z_0-9]*$");
  def _key: if _is_ident | not then tojson end;
  def _f($opts; $indent):
    def _r($prefix):
      ( type as $t
      | if $t == "null" then tojson
        elif $t == "string" then tojson
        elif $t == "number" then tojson
        elif $t == "boolean" then tojson
        elif $t == "array" then
          if length == 0 then "[]"
          else
            [ "[", $opts.compound_newline
            , ( [ .[]
                | $prefix, $indent
                , _r($prefix+$indent), $opts.array_sep
                ]
              | .[0:-1]
              )
            , $opts.compound_newline
            , $prefix, "]"
            ]
          end
        elif $t == "object" then
          if length == 0 then "{}"
          else
            [ "{", $opts.compound_newline
            , ( [ to_entries[]
                | $prefix, $indent
                , (.key | _key), $opts.key_sep
                , (.value | _r($prefix+$indent)), $opts.object_sep
                ]
              | .[0:-1]
              )
            , $opts.compound_newline
            , $prefix, "}"
            ]
          end
        else error("unknown type \($t)")
        end
      );
    _r("");
  ( _f($opts; $opts.indent * " ")
  | if _is_array then flatten | join("") end
  );
def to_jq($opts):
  _to_jq(
    ( { indent: 0
      , key_sep: ":"
      , object_sep: ","
      , array_sep: ","
      , compound_newline: "",
      } + $opts
    | if .indent > 0  then
        ( .key_sep = ": "
        | .object_sep = ",\n"
        | .array_sep = ",\n"
        | .compound_newline = "\n"
        )
      end
    )
  );
def to_jq: to_jq(null);

# from jq-flavoured json
def from_jq:
  def _f:
    ( . as $v
    | .term.type
    | if . == "TermTypeNull" then null
      elif . == "TermTypeTrue" then true
      elif . == "TermTypeFalse" then false
      elif . == "TermTypeString" then
        if $v.term.str.queries then error("string interpolation")
        else $v.term.str.str
        end
      elif . == "TermTypeNumber" then $v.term.number | tonumber
      elif . == "TermTypeObject" then
        ( $v.term.object.key_vals // []
        | map(
            { key: (.key // .key_string.str)
            , value: (.val | _f)
            }
          )
        | from_entries
        )
      elif . == "TermTypeArray" then
        ( def _a: if .op then .left, .right | _a end;
          [$v.term.array.query // empty | _a | _f]
        )
      else error("unsupported term \($v.term.type)")
      end
    );
  try
    (_query_fromstring | _f)
  catch
    error("from_jq only supports constant literals: \(.)");
