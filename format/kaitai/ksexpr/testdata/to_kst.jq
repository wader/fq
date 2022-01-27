def to_kst:
  { id: "test"
  , data: ""
  , asserts:
      ( to_entries
      | sort_by(.key)
      | map(
          ( . as {$key, $value}
          | { actual: $key
            , expected:
                ( $value
                # TODO: hack to workaround kaitai ruby byte array output
                | if ($key | startswith("array")) and ($value | type == "string") then
                    split(" ") | map(tonumber)
                  else .
                  end
                | tojson
                )
            }
          )
        )
      )
  };
