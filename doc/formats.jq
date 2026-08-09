#!/usr/bin/env fq -rnf

def formats_list:
  [ formats[] as $f
  | ({} | _help_format_enrich("fq"; $f; false)) as $fhelp
  | $f.name
  ] | join(",\n");
