def bson_torepr:
  def _torepr:
   ( if .type == null or .type == "array" then
      ( .value.elements
      | map(_torepr)
      )
     elif .type == "document" then
      ( .value.elements
      | map({key: .name, value: _torepr})
      | from_entries
      )
     elif .type == "boolean" then .value != 0
     else .value
     end
   );
  ( {type: "document", value: .}
  | _torepr
  );
