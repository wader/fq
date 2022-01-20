#!/usr/bin/env fq -rnf

def _formats_dot:
  def _record($title; $fields):
    [  "<"
    , "<table bgcolor=\"paleturquoise\" border=\"0\" cellspacing=\"0\">"
    , "<tr><td port=\"\($title)\">\($title)</td></tr>"
    , [$fields | flatten | map("<tr><td align=\"left\" bgcolor=\"lightgrey\" port=\"\(.)\">\(.)</td></tr>")]
    , "</table>"
    , ">"
    ] | flatten | join("");
  ( "# ... | dot -Tsvg -o formats.svg"
  , "digraph formats {"
  , "  rankdir=TB"
  , "  node [shape=\"none\" style=\"\"]"
  , ( .[]
    | . as $f
    | .dependencies
    | flatten?
    | .[]
    | "  \"\($f.name)\":\(.):e -> \(.):n"
    )
  , ( .[]
    | .name as $name
    | .groups[]?
    | "  \(.) -> \"\($name)\":\($name):n"
    )
  , ( to_entries[]
    | "  \(.key) [color=\"paleturquoise\", label=\(_record(.key; (.value.dependencies // [])))]"
    )
  , ( [.[].groups[]?]
    | unique[]
    | "  \(.) [shape=\"record\",style=\"rounded,filled\",color=\"palegreen\"]"
    )
  , "}"
  );

formats | _formats_dot
