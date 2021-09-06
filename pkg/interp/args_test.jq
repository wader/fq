include "assert";
include "args";

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
][] | assert(.name; _args_parse(.args; .opts); .expected)
