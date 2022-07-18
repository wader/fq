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

def _cbor__help:
  { links: [
      {url: "https://en.wikipedia.org/wiki/CBOR"},
      {url: "https://www.rfc-editor.org/rfc/rfc8949.html"}
    ]
  };
