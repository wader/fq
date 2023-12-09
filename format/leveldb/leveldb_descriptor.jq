# TK(2023-12-07): in the output, object-properties are
# sorted alphabetically... how can one prevent this
# and keep the original order?
def _leveldb_descriptor_torepr:
  def _f:
    if .type == "root" then
      [ .value.blocks[].records[]
      | {type: "record", value: .}
      | _f
      ]
    elif .type == "record" then
      if .value.header.record_type != "full" then
        empty
      else
        [ .value.data.tags[]
        | {type: "tag", value: {(.key): .value}}
        | _f
        ]
      end
    elif .type == "tag" then
      ( .value
      | if .comparator then .comparator |= .data else . end
      | if .new_file then
          ( .new_file.largest_internal_key |= .data
          | .new_file.smallest_internal_key |= .data
          )
        else .
        end
      | if .compact_pointer then
          .compact_pointer.internal_key |= .data
        else .
        end
      )
    end;
  ( {type: "root", value: .}
  | _f
  );