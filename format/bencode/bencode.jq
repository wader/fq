def bencode_torepr:
  def _f:
    if .type == "string" then .value
    elif .type == "integer" then .value
    elif .type == "list" then .values | map(_f)
    elif .type == "dictionary" then
      ( .pairs
      | map({key: (.key | _f), value: (.value | _f)})
      | from_entries
      )
    else error("unknown type \(.type)")
    end;
  _f;
