def _bplist_torepr:
  def _f:
    ( if .type == "singleton" then .value | tovalue
      elif .type == "int" then .value | tovalue
      elif .type == "real" then .value | tovalue
      elif .type == "date" then .value | tovalue
      elif .type == "data" then .value | tovalue
      elif .type == "ascii_string" then .value | tovalue
      elif .type == "unicode_string" then .value | tovalue
      elif .type == "uid" then {"cfuid": .value | tovalue}
      elif .type == "array" then
        ( .entries
        | map(_f)
        )
      elif .type == "set" then
        ( .entries
        | map(_f)
        )
      elif .type == "dict" then
        ( .entries
        | map({key: (.key | _f), value: (.value | _f)})
        | from_entries
        )
      else  error("unknown type: \(.type)")
      end
    );
  ( .objects
  | _f
  );

