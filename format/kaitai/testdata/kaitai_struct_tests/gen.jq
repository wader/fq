"$ fq -d kaitai -o source=@formats/\($n | gsub(".kst$";".ksy")) \"d,tovalue({skip_gaps: true})\" src/\(.data)",
( reduce
    ( .asserts[]?
    # "some.as<type>.path" -> ".some.path"
    | { path: ("." + .actual | gsub(".as<.*>"; "") | expr_to_path)
      , value:
          ( .expected
          | if type == "string" then
              from_jq? // .
            end
          )
      }
    ) as {$path,$value} (
      {};
      setpath($path; $value)
    )
)
