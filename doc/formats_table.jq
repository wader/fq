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
        uses: "<sub>\((((.value.dependencies | flatten | map(code)) | join(" "))? // ""))</sub>"
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
          uses: "<sub>\(((.value | map(code)) | join(" ")))</sub>"
      }
    )
]
| table(
    [.name, .desc, .uses];
    [ ""
    , (.[0] | . as $rc | $rc.string | rpad(" ";$rc.maxwidth))
    , (.[1] | . as $rc | $rc.string | rpad(" ";$rc.maxwidth))
    , .[2].string
    , ""
    ] | join("|")
  )


