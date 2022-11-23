def _apple_bookmarkdata_torepr:
  def _f:
    ( if .type == 0x101 then .value | tovalue
      elif .type == "String" then .data | tovalue
      elif .type == "Data" then .data | tovalue
      elif .type == "Byte" then .data | tovalue
      elif .type == "Short" then .data | tovalue
      elif .type == "Int" then .data | tovalue
      elif .type == "Long" then .data | tovalue
      elif .type == "Float" then .data | tovalue
      elif .type == "Double" then .data | tovalue
      elif .type == "BooleanFalse" then false
      elif .type == "BooleanTrue" then false
      elif .type == "Array" then 
        ( .data
        | map(_f)
        )
      elif .type == "Dictionary" then
        ( .data
        | map({key: (.key | _f), value: (.value | _f)})
        | from_entries
        )
      elif .type == "UUID" then .data | tovalue
      elif .type == "URL" then .data | tovalue
      elif .type == "RelativeURL" then .data | tovalue
      end
     );
  ( .bookmark_entries
  | map({(.key): (.record | _f)})
  );

