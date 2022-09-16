def _msgpack_torepr:
  if .type | . == "fixmap" or . == "map16" or . == "map32" then
    ( .pairs
    | map({key: (.key | _msgpack_torepr), value: (.value | _msgpack_torepr)})
    | from_entries
    )
  elif .type | . == "fixarray" or . == "array16" or . == "array32" then .elements | map(_msgpack_torepr)
  elif .type | . == "bin8" or . == "bin16" or . == "bin32" then .value | tostring
  else .value | tovalue
  end;

