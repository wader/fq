def _cbor_torepr:
  def _f:
    ( if .major_type == "map" then
        ( .pairs
        | map({key: (.key | _f), value: (.value | _f)})
        | from_entries
        )
      elif .major_type == "array" then .elements | map(_f)
      elif .major_type == "bytes" then .value | tostring
      else .value | tovalue
      end
    );
  _f;
