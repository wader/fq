#!/usr/bin/env fq -f

def recurse_depth(f; cond):
  def _r($depth):
    ( ( .
      | select(cond)
      | {depth: $depth, value: .}
      )
    , ( [f]
      | to_entries[] as $e
      | $e.value
      | _r($depth + if cond then 1 else 0 end)
      )
    );
  _r(0);

[ ( recurse_depth(
      .[]?;
      format
    )
  | . + {
      norm_path: (.value._path | map(if type == "number" then "index" end)),
    }
  )
]
| streaks_by(.norm_path)
| map(.[0] + {count: length})
| .[]
| [ if .depth > 0 then "  "*.depth else empty end
  , ((.value | format) + if .count > 1 then "*\(.count) " else " " end)
  , (.value._path | path_to_expr)
  ]
| join("")
| println
