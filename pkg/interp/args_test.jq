include "args";

def assert($name; $a; $b):
  ( if $a == $b then "PASS \($name)\n"
    else
      ( "FAIL \($name) \($a) != \($b)\n"
      , (null | halt_error(1))
      )
    end
  , empty
  );

[ { name: "Basic parse",
    args: ["-a", "123", "b"],
    opts: {
      "a": {
        short: "-a",
        long: "--abc",
        description: "Set abc",
        string: true
      }
    },
    expected: {   
      "parsed": {
        "a": "123"
      },
      "rest": ["b"]
    }
  }
][] | assert(.name; args_parse(.args; .opts); .expected)
