# note this is a "dynamic" include, output string will be used as source

# generates a _format_func function that can be used to implement format overloaded
# functions like torepr, _format_func("msgpack", "torepr") calls _msgpack_torepr

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
