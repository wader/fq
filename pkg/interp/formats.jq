# note this is a "dynamic" include, outputted string will be used as source

def _formats_source:
  ( [ ( _registry.groups
      | to_entries[]
      # TODO: nicer way to skip "all" which also would override builtin all/*
      | select(.key != "all")
      | "def \(.key)($opts): decode(\(.key | tojson); $opts);"
      , "def \(.key): decode(\(.key | tojson); {});"
      )
    , ( _registry.formats[]
      | select(.files)
      | .files[]
      )
    , ( "def torepr:"
      , "  ( format as $f"
      , "  | if $f == null then error(\"value is not a format root\") end"
      , "  | if false then error(\"unreachable\")"
      , ( _registry.formats[]
        | select(.to_repr != "")
        | "    elif $f == \(.name | tojson) then \(.to_repr)"
        )
      , "    else error(\"format has no torepr\")"
      , "    end"
      , "  );"
      )
    ]
  | join("\n")
  );

_formats_source
