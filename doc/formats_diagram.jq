#!/usr/bin/env fq -rnf

def color:
  to_md5 | [.[range(3)]] | map(band(.; 0x7f)+60 | to_radix(16) | "0"[length:]+.) | join("");

def _formats_dot:
  def _record($title; $fields):
    [  "<"
    , "<table bgcolor=\"paleturquoise\" border=\"0\" cellspacing=\"0\">"
    , "<tr><td port=\"\($title)\"><font point-size=\"20\">\($title)</font></td></tr>"
    , [$fields | flatten | map("<tr><td align=\"left\" bgcolor=\"lightgrey\" port=\"\(.)\">\(.)</td></tr>")]
    , "</table>"
    , ">"
    ] | flatten | join("");
  ( "# ... | dot -Tsvg -o formats.svg"
  , "digraph formats {"
  , "  nodesep=0.2"
  , "  ranksep=1"
  , "  rankdir=TB"
  , "  node [penwidth=2 shape=\"none\" style=\"\"]"
  , "  edge [penwidth=2]"
  , ( .[]
    | . as $f
    | .dependencies
    | flatten?
    | .[]
    | "  \"\($f.name)\":\(.):e -> \(.):n [color=\"#\($f.name | color)\"]"
    )
  , ( .[]
    | .name as $name
    | .groups[]?
    | "  \(.) -> \"\($name)\":\($name):n [color=\"#\(. | color)\"]"
    )
  , ( to_entries[]
    | "  \(.key) [color=\"paleturquoise\", label=\(_record(.key; (.value.dependencies // [])))]"
    )
  , ( [.[].groups[]?]
    | unique[]
    | "  \(.) [shape=\"record\",style=\"rounded,filled\",fontsize=\"25\"color=\"palegreen\"]"
    )
  , "}"
  );

formats | _formats_dot
