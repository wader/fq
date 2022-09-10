def _cbor_torepr:
  if .major_type == "map" then
    ( .pairs
    | map({key: (.key | _cbor_torepr), value: (.value | _cbor_torepr)})
    | from_entries
    )
  elif .major_type == "array" then .elements | map(_cbor_torepr)
  elif .major_type == "bytes" then .value | tostring
  else .value | tovalue
  end;
