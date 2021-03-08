#!/usr/bin/env fq -rnf

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


