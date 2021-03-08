#!/usr/bin/env fq -rnf

# {columns: [{name: "name", "title": "Name"}, ...], rows: [{name: "Abc", ...}, ...]}
def table:
    def _rpad($s;$w): . + ($s * ($w+1-length))[1:];
    def _column_widths:
        [. as {columns: $cs, rows: $rs}
            | $cs[]
            | . as $c
            | [$c.title, ($rs[] | .[$c.name])]
            | {
            ($c.name): (. | map(length) | max)
            }
        ] | add;
    def _row: ([""] + .  + [""]) | join("|");
    _column_widths as $cw
    | . as {columns: $cs, rows: $rs}
    | ( ($cs | map(. as $c | .title | _rpad(" ";$cw[$c.name])) | _row)
      , ([range($cs | length) | "-"]  | _row)
      , ($rs[]
        | . as $r
        | [ ($cs[]
            | . as $c
            | ($r[$c.name] | _rpad(" ";$cw[$c.name]))
            )
          ]
        | _row
        )
      );


def code: "`" + . + "`";
def nbsp: gsub(" ";"&nbsp;");

{
    columns: [
        {name: "name", "title": "Name"},
        {name: "desc", "title": "Description"},
        {name: "uses", "title": "Uses"}
    ],
    rows: [
        ( formats
          | to_entries[]
          | {
            name: ((.key | code) + " "),
            desc: ((.value.description | nbsp) + " "),
            uses: (((.value.dependencies | flatten | map(code)) | join(", "))? // "")
          }
        ),
        (
        [ formats
          | to_entries[]
          | . as $e
          | select(.value.groups)
          | .value.groups[] | {key: ., value: $e.key}
        ]
        | reduce .[] as $e ({}; .[$e.key] += [$e.value])
        | to_entries[]
        | {
            name: ((.key | code) + " "),
            desc: "Group",
            uses: ((.value | map(code)) | join(", "))
        }
        )
    ]
} | table
