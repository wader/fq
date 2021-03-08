#!/usr/bin/env fq -rnf

# [{a: 123, ...}, ...]
# colmap maps something into [col, ...]
# render maps [{string: "coltext", maxwidth: 12}, ..] into a row string
def table(colmap;render):
    def _column_widths:
        [ . as $rs
          | range($rs[0] | length) as $i
          | [$rs[] | colmap | (.[$i] | length)]
          | max
        ];
    if (. | length) == 0 then ""
    else
      _column_widths as $cw
      | . as $rs
      | ( ($rs[]
          | . as $r
          | [ range($r | length) as $i
              | ($r | colmap | {string: .[$i], maxwidth: $cw[$i]})
            ]
          | render
          )
        )
      end;


def rpad($s;$w): . + ($s * ($w+1-length))[1:];
def code: "`" + . + "`";
def nbsp: gsub(" ";"&nbsp;");

[
    {
        name: "Name",
        desc: "Description",
        uses: "Uses"
    },
    {
        name: "-",
        desc: "-",
        uses: "-"
    },
    ( formats
      | to_entries[]
      | {
        name: ((.key | code) + " "),
        desc: ((.value.description | nbsp) + " "),
        uses: (((.value.dependencies | flatten | map(code)) | join(", "))? // "")
      }
    ),
    ( [ formats
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
| table(
    [.name, .desc, .uses];
    ([""] + map(. as $rc | $rc.string | rpad(" ";$rc.maxwidth)) + [""]) | join("|")
  )
