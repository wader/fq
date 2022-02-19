#!/usr/bin/env fq -rnf

[ (formats | keys[]) as $format
| if ($doc_formats | indices($format)) != [] then "[\($format)](doc/formats.md#\($format))"
  else $format
  end
] | join(",\n")
