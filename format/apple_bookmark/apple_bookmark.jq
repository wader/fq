def _apple_bookmark_torepr:
  def _f:
    ( if .type == "String" then .data | tovalue
      elif .type == "Data" then .data | tovalue
      elif .type == "Byte" then .data | tovalue
      elif .type == "Short" then .data | tovalue
      elif .type == "Int" then .data | tovalue
      elif .type == "Long" then .data | tovalue
      elif .type == "Float" then .data | tovalue
      elif .type == "Double" then .data | tovalue
      elif .type == "Date" then .data | tovalue
      elif .type == "BooleanFalse" then false
      elif .type == "BooleanTrue" then true
      elif .type == "Array" then 
        ( .data
        | map(.record | _f)
        )
      elif .type == "Dictionary" then
        ( .data
        | map({key: (.key | _f), value: (.value | _f)})
        | from_entries
        )
      elif .type == "UUID" then .data | tovalue
      elif .type == "URL" then .data | tovalue
      elif .type == "RelativeURL" then
		.data | map(.record.data)
      end
     );
  ( .bookmark_entries
  | map({key: (.key_string?.record.data // .key|tostring), value: (.record | _f)})
  | from_entries
  );


