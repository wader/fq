def log: if env.VERBOSE then stderr else empty end;

def assert($name; $expected; $actual):
  ( if $expected == $actual then
      "PASS \($name)\n" | log
    else
      ( "FAIL \($name): expected \($expected) got \($actual)\n" | stderr
      , (null | halt_error(1))
      )
    end
  );
