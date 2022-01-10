def msgpack_torepr:
  def _f:
    ( if .type | . == "fixmap" or . == "map16" or . == "map32" then
        ( .pairs
        | map({key: (.key | _f), value: (.value | _f)})
        | from_entries
        )
      elif .type | . == "fixarray" or . == "array16" or . == "array32" then .elements | map(_f)
      elif .type | . == "bin8" or . == "bin16" or . == "bin32" then .value | tostring
      else .value
      end
    );
  _f;
