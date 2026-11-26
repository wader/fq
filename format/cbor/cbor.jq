def _cbor_torepr:
  if .major_type == "map" then
    ( .pairs
    | map({key: (.key | _cbor_torepr), value: (.value | _cbor_torepr)})
    | from_entries
    )
  elif .major_type == "array" then .elements | map(_cbor_torepr)
  elif .major_type == "bytes" then .value | tostring
  elif .major_type == "semantic" then
    ( (.tag | toactual) as $tag
    | .value
    | if ($tag == 2 or $tag == 3) and .major_type == "bytes" then
        ( .value
        | tobytes
        | explode
        | reduce .[] as $b (0; . * 256 + $b)
        | if $tag == 3 then -1 - . end
        )
      else _cbor_torepr
      end
    )
  else .value | tovalue
  end;
