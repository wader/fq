def assert($name; $expected; $actual):
  ( if $expected == $actual then
      "PASS \($name)\n" | stderr
    else
      ( "FAIL \($name): expected \($expected) got \($actual)\n" | stderr
      , (null | halt_error(1))
      )
    end
  );
