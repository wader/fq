def _markdown__todisplay: tovalue;

def _word_break($width):
  def _f($a; $acc; $l):
    ( $a[0] as $w
    | ($w // "" | length+1) as $wl
    | if $w == null then $acc
      elif ($l + $wl) >= $width then
        ( $acc
        , _f($a[1:]; [$w]; $wl)
        )
      else _f($a[1:]; $acc+[$w]; $l+$wl)
      end
    );
  ( [_f([splits("\\s{1,}")]; []; 0)]
  | map(join(" "))
  );


# for a document output {heading: <heading>, children: [<children until next heading>]}
# heading can be null for a document with children before a heading
def _markdown_split_headings:
  foreach
    ( .children[]
    , {type:"heading"} # dummy heading to flush
    ) as $c
  (
    {heading: null, children: null, extract: null};
    if $c.type == "heading" then
      ( .extract = {heading,children}
      | .heading = $c
      | .children = null
      )
    else
      ( .children += [$c]
      | .extract = null
      )
    end;
    .extract | select(.heading or .children)
  );

def _markdown_children_to_text($width):
  def lb: if $width > 0 then _word_break($width) | join("\n") end;
  def _f:
    if type == "string" then gsub("\n"; " ")
    elif .type == "document" then .children[] | _f
    elif .type == "heading" then .children[] | _f
    elif .type == "paragraph" then
      ( [.children[] | _f]
      | join("")
      | lb
      )
    elif .type == "link" then
      ( ( [ .children[]
          | _f
          ]
        | join("")
        ) as $text
      | if $text == .destination then $text
        else "\($text) (\(.destination))"
        end
      )
    elif .type == "code_block" then .literal | rtrimstr("\n") | split("\n") | "  " + join("\n  ")
    elif .type == "code" then .literal
    elif .type == "list" then ([.children[] | _f] | join("\n")) # TODO: delim
    elif .type == "list_item" then "\(.bullet_char) \(.children[] | _f)"
    elif .type == "html_span" then .literal | gsub("<br>"; "\n") # TODO: more?
    else empty
    end;
  [_f] | join("\n");

def _markdown_to_text($width; $header_depth):
  [ _markdown_split_headings
  | if .heading then
      ( (.heading | _markdown_children_to_text($width)) as $h
      | $h
      , ("=" * ($h | length))
      )
    else empty
    end
  , ( .children
    | if length == 0 then ""
      else
        ( .[]
        | _markdown_children_to_text($width)
        | select(. != "")
        | .
        , ""
        )
      end
    )
  ][:-1] | join("\n");
def _markdown_to_text:
  _markdown_to_text(-1; 0);
