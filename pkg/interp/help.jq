include "internal";
include "interp";
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
  { length:
      { summary: "Length of string, array, object, etc"
      , doc:
"- For string number of unicode codepoints
- For array number of elements in array
- For object number of key-value pairs
- For null zero
- For number the number itself
- For boolean is an error
"
      , examples:
          [ [[1,2,3], "length"]
          , ["abc", "length"]
          , [{a: 1, b: 2}, "length"]
          , [null, "length"]
          , [123, "length"]
          , [true, "length"]
          ]
      }
  , "..":
      { summary: "Recursive descent of ."
      , doc:
"Recursively descend . and output each value.
Same as recurse without argument.
"
      , examples:
          [ ["a", ".."]
          , [[1,2,3], ".."]
          , [{a: 1, b: {c: 3}}, ".."]
          ]
      }
  , empty:
      { summary: "Output nothing"
      , doc:
"Output no value, not even null, and cause backtrack.
"
      , examples:
          [ ["empty"]
          , ["[1,empty,2]"]
          ]
      }
  };

def _help_format_enrich($arg0; $f; $include_basic):
  ( if $include_basic then
      .examples +=
        [ {comment: "Decode file as \($f.name)", shell: "fq -d \($f.name) . file"}
        , {comment: "Decode value as \($f.name)", expr: "\($f.name)"}
        ]
    end
  | if $f.decode_in_arg then
      .examples +=
        [ { comment: "Decode file using \($f.name) options"
          , shell: "\($arg0) -d \($f.name)\($f.decode_in_arg | to_entries | map(" -o ", .key, "=", (.value | tojson)) | join("")) . file"
          }
        , { comment: "Decode value as \($f.name)"
          , expr: "\($f.name)(\($f.decode_in_arg | to_jq))"
          }
        ]
    end
  );

# trailing help gets rewritten to _help_slurp, these are here to catch other variants
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
      , "  fq -V '.path[1].value' file"
      , "  fq tovalue file"
      , "  fq -r to_toml file.yml"
      , "  fq -s -d html 'map(.html.head.title?)' *.html"
      , "  cat file.cbor | fq -d cbor torepr"
      , "  fq 'grep(\"^main$\") | parent' /bin/ls"
      , "  fq -i"
      )
    elif . == "banner" then
      ( "fq - jq for binary formats"
      , "Tool, language and decoders for working with binary data."
      , "For more information see https://github.com/wader/fq"
      )
    elif . == "args" then
      args_help_text(_opt_cli_opts)
    elif . == "options" then
      ( [ ( options
          | _opt_cli_arg_from_options
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
    elif _registry.formats | has($topic) then
      ( _registry.formats[$topic] as $f
      | (_format_func($f.name; "_help")? // {} | _help_format_enrich($arg0; $f; true)) as $fhelp
      | ((_registry.files[][] | select(.name=="\($topic).md").data) // false) as $doc
      | "\($f.name): \($f.description) decoder"
      , ""
      , if $f.decode_in_arg then
          ( $f.decode_in_arg
          | to_entries
          | map(["  \(.key)=\(.value | tojson)  ", $f.decode_in_arg_doc[.key]])
          | "Options"
          , "======="
          , ""
          , table(
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
          , ""
          )
        else empty
        end
      , "Decode examples"
      , "==============="
      , ""
      , ( $fhelp.examples[]
        | "  # \(.comment)"
        , if .shell then "  $ \(.shell)"
          elif .expr then "  ... | \(.expr)"
          else empty
          end
        )
      , ""
      , if $doc then $doc | markdown | _markdown_to_text(options.width; -2)
        else empty
        end
      )
    elif _help_functions | has($topic) then
      ( _help_functions[$topic] as $hf
      | "\($topic): \($hf.summary)"
      , $hf.doc
      , if $hf.examples then
          ( "Examples:"
          , ( $hf.examples[]
            | . as $e
            | if length == 1 then
                ( "> \($e[0])"
                , (null | try (_eval($e[0]; {}) | tojson) catch "error: \(.)")
                )
              else
                ( "> \($e[0] | tojson) | \($e[1])"
                , ($e[0] | try (_eval($e[1]; {}) | tojson) catch "error: \(.)")
                )
              end
            )
          )
        end
      )
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
      | if $builtins | has($topic) then
          "\($topic) is builtin function"
        elif $scope | has($topic) then
          "\($topic) is a function or variable"
        else
          "don't know what \($topic) is "
        end
      | println
      )
    end
  );

# TODO: refactor
 def _help_slurp($query):
  def _name:
    if _query_is_func then _query_func_name
    elif _query_is_string then _query_string_str
    else _query_tostring
    end;
  if $query.orig | _query_is_func then
    ( ($query.orig | _query_func_args) as $args
    | ($args | length) as $argc
    | if $args == null then
        # help
        ( "Type jq expression to evaluate"
        , "help(...)   Help for topic. Ex: help(mp4), help(\"mp4\")"
        , "\\t          Completion"
        , "Up/Down     History"
        , "... | repl  Start a new REPL"
        , "^C          Interrupt execution"
        , "^D          Exit REPL"
        ) | println
      elif $argc == 1 then
        ( _help("fq"; $args[0] | _name)
        | println
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
