def log: if env.VERBOSE then printerrln else empty end;

def assert($name; $expected; $actual):
  ( if $expected == $actual then
      "PASS \($name)" | log
    else
      ( "FAIL \($name): expected \($expected) got \($actual)" | printerrln
      , (null | halt_error(1))
      )
    end
  );
