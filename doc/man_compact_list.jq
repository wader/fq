# post-process asciidoctor manpage roff output:
# use "-" as bullet marker with one space to text and
# remove blank lines between items.

# replace \(bu bullet marker with \- (ascii "-") and one space
( split("\\h'-04'\\(bu\\h'+03'\\c")
| join("\\h'-04'\\-\\h'+01'\\c")
| split(".IP \\(bu 2.3")
| join(".IP \\- 2.3")

# track list nesting by counting .RS/.RE blocks
| split("\n")
| . as $l
| foreach range(length) as $i (
    {d: 0, o: []};
    ( .d +=
        if $l[$i] == ".RS 4" or $l[$i] == ".if n .RS 4" then 1
        elif $l[$i] == ".RE" or $l[$i] == ".if n .RE" then -1
        else 0
        end
    # .sp directly before a bullet list start?
    | if $l[$i] == ".sp" and
        ($l[$i+1] // "") == ".RS 4" and
        (($l[$i+2] // "") | startswith(".ie n ")) then
        if ($l[$i-1] // "") == ".RE" then .o = ".br"      # between items
        elif .d > 0 then .o = null    # start of nested list, remove blank line
        else .o = $l[$i] # top-level list, keep spacing
        end
      else .o = $l[$i]
      end
    );
    .o | values
  )
)
