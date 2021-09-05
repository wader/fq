include "assert";
include "query";

(
  ([
    ["", "."],
    [".", "."],
    ["a", "a"],
    ["1,  2", "1, 2"],
    ["1 | 2", "2"],
    ["1 | 2 | 3", "3"],
    ["(1 | 2) | 3", "3"],
    ["1 | (2 | 3)", "(2 | 3)"],
    ["1 as $_ | 2", "2"],
    ["def f: 1; 1", "def f: 1; 1"],
    ["def f: 1; 1 | 2", "2"],
    empty
  ][] | assert(
      "\(.) | _query_pipe_last";
      .[1];
      .[0] | _query_fromstring | _query_pipe_last | _query_tostring
    )
  )
,
  ([
    ["", "map(.) | ."],
    [".", "map(.) | ."],
    ["a", "map(.) | a"],
    ["1,  2", "map(.) | 1, 2"],
    ["1 | 2", "map(1 | .) | 2"],
    ["1 | 2 | 3", "map(1 | 2 | .) | 3"],
    ["(1 | 2) | 3", "map((1 | 2) | .) | 3"],
    ["1 | (2 | 3)", "map(1 | .) | (2 | 3)"],
    ["1 as $_ | 2", "map(1 as $_ | .) | 2"],
    ["def f: 1; 1", "map(.) | def f: 1; 1"],
    ["def f: 1; 1 | 2", "map(def f: 1; 1 | .) | 2"],
    ["module {a:1};\ninclude \"a\";\n1", "module { a: 1 };\ninclude \"a\";\nmap(.) | 1"],
    empty
  ][] | assert(
      "\(.) | _query_slurp_wrap";
      .[1];
      .[0] | _query_fromstring | _query_slurp_wrap(.) | _query_tostring
    )
  )
)
