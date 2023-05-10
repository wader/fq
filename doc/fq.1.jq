( ".TH fq 1"
, ".SH NAME"
, _help("fq"; "fq_banner")

, ".SH SYNOPSIS"
, _help("fq"; "fq_usage")
, ""
, _help(""; "fq_example_usage")

, ".SH DESCRIPTION"

, _help(""; "fq_summary")
, ""
, "fq is inspired by the well known jq tool and language and allows you to work with binary formats the same way you would using jq. In addition it can present data like a hex viewer, transform, slice and concatenate binary data. It also supports nested formats and has an interactive REPL with auto-completion."

, ".SH OPTIONS"
, ( _opt_cli_args
  | to_entries[] as $a
  | $a.value
  | ".TP"
  , ( "\(.long)\(.short | if . then ",\(.)" else "" end)"
    + ( .string // .array // .object // .pairs
      | if . then " \(.)"
        else ""
        end
      )
    )
  , .description
  , ( select($a.key == "option")
    | ( _opt_options
      | to_entries[] as $o
      | select($o.value.help != null)
      | ".TP"
      , "  -o \($o.key)=<\($o.value.type)>"
      , "    \($o.value.help)"
      )
    )
  )

, ".SH ENVIRONMENT"
, ( _opt_environment
  | to_entries[] as $e
  | ".TP"
  , $e.key
  , $e.value
  )

, ".SH SUPPORTED FORMATS"
, "By default fq will try to probe input format. If this does not work"
, "a format can by specified by using -d <NAME>."
, "To see more details like options and example about a format use -h <NAME>."
, ""
, ".EX"
, "$ fq -d msgpack d file  # decode as msgpack"
, "$ fq -h msgpack         # see msgpack help"
, "$ fq -h formats         # list formats"
, ".EE"
, ""
, ".EX"
, _help(""; "formats")
, ".EE"

, ".SH SEE ALSO"
, "jq(1)"
, "dd(1)"
, ".SH AUTHOR"
, "Mattias Wadman (mattias.wadman@gmail.com)"
)