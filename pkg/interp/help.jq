include "internal";
include "query";
include "eval";
include "repl";
include "decode";
include "funcs";
include "options";
include "args";

# TODO: variants, values, keywords?
# TODO: store some other way?
def _help_functions:
  { length: {
      summary: "Length of string, array, object, etc",
      doc:
"- For string number of unicode codepoints
- For array number of elements in array
- For object number of key-value pairs
- For null zero
- For number the number itself
- For boolean is an error
",
      examples:
        [ [[1,2,3], "length"]
        , ["abc", "length"]
        , [{a: 1, b: 2}, "length"]
        , [null, "length"]
        , [123, "length"]
        , [true, "length"]
        ]
    },
    "..": {
      summary: "Recursive descent of .",
      doc:
"Recursively descend . and output each value.
Same as recurse without argument.
",
      examples:
        [ ["a", ".."]
        , [[1,2,3], ".."]
        , [{a: 1, b: {c: 3}}, ".."]
        ]
    },
    empty: {
      summary: "Output nothing",
      doc:
"Output no value, not even null, and cause backtrack.
",
      examples:
        [ ["empty"]
        , ["[1,empty,2]"]
        ]
    }
  };

def help($_): error("help must be alone or last in pipeline. ex: help(length) or ... | help");
def help: help(null);

def _help($arg0; $topic):
  ( $topic
  | if  . == "usage" then
      "Usage: \($arg0) [OPTIONS] [--] [EXPR] [FILE...]"
    elif . == "example_usage" then
      ( "Example usages:"
      , "  fq . file"
      , "  fq d file"
      , "  fq tovalue file"
      , "  cat file.cbor | fq -d cbor torepr"
      , "  fq 'grep(\"^main$\") | parent' /bin/ls"
      , "  fq 'grep_by(format == \"exif\") | d' *.png *.jpeg"
      )
    elif . == "banner" then
      ( "fq - jq for binary formats"
      , "Tool, language and decoders for inspecting binary data."
      , "For more information see https://github.com/wader/fq"
      )
    elif . == "args" then
      args_help_text(_opt_cli_opts)
    elif  . == "options" then
      ( [ ( options
          | _opt_cli_arg_fromoptions
          )
        | to_entries[]
        | [(.key+"  "), .value | tostring]
        ]
      | table(
          .;
          map(
            ( . as $rc
            # right pad format name to align description
            | if .column == 0 then .string | rpad(" "; $rc.maxwidth)
              else $rc.string
              end
            )
          ) | join("")
        )
      )
    elif . == "formats" then
      ( [ formats
      | to_entries[]
      | [(.key+"  "), .value.description]
      ]
      | table(
          .;
          map(
            ( . as $rc
            # right pad format name to align description
            | if .column == 0 then .string | rpad(" "; $rc.maxwidth)
              else $rc.string
              end
            )
          ) | join("")
        )
      )
    else
      error("unknown topic: \($topic)")
    end
  );

# TODO: refactor
 def _help_slurp($query):
  def _name:
    if _query_is_func then _query_func_name
    else _query_tostring
    end;
  if $query.orig | _query_is_func then
    ( ($query.orig | _query_func_args) as $args
    | ($args | length) as $argc
    | if $args == null then
        # help
        ( "Type expression to evaluate"
        , "\\t          Completion"
        , "Up/Down     History"
        , "^C          Interrupt execution"
        , "... | repl  Start a new REPL"
        , "^D          Exit REPL"
        ) | println
      elif $argc == 1 then
        # help(...)
        ( ($args[0] | _name) as $name
        | _help_functions[$name] as $hf
        | if $hf then
            # help(name)
            ( "\($name): \($hf.summary)"
            , $hf.doc
            , if $hf.examples then
                ( "Examples:"
                , ( $hf.examples[]
                  | . as $e
                  | if length == 1 then
                      ( "> \($e[0])"
                      , (null | try (_eval($e[0]) | tojson) catch "error: \(.)")
                      )
                    else
                      ( "> \($e[0] | tojson) | \($e[1])"
                      , ($e[0] | try (_eval($e[1]) | tojson) catch "error: \(.)")
                      )
                    end
                  )
                )
              end
            ) | println
          else
            # help(unknown)
            # TODO: check builtin
            ( ( . # TODO: extract
              | builtins
              | map(split("/") | {key: .[0], value: true})
              | from_entries
              ) as $builtins
            | ( . # TODO: extract
              | scope
              | map({key: ., value: true})
              | from_entries
              ) as $scope
            | if $builtins | has($name) then
                "\($name) is builtin function"
              elif $scope | has($name) then
                "\($name) is a function or variable"
              else
                "don't know what \($name) is "
              end
            | println
            )
          end
        )
      else
        _eval_error("compile"; "help must be last in pipeline. ex: help(length) or ... | help")
      end
    )
  else
    # ... | help
    # TODO: check builtin
    ( _repl_slurp_eval($query.rewrite) as $outputs
    | "value help"
    , $outputs
    | display
    )
  end;
