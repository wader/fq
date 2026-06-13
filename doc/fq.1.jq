( ".TH fq 1"
, ".SH NAME"
, _help("fq"; "banner")

, ".SH SYNOPSIS"
, _help("fq"; "usage")
, ""
, _help(""; "example_usage")

, ".SH DESCRIPTION"

, _help(""; "summary")
, ""
, "fq is inspired by the well known jq tool and language and allows you to work with binary formats the same way you would using jq. In addition it can present data like a hex viewer, transform, slice and concatenate binary data. It also supports nested formats and has an interactive REPL with auto-completion."

, ".SH EXPRESSION"
, "See "
, ".BR jq (1)"

, ".SH OPTIONS"
, ( _opt_cli_opts
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
      | ".TP"
      , "  -o \($o.key)=<\($o.value)>"
      )
    )
  )

, ".SH ENVIRONMENT"
, ".TP"
, "NO_COLOR"
, "Don't use color output"
, ".TP"
, "CLIUNICODE"
, "Use unicode output"
, ""

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