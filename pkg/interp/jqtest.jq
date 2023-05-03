# read jq test format:
# # comment
# expr
# input
# output*
# <blank>+
# ...
# <next test>
def from_jqtest:
  [ foreach (split("\n")[], "") as $l (
      { current_line: 0
      , nr: 1
      , emit: true
      };
      ( .current_line += 1
      | if .emit then
          ( .expr = null
          | .input = null
          | .output = []
          | .fail = null
          | .emit = null
          | .error = null
          )
        else .
        end
      | if $l | test("^\\s*#") then .
        elif $l | test("^\\s*$") then
          if .expr then
            ( .emit =
                { line
                , nr
                , expr
                , input
                , output
                , fail
                , error
                }
            | .nr += 1
            )
          else .
          end
        elif $l | test("^\\s*%%FAIL") then
          .fail = $l
        else
          if .expr == null then
            ( .line = .current_line
            | .expr = $l
            )
          elif .fail and .error == null then .error = $l
          elif .input == null then .input = $l
          else .output += [$l]
          end
        end
      );
      if .emit then .emit
      else empty
      end
    )
  ];

def run_tests:
  def _f:
    ( from_jqtest[]
    | . as $c
    | try
        if .error | not then
          ( ( .input
            , .output[]
            ) |= fromjson
          )
        else .
        end
      catch
        ( . as $err
        | $c
        | .fromjson_error = $err
        )
    | select(.fromjson_error | not)
    | "\(.nr) (line \(.line)): [\(.input | tojson) | \(.expr)] -> \(.output | tojson)" as $test_name
    | . as $test
    | try
        ( $test.input
        | [ eval($test.expr)] as $actual_output
        | if $test.output == $actual_output then
            ( empty # "OK: \($test_name)"
            , {ok: true}
            )
          else
            ( "DIFF: \($test_name)"
            , "  Expected: \($test.output | tojson)"
            , "    Actual: \($actual_output | tojson)"
            , {error: true}
            )
          end
        )
      catch
        if $test.fail then
          if . == $test.error then
            ( empty #"OK: \($test_name)"
            , {ok: true}
            )
          else
            ( "FAIL DIFF: \($test_name)"
            , "  Expected: \($test.error)"
            , "    Actual: \(.)"
            , {error: true}
            )
          end
        else
          ( "ERROR: \($test_name)"
          , "  \(.)"
          , {error: true}
          )
        end
    );
  # this mess make it possible to run all tests and exit with non-zero if any test failed
  ( foreach (_f, {end: true}) as $l (
      { errors: 0
      , oks: 0
      };
      if ($l | type) == "object" then
        ( .line = false
        | if $l.error then .errors +=1
          elif $l.ok then .oks += 1
          elif $l.end then .end = true
          else .
          end
        )
      else . + {line: $l}
      end;
      .
    )
  | if .end then
      ( "\(.oks) of \(.oks + .errors) tests passed"
      , if .errors > 0 then null | halt_error(1) else empty end
      )
    elif .line then .line
    else empty
    end
  );
