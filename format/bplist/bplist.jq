def _bplist_torepr:
  def _f:
    ( if .type == "singleton" then .value
      elif .type == "int" then .value
      elif .type == "real" then .value
      elif .type == "date" then .value
      elif .type == "data" then .value.data
      elif .type == "ascii_string" then .value
      elif .type == "unicode_string" then .value
      elif .type == "uid" then .value
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

