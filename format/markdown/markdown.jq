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

def _markdown_to_text($width; $header_depth):
  def lb: if $width > 0 then _word_break($width) | join("\n") end;
  def _f:
    if type == "string" then gsub("\n"; " ")
    elif .type == "document" then .children[] | _f
    elif .type == "heading" then
      ( (.children[] | _f) as $title
      | "\($title)\n\("=" * ($title | length))"
      )
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
  [_f] | join("\n\n");
def _markdown_to_text:
  _markdown_to_text(-1; 0);
