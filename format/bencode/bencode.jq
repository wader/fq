def _bencode_torepr:
  if .type == "string" then .value | tovalue
  elif .type == "integer" then .value | tovalue
  elif .type == "list" then .values | map(_bencode_torepr)
  elif .type == "dictionary" then
    ( .pairs
    | map({key: (.key | _bencode_torepr), value: (.value | _bencode_torepr)})
    | from_entries
    )
  else error("unknown type \(.type)")
  end;
