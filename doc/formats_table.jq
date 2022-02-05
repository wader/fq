#!/usr/bin/env fq -rnf

def code: "`\(.)`";
def nbsp: gsub(" "; "&nbsp;");

def format_table:
  ( ($doc_formats | split(" ")) as $doc_formats
  | [ {
        name: "Name",
        desc: "Description",
        uses: "Dependencies"
      },
      {
        name: "-",
        desc: "-",
        uses: "-"
      },
      ( formats
      | to_entries[]
      | {
          name:
            ( ( .key as $format
              | if ($doc_formats | indices($format)) != [] then "[\($format | code)](#\($format))"
                else $format | code
                end
              )
            + " "
            ),
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
      , (.[0] | . as $rc | $rc.string | rpad(" "; $rc.maxwidth))
      , (.[1] | . as $rc | $rc.string | rpad(" "; $rc.maxwidth))
      , .[2].string
      , ""
      ] | join("|")
    )
  );

format_table
