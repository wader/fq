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
)
