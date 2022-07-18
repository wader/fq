# note this is a "dynamic" include, output string will be used as source

[ "def _format_func($format; $func):"
, "  ( [$format, $func] as $ff"
, "  | if false then error(\"unreachable\")"
, ( _registry.formats[] as $f
  | $f.functions[]?
  | "    elif $ff == \([$f.name, .] | tojson) then _\($f.name)_\(.)"
  )
  , "    else error(\"\\($format) has no \\($func)\")"
  , "    end"
  , "  );"
] | join("\n")
