def assert($name; $a; $b):
  ( if $a == $b then "PASS \($name)\n"
    else
      ( "FAIL \($name) \($a) != \($b)\n"
      , (null | halt_error(1))
      )
    end
  , empty
  );
